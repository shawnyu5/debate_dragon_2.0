package generatedocs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Generate() {
	// filepath.SkipDir
	err := filepath.Walk("./commands",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Name())
			return nil
		})
	if err != nil {
		log.Println(err)
	}

}
