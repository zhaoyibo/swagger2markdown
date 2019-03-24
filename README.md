# swagger2markdown
Swagger 文档转 Markdown

## Usage

```
swagger2markdown version: 1.0
Usage: swagger2markdown [-h] [--project project] [--url url] [--path] [--file filename]
Example: swagger2markdown --project wim-manager --path /text/list

Options:
  -file 文件名
    	保存结果的文件名，例如：wim-manager.md 若文件已存在则会被覆盖
  -h	帮助
  -path path 路径
    	要导出 path 路径，例如：/text/list
  -project 工程名
    	要导出文档要的工程名, 例如：wim-manager（若 -project 和 -url 同时存在，则 -url 的优先级更高）
  -url Swagger API 的 url
    	要导出文档的Swagger API 的 url, 例如：http://testmanager.wb-intra.com/wim-manager/v2/api-docs

```