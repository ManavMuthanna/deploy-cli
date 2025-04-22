package dockergen

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

// Function to detect the project type
func DetectProjectType() string {
	items := []string{"Go", "Node", "Python"}
	prompt := promptui.Select{
		Label: "Select your Project Type",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Println("Error selecting project type:", err)
		return ""
	}
	return strings.ToLower(result)
}

// Helper function to check if a file exists
func DockerfileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
