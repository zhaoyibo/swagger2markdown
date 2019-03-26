package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strings"
	"windmt.com/swagger2markdown/swagger"
	"windmt.com/swagger2markdown/tool"
)

var opts struct {
	Url     string `short:"u" long:"url" value-name:"URL" description:"要导出文档的 Swagger API 的 url" long-description:"要导出文档的Swagger API 的 url, 例如：\nhttp://testmanager.wb-intra.com/wim-manager/v2/api-docs"`
	Project string `short:"p" long:"project" value-name:"PROJECT" description:"要导出文档要的工程名, 例如：wim-manager\np.s. -u 的优先级更高"`
	File    string `short:"f" long:"file" description:"保存结果的文件名，例如：wim-manager.md 若文件已存在则会被覆盖" value-name:"FILE"`
	Path    string `short:"P" long:"path" value-name:"PATH" description:"要导出 path 路径，例如：/text/list"`
}

func main() {
	var parser = flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	var url string
	if opts.Project == "" && opts.Url == "" {
		fmt.Fprintf(os.Stderr, "Missing flag [-u|-p]\nRun '%s -h' for usage.\n", os.Args[0])
		os.Exit(1)
	} else if opts.Url != "" {
		url = opts.Url
	} else if opts.Project != "" {
		url = "http://testmanager.wb-intra.com/" + opts.Project + "/v2/api-docs"
	}

	file := opts.File
	path := opts.Path
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
