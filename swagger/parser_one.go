package swagger

import (
	_ "encoding/json"
	"strings"
	"windmt.com/swagger2markdown/model"
	"windmt.com/swagger2markdown/tool"
)

func ParseOne(swaggerApiUrl, targetPath, file string) error {
	testLogin()
	root, err := getRootFromUrl(swaggerApiUrl)
	if err != nil {
		return err
	}
	basePath := root.BasePath

	builder := &tool.StringBuilder{}

	paths := root.Paths[targetPath]

	methods := make([]string, 0, len(paths))
	var path model.Path
	for method, v := range paths {
		if len(methods) == 0 {
			path = v
		}
		methods = append(methods, strings.ToUpper(method))
	}

	baseTitleLevel := 1

	methodsStr := strings.Join(methods, ", ")
	desc := getDescMarkdown(path, basePath, targetPath, methodsStr, baseTitleLevel)
	builder.Append(desc).Br2()

	request, exampleParams := getRequestMarkdown(path, root, baseTitleLevel)
	builder.Append(request).Br2()

	response := getResponseMarkdown(path, root, getRequestUrl(basePath, targetPath), exampleParams, methodsStr, baseTitleLevel)
	builder.Append(response).Br2()

	return tool.SaveToFile(file, builder.String())
}
