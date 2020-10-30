package swagger

import (
	"fmt"
	"github.com/virtuald/go-ordered-json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"windmt.com/swagger2markdown/model"
	"windmt.com/swagger2markdown/tool"
)

const testDataFile string = "data/swagger.json"

func interfaceToString(i interface{}) string {
	if i == nil {
		return ""
	}

	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return fmt.Sprintf("%d", i.(int))
	case float64:
		return fmt.Sprintf("%d", int(i.(float64)))
	}

	if bytes, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return fmt.Sprintf("%s", bytes)
	}
}

func getRootFromLocalFile() (*model.Root, error) {
	file, err := os.Open(testDataFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var root *model.Root
	err = json.NewDecoder(file).Decode(&root)
	return root, err
}

func getRootFromUrl(url string) (*model.Root, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error：", err)
	}

	defer resp.Body.Close()

	var root *model.Root
	err = json.NewDecoder(resp.Body).Decode(&root)

	return root, err
}

func yOrN(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}

func contains(arr []string, target string) bool {
	if arr == nil || len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

var cookie = ""
var contentType = "application/json;charset=UTF-8"

func testLogin() {
	if cookie != "" {
		return
	}
	url := tool.GetDomain() + "/" + tool.GetProject() + "/testLogin"
	resp, err := http.Post(url, contentType, strings.NewReader(""))
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer resp.Body.Close()

	cookies := make([]string, 0, len(resp.Cookies()))
	for _, v := range resp.Cookies() {
		cookies = append(cookies, v.Name+"="+v.Value)
	}

	cookie = strings.Join(cookies, "; ")
}

func getResponseExample(apiUrl string, param string) string {
	testLogin()

	client := &http.Client{}

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(param))
	if err != nil {
		log.Fatal("Error:", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error:", err)
	}

	tmpMap := make(map[string]interface{})
	if err = json.Unmarshal(body, &tmpMap); err != nil {
		log.Fatal("Error:", err)
	}

	if data, ok := tmpMap["data"].(map[string]interface{}); ok {
		if list, ok := data["list"]; ok {
			if list, ok := list.([]interface{}); ok {
				if len(list) > 2 {
					data["list"] = list[:2]
				}
			}
		}
	}

	if bytes, err := json.MarshalIndent(tmpMap, "", "    "); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

func getDescMarkdown(path model.Path, basePath string, targetPath string, methods string, baseTitleLevel int) string {
	builder := &tool.StringBuilder{}
	builder.Append(strings.Repeat("#", baseTitleLevel)).Append(" ").Append("接口描述").Br()
	builder.Append("功能：").Append(path.Summary).Br2()
	builder.Append("地址：").Append(basePath).Append(targetPath).Br2()
	builder.Append("方法：").Append("POST").Br2()
	return builder.String()
}

func getRequestMarkdown(path model.Path, root *model.Root, baseTitleLevel int) (markdown, exampleParams string) {
	builder := &tool.StringBuilder{}

	builder.Append(strings.Repeat("#", baseTitleLevel)).Append(" ").Append("请求参数").Br2()
	parameters := path.Parameters
	if parameters == nil {

		builder.Append("无").Br2()
	} else if len(parameters) > 0 {
		firstParameter := parameters[0]
		isBody := firstParameter.IsBody()

		if path.Consumes == nil {
			builder.Append("Content-Type：").Append("application/x-www-form-urlencoded").Br2()
		} else {
			builder.Append("Content-Type：`").Append(strings.Join(path.Consumes, "`, `")).Append("`").Br2()
		}

		if isBody {

			definition := root.Definitions[firstParameter.TypeName()]

			buildModels(true, builder, definition, root)

			builder.Br2()
			builder.Append("请求示例：").Br2()

			requestExampleMap := buildRequestExampleMap(definition, root)

			if bytes, err := json.Marshal(requestExampleMap); err != nil {
				panic(err)
			} else {
				exampleParams = string(bytes)
			}

			builder.Append("```json").Br()

			if bytes, err := json.MarshalIndent(requestExampleMap, "", "    "); err != nil {
				panic(err)
			} else {
				builder.Append(string(bytes)).Br()
			}

			builder.Append("```").Br2()

		} else {
			requestExampleMap := make(map[string]interface{})

			builder.Append("|参数|类型|必须|说明|示例|默认值|").Br()
			builder.Append("|----|----|----|----|----|----|").Br()
			for _, p := range parameters {
				builder.Append("|").Append(p.Name)
				builder.Append("|").Append(p.TypeName())
				builder.Append("|").Append(yOrN(p.Required))
				builder.Append("|").Append(p.Description)
				builder.Append("|").Append(interfaceToString(p.Example))
				if p.Default == nil {
					builder.Append("|").Append("无")
				} else {
					builder.Append("|").Append(interfaceToString(p.Default))
				}
				builder.Append("|").Br()

				requestExampleMap[p.Name] = p.Example
			}

			builder.Br2()
			builder.Append("请求示例：").Br2()

			builder.Append("```plain").Br()

			params := make([]string, 0, len(requestExampleMap))
			for k, v := range requestExampleMap {
				params = append(params, k+"="+interfaceToString(v))
			}

			exampleParams = strings.Join(params, "&")
			builder.Append(exampleParams).Br()

			builder.Append("```").Br2()
		}
	}
	return builder.String(), exampleParams
}

func buildRequestExampleMap(definition model.Definition, root *model.Root) map[string]interface{} {
	requestExampleMap := make(map[string]interface{})
	for _, v := range definition.Properties() {
		if !v.ReadOnly {
			if v.RefRaw != "" {
				requestExampleMap[v.Name] = buildRequestExampleMap(root.Definitions[v.Ref()], root)
			} else {
				requestExampleMap[v.Name] = v.Example
			}
		}
	}
	return requestExampleMap
}

func getResponseMarkdown(path model.Path, root *model.Root, requestUrl string, exampleParams string, methods string, baseTitleLevel int) string {
	builder := &tool.StringBuilder{}

	builder.Append(strings.Repeat("#", baseTitleLevel)).Append(" ").Append("返回响应").Br2()
	builder.Append("Accept：`").Append(strings.Join(path.Produces, "`, `")).Append("`").Br2()

	builder.Append(strings.Repeat("#", baseTitleLevel+1)).Append(" ").Append("响应码").Br2()
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

	builder.Append(strings.Repeat("#", baseTitleLevel+1)).Append(" ").Append("响应体").Br2()
	schema := path.Responses["200"].Schema

	if strings.EqualFold(schema.Ref(), "ManagerResponse") || !strings.Contains(methods, "GET") {
		builder.Append("```json").Br()
		responseExampleMap := make(map[string]interface{})
		responseExampleMap["code"] = 0
		responseExampleMap["msg"] = "success"
		responseExampleMap["data"] = make(map[string]interface{})
		if bytes, err := json.MarshalIndent(responseExampleMap, "", "    "); err != nil {
			panic(err)
		} else {
			builder.Append(string(bytes)).Br()
		}
		builder.Append("```").Br2()
	} else {
		if definition, ok := root.Definitions[schema.Ref()]; ok {
			buildModels(false, builder, definition, root)
		}

		builder.Append("响应示例：").Br2()
		builder.Append("```json").Br()
		example := getResponseExample(requestUrl, exampleParams)
		builder.Append(example).Br()
		builder.Append("```").Br2()
	}

	return builder.String()
}

func buildModels(isRequest bool, builder *tool.StringBuilder, definition model.Definition, root *model.Root) {
	subDefinitions := make([]model.Definition, 0)
	builder.Append("**").Append(definition.Title).Append("：**").Br2()
	if len(definition.Properties()) == 0 {
		builder.Append("无")
		return
	}

	if isRequest {
		builder.Append("|字段|类型|必须|说明|示例|").Br()
	} else {
		builder.Append("|字段|类型|非空|说明|示例|").Br()
	}

	builder.Append("|----|----|----|----|----|").Br()

	for _, property := range definition.Properties() {
		if isRequest && property.ReadOnly {
			continue
		}

		builder.Append("|").Append(property.Name)
		if property.Type != "" {
			builder.Append("|").Append(property.Type)
		} else if property.Ref() != "" {
			builder.Append("|").Append(property.Ref())
			subDefinitions = append(subDefinitions, root.Definitions[property.Ref()])
		} else {
			builder.Append("|")
		}
		builder.Append("|").Append(yOrN(contains(definition.Required, property.Name)))
		builder.Append("|").Append(property.Description)
		builder.Append("|").Append(interfaceToString(property.Example))
		builder.Append("|").Br()
	}
	builder.Br2()

	if len(subDefinitions) > 0 {
		for _, subDefinition := range subDefinitions {
			buildModels(isRequest, builder, subDefinition, root)
		}
	}
}

func getRequestUrl(basePath, path string) string {
	return tool.GetDomain() + basePath + path
}

func sortedMap(m map[string]interface{}, f func(k string, v interface{})) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}
