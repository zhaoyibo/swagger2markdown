package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"windmt.com/swagger2markdown/swagger"
	"windmt.com/swagger2markdown/tool"
)

//const url = "http://testmanager.wb-intra.com/wim-manager/v2/api-docs"

var (
	h       bool
	url     string
	project string
	file    string
	path    string
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助")
	flag.StringVar(&url, "url", "", "要导出文档的`Swagger API 的 url`, 例如：http://testmanager.wb-intra.com/wim-manager/v2/api-docs")
	flag.StringVar(&project, "project", "", "要导出文档要的`工程名`, 例如：wim-manager（若 -project 和 -url 同时存在，则 -url 的优先级更高）")
	flag.StringVar(&file, "file", "", "保存结果的`文件名`，例如：wim-manager.md 若文件已存在则会被覆盖")
	flag.StringVar(&path, "path", "", "要导出 `path 路径`，例如：/clue/createOrUpdate")

	// 改变默认的 Usage
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	if project == "" && url == "" {
		flag.Usage()
		return
	} else if url != "" {

	} else if project != "" {
		url = "http://testmanager.wb-intra.com/" + project + "/v2/api-docs"
	}

	if file != "" {
		if tool.CheckFileIsExist(file) {
			var overwrite string
			fmt.Println("文件已存在，是否确定覆盖: y/N")
			if _, err := fmt.Scanf("%s", &overwrite); err != nil {
				fmt.Println("exit, bye!")
				return
			}
			if !strings.EqualFold(overwrite, "y") && !strings.EqualFold(overwrite, "yes") {
				fmt.Println("exit, bye!")
				return
			}
		}
	} else {
		file = tool.ExtractFilename(url, path)
	}

	log.Printf("Swagger API 地址：%s\n", url)
	log.Printf("保存的文件名：%s\n", file)

	if path == "" {
		if err := swagger.ParseAll(url, file); err != nil {
			log.Fatal("Error:", err)
		}
	} else {
		if err := swagger.ParseOne(url, path, file); err != nil {
			log.Fatal("Error:", err)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `swagger2markdown version: 1.0
Usage: swagger2markdown [-h] [-project project] [-url url] [-file filename]
Example: swagger2markdown -project wim-manager -file wim-manager.md

Options:
`)
	flag.PrintDefaults()
}
