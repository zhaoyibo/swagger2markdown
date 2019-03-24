package swagger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"windmt.com/swagger2markdown/model"
)

const testDataFile string = "data/swagger.json"

func interfaceToString(i interface{}) string {
	if i == nil {
		return ""
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
		return false;
	}
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
