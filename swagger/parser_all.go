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
	pathMapGroupByTag := make(map[string]map[string]model.Path)

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
					pathMap = make(map[string]model.Path)
					pathMapGroupByTag[tag] = pathMap
				}
				pathMap[uri] = path
			}
			methods = append(methods, strings.ToUpper(method))
		}

		pathMethodMap[uri] = strings.Join(methods, ", ")
	}

	//fmt.Printf("%+v\n", pathMapGroupByTag)
	//fmt.Printf("%+v\n", pathMethodMap)

	for k, pathMap := range pathMapGroupByTag {
		builder.Append("## ").Append(k).Br2()

		for uri, path := range pathMap {
			builder.Append("### ").Append(basePath).Append(uri).Br2()

			builder.Append("#### 接口描述").Br2()
			builder.Append(path.Summary).Br2()
			builder.Append("地址：").Append(basePath).Append(uri).Br2()
			builder.Append("协议：").Append(pathMethodMap[uri]).Br2()
			if path.Consumes == nil {
				builder.Append("consumes：").Append("application/x-www-form-urlencoded").Br2()
			} else {
				builder.Append("consumes：`").Append(strings.Join(path.Consumes, "`, `")).Append("`").Br2()
			}
			builder.Append("produces：`").Append(strings.Join(path.Produces, "`, `")).Append("`").Br2()

			builder.Append("#### 请求参数").Br2()
			if path.Parameters == nil {
				builder.Append("无").Br2()
			} else {
				builder.Append("|请求类型|参数|类型|必须|说明|默认值|").Br()
				builder.Append("|----|----|----|----|----|----|").Br()

				for _, parameter := range path.Parameters {
					builder.Append("|").Append(parameter.In)
					builder.Append("|").Append(parameter.Name)
					builder.Append("|").Append(parameter.TypeName())
					builder.Append("|").Append(yOrN(parameter.Required))
					builder.Append("|").Append(parameter.Description)

					if parameter.IsBody() {
						builder.Append("（见下方模型定义中的`").Append(parameter.TypeName()).Append("`）")
					}

					if parameter.Default == nil {
						builder.Append("|").Append("无")
					} else {
						builder.Append("|").Append(interfaceToString(parameter.Default))
					}
					builder.Append("|").Br()
				}
				builder.Br2()
			}

			builder.Append("#### 返回响应").Br2()
			builder.Append("##### 响应码").Br2()
			builder.Append("|code|说明|").Br()
			builder.Append("|----|----|").Br()

			for code, resp := range path.Responses {
				switch code {
				case "200":
					builder.Append("|0|").Append(resp.Description).Append("|").Br()
				case "500":
					builder.Append("|-1|").Append(resp.Description).Append("|").Br()
				}
			}
			builder.Br()

			builder.Append("##### 响应体").Br2()
			schema := path.Responses["200"].Schema
			if schema.IsVoid() {
				builder.Append("无").Br2()
			} else {
				builder.Append("|字段|类型|说明|").Br()
				builder.Append("|----|----|----|").Br()
				if schema.IsList() {
					builder.Append("|data.list|").Append("List\\<`" + schema.Ref() + "`>").Append("|见下方模型定义中的`" + schema.Ref() + "`|").Br()
					builder.Append("|data.max|int|始终为0|").Br()
				} else {
					builder.Append("|data.one|`").Append(schema.Ref()).Append("`|见下方模型定义中的`" + schema.Ref() + "`|").Br()
				}
			}

			builder.Br2()
		}

		builder.Br2()
	}

	builder.Append("## 模型定义").Br2()
	for name, definition := range root.Definitions {
		if strings.EqualFold(name, "MsResponse") {
			continue
		}
		builder.Append("### ").Append(name).Br2()
		builder.Append("|字段|类型|必须|只读|说明|示例|").Br()
		builder.Append("|----|----|----|----|----|----|").Br()

		for field, property := range definition.Properties {
			builder.Append("|").Append(field)
			builder.Append("|").Append(property.Type)
			builder.Append("|").Append(yOrN(contains(definition.Required, field)))
			builder.Append("|").Append(yOrN(property.ReadOnly))
			builder.Append("|").Append(property.Description)
			builder.Append("|").Append(interfaceToString(property.Example))
			builder.Append("|").Br()
		}
		builder.Br2()
	}

	//fmt.Println(builder.String())

	return tool.SaveToFile(dst, builder.String())
}