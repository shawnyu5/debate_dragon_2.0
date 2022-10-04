package generatedocs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generate read all readmes from the command files, and generate a repo readme
func Generate() {
	commandDescriptions := readAllReadmes()
	t := template.Must(template.ParseFiles("generate_docs/templates/readme.md"))
	f, err := os.OpenFile("README.md", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(f, commandDescriptions)

}

// readAllReadmes read all readme from all command folders
// Return an array of strings of all the readme contents, with leading and trailing white space trimmed
func readAllReadmes() []string {
	commandDescriptions := []string{}
	// read all command readmes
	err := filepath.Walk("./commands",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// find all readmes in the commands folder, and store their path
			if strings.ToLower(info.Name()) == "readme.md" {
				f, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				commandDescriptions = append(commandDescriptions, strings.TrimSpace(string(f)))
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return commandDescriptions
}
