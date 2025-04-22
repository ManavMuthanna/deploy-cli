package dockergen

import (
	"fmt"
	"html/template"
	"os"
)

// GenerateDockerfile generates a Dockerfile based on the provided language
func GenerateDockerfile(lang string) error {
	dockerTemplate, err := template.ParseFiles(
		fmt.Sprintf("features/dockergen/templates/%s.dockerfile", lang),
	)
	if err != nil {
		return err
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	defer file.Close()
	return dockerTemplate.Execute(file, nil)
}
