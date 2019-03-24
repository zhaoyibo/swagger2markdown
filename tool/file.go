package tool

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

const dir string = "api-documents/"

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func SaveToFile(dstFile, content string) error {
	if err := mkdirIfNotExist(dir); err != nil {
		log.Fatal("Error:", err)
	}

	f, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(content)
	if err != nil {
		return err
	}
	return w.Flush()
}

func formatFilename(filename string) string {
	lastIndex := strings.LastIndex(filename, ".")
	if lastIndex > 0 {
		now := time.Now().Format("20060102150405")
		pre := string([]rune(filename)[:lastIndex])
		suf := string([]rune(filename)[lastIndex:])
		return pre + "_" + now + suf
	} else {
		return filename
	}
}

func ExtractFilename(url, path string) string {
	url = strings.ReplaceAll(url, "/v2/api-docs", "")
	lastIndex := strings.LastIndex(url, "/")

	var filename string
	if path == "" {
		filename = string([]rune(url)[lastIndex+1:]) + ".md"
	} else {
		filename = string([]rune(url)[lastIndex+1:]) + "_" + strings.ReplaceAll(path, "/", "|") + ".md"
	}
	filename = dir + formatFilename(filename)
	return filename
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mkdirIfNotExist(dir string) error {
	exist, err := pathExists(dir)
	if err != nil {
		return err
	}

	if !exist {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}
