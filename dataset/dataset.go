package dataset

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadDataset(dir string, itemFunc func(body, category string)) (map[string]int, error) {
	categories := map[string]int{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == ".DS_Store" {
			return nil
		}

		parts := strings.Split(path, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}
		category := parts[len(parts)-2]
		categories[category]++
		b, err := ioutil.ReadFile(path)
		body := string(b)
		if err != nil {
			log.Println("[ERROR] reading the file ", path, err)
			return nil
		}
		// skip message headers
		skip := strings.Index(body, "\n\n")
		body = body[skip+1 : len(body)]
		itemFunc(sanitize(body), category)
		return nil
	})
	return categories, err
}

func sanitize(s string) string {
	r := strings.NewReplacer("\n", "", "\r", "", "\t", " ", "|", " ", ":", " ", "_", " ")
	return r.Replace(s)
}
