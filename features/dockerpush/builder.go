package dockerpush

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

// getDockerCredentials loads from env or prompts the user if missing
func getDockerCredentials() (string, string, error) {
	_ = godotenv.Load("deploy-config.env") // Ignore error if file doesn't exist

	dockerUsername := os.Getenv("DOCKERHUB_USERNAME")
	dockerPassword := os.Getenv("DOCKERHUB_PASSWORD")

	if dockerUsername == "" {
		prompt := promptui.Prompt{Label: "Docker Hub Username"}
		result, err := prompt.Run()
		if err != nil {
			return "", "", fmt.Errorf("failed to get Docker Hub username: %v", err)
		}
		dockerUsername = result
	}

	if dockerPassword == "" {
		prompt := promptui.Prompt{Label: "Docker Hub Password", Mask: '*'}
		result, err := prompt.Run()
		if err != nil {
			return "", "", fmt.Errorf("failed to get Docker Hub password: %v", err)
		}
		dockerPassword = result
	}

	return dockerUsername, dockerPassword, nil
}

// isDockerAuthenticated returns true if already logged in to Docker Hub
func isDockerAuthenticated() bool {
	authCmd := exec.Command("docker", "info", "--format={{.RegistryConfig.IndexConfigs.docker.io.Auth}}")
	authOutput, _ := authCmd.CombinedOutput()
	return strings.TrimSpace(string(authOutput)) != ""
}

// BuildDockerImage handles image building with version management
func BuildDockerImage() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter image name (e.g. myrepo/app)",
	}
	baseName, err := prompt.Run()
	if err != nil || baseName == "" {
		return "", fmt.Errorf("error getting image name: %v", err)
	}

	imageName := baseName + ":1.0"
	if imageExistsLocally(imageName) {
		fmt.Println("Image already exists locally, incrementing version...")
		inc_check := selectIncType()
		imageName = incrementVersion(imageName, inc_check)
	}

	fmt.Println("Running Docker Build...")

	buildCmd := exec.Command("docker", "build", "-t", imageName, ".")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build failed: %s\n%s", err, output)
	}

	return imageName, nil
}

// PushImage builds a Docker image and pushes it to Docker Hub
func PushDockerImage(imageName string) error {
	if isDockerAuthenticated() {
		fmt.Println("Docker Hub is already authenticated, skipping login...")
	} else {
		dockerUsername, dockerPassword, err := getDockerCredentials()
		if err != nil {
			return err
		}
		loginCmd := exec.Command("docker", "login", "-u", dockerUsername, "-p", dockerPassword)
		loginOutput, loginErr := loginCmd.CombinedOutput()
		if loginErr != nil {
			return fmt.Errorf("failed to login to Docker Hub: %s", loginOutput)
		}
	}

	fmt.Printf("About to push image: %s\n", imageName)

	// Push the image to Docker Hub
	pushCmd := exec.Command("docker", "push", imageName)
	pushOutput, pushErr := pushCmd.CombinedOutput()
	if pushErr != nil {
		return fmt.Errorf("failed to push Docker image: %s", pushOutput)
	}

	return nil
}
