# swagger2markdown
Swagger 文档转 Markdown

## Usage

```
swagger2markdown version: 1.0.1
Usage:
  swagger2markdown [OPTIONS]

Application Options:
  -u, --url=URL            要导出文档的 Swagger API 的 url
  -p, --project=PROJECT    要导出文档要的工程名, 例如：wim-manager
                           p.s. -u 的优先级更高
  -f, --file=FILE          保存结果的文件名，例如：wim-manager.md
                           若文件已存在则会被覆盖
  -P, --path=PATH          要导出 path 路径，例如：/text/list

Help Options:
  -h, --help               Show this help message
```