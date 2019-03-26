package swagger

import (
	_ "encoding/json"
	"strings"
	"windmt.com/swagger2markdown/model"
	"windmt.com/swagger2markdown/tool"
)

func ParseAll(swaggerApiUrl, dst string) error {
	root, err := getRootFromUrl(swaggerApiUrl)
	if err != nil {
		return err
	}

	basePath := root.BasePath

	builder := &tool.StringBuilder{}

	builder.Append("# 接口文档").Br2()
	builder.Append("BasePath: ").Append(root.BasePath).Br2()
	builder.Append("Swagger UI: ").Append(strings.Replace(swaggerApiUrl, "/v2/api-docs", "/swagger-ui.html#/", 1)).Br2()

	tagMap := make(map[string]model.Tag)
	for _, tag := range root.Tags {
		tagMap[tag.Name] = tag
	}

	// tag -> <pathStr, path>
	pathMapGroupByTag := make(map[string]interface{})

	pathMethodMap := make(map[string]string)

	for uri, paths := range root.Paths {
		if strings.HasPrefix(uri, "/_") || strings.EqualFold(uri, "/ms/*/*") || strings.Contains(uri, "testLogin") || strings.Contains(uri, "managerSchedule") {
			continue
		}

		methods := make([]string, 0, len(paths))
		for method, path := range paths {
			if len(methods) == 0 {
				tag := path.Tags[0]
				pathMap := pathMapGroupByTag[tag]
				if pathMap == nil {
					pathMap = make(map[string]interface{})
					pathMapGroupByTag[tag] = pathMap
				}
				pathMap.(map[string]interface{})[uri] = path
			}
			methods = append(methods, strings.ToUpper(method))
		}

		pathMethodMap[uri] = strings.Join(methods, ", ")
	}

	//fmt.Printf("%+v\n", pathMapGroupByTag)
	//fmt.Printf("%+v\n", pathMethodMap)

	baseTitleLevel := 4

	sortedMap(pathMapGroupByTag, func(tag string, pathMap interface{}) {
		builder.Append("## ").Append(tag).Br2()

		sortedMap(pathMap.(map[string]interface{}), func(uri string, v interface{}) {
			path := v.(model.Path)
			builder.Append("### ").Append(basePath).Append(uri).Br2()

			desc := getDescMarkdown(path, basePath, uri, pathMethodMap[uri], baseTitleLevel)
			builder.Append(desc).Br2()

			request, exampleParams := getRequestMarkdown(path, root, baseTitleLevel)
			builder.Append(request).Br2()

			response := getResponseMarkdown(path, root, getRequestUrl(basePath, uri), exampleParams, pathMethodMap[uri], baseTitleLevel)
			builder.Append(response).Br2()
		})

		builder.Br2()
	})

	return tool.SaveToFile(dst, builder.String())
}
