package dockerpush

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func SelectDockerImage() (string, error) {
	// List images using docker CLI
	out, err := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error listing Docker images: %s", string(out))
	}
	images := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(images) == 0 || (len(images) == 1 && images[0] == "") {
		return "", fmt.Errorf("no Docker images found")
	}

	prompt := promptui.Select{
		Label: "Select Docker image to push",
		Items: images,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}
	return result, nil
}

// incrementVersion increments the version number in the image tag
func incrementVersion(imageName string, incrementMajor bool) string {
	parts := strings.Split(imageName, ":")
	if len(parts) != 2 {
		return imageName + ":1.0"
	}

	repo, version := parts[0], parts[1]
	verParts := strings.Split(version, ".")
	if len(verParts) != 2 {
		return fmt.Sprintf("%s:%s.1", repo, version)
	}

	major, minor := verParts[0], verParts[1]
	majorInt, _ := strconv.Atoi(major)
	minorInt, _ := strconv.Atoi(minor)

	if incrementMajor {
		majorInt++
		minorInt = 0
	} else {
		minorInt++
	}

	return fmt.Sprintf("%s:%d.%d", repo, majorInt, minorInt)
}

// imageExistsLocally checks if the image exists locally
func imageExistsLocally(imageName string) bool {
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	images := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, img := range images {
		if img == imageName {
			return true
		}
	}
	return false
}

func selectIncType() bool {
	items := []bool{true, false}
	prompt := promptui.Select{
		Label: "Is this a major version increment?",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Println("Error selecting Increment Type:", err)
		return false // or return true, depending on the default value you want
	}

	resultMap := map[string]bool{"true": true, "false": false}
	return resultMap[result]
}
