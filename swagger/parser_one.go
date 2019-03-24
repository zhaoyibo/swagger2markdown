package swagger

import (
	_ "encoding/json"
	"strings"
	"windmt.com/swagger2markdown/model"
	"windmt.com/swagger2markdown/tool"
)



func ParseOne(swaggerApiUrl, targetPath, file string) error {
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

	builder.Append("# 接口描述").Br2()
	builder.Append(path.Summary).Br2()
	builder.Append("地址：").Append(basePath).Append(targetPath).Br2()
	builder.Append("协议：").Append(strings.Join(methods, ", ")).Br2()
	if path.Consumes == nil {
		builder.Append("consumes：").Append("application/x-www-form-urlencoded").Br2()
	} else {
		builder.Append("consumes：`").Append(strings.Join(path.Consumes, "`, `")).Append("`").Br2()
	}
	builder.Append("produces：`").Append(strings.Join(path.Produces, "`, `")).Append("`").Br2()

	builder.Append("# 请求参数").Br2()
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
				builder.Append("（见 ").Append(parameter.TypeName()).Append(" ）")
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

	builder.Append("# 返回响应").Br2()
	builder.Append("## 返回码").Br2()
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

	builder.Append("## 响应体").Br2()
	schema := path.Responses["200"].Schema
	if schema.IsVoid() {
		builder.Append("无").Br2()
	} else {
		builder.Append("|字段|类型|说明|").Br()
		builder.Append("|----|----|----|").Br()
		if schema.IsList() {
			builder.Append("|data.list|").Append("List\\<" + schema.Ref() + ">").Append("|见 " + schema.Ref() + " |").Br()
			builder.Append("|data.max|int|始终为0|").Br()
		} else {
			builder.Append("|data.one|").Append(schema.Ref()).Append("|见 " + schema.Ref() + " |").Br()
		}
	}

	builder.Br2()

	return tool.SaveToFile(file, builder.String())
}
