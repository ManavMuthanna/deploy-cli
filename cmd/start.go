package cmd

import (
	"deploy-cli/features/dockergen"
	"deploy-cli/features/dockerpush"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deploy-cli",
	Short: "deploy-cli is a cli tool for deploying your projects through the terminal",
	Long:  "deploy-cli is a cli tool for deploying your projects from the terminal in a fast, easy and secure way",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Dockerfile generation
var dockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "Generate a Dockerfile for your project",
	Run: func(cmd *cobra.Command, args []string) {
		projectType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Println("Error getting project type:", err)
			return
		}

		if projectType == "" {
			projectType = dockergen.DetectProjectType()
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Println("Error getting force flag:", err)
			return
		}

		if !force && dockergen.DockerfileExists("Dockerfile") {
			log.Println("Error: Dockerfile already exists. Use --force to overwrite.")
			return
		}

		err = dockergen.GenerateDockerfile(projectType)
		if err != nil {
			fmt.Println("Error generating Dockerfile:", err)
		} else {
			fmt.Println("Dockerfile generated successfully.")
		}
	},
}

// Docker image push
var dockerpushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a Docker Image to DockerHub",
	Run: func(cmd *cobra.Command, args []string) {
		newImage, _ := cmd.Flags().GetBool("new")
		var imageName string
		var err error

		if err != nil {
			fmt.Println("Error getting image name:", err)
			return
		}

		if newImage {
			// Build new versioned image
			imageName, err = dockerpush.BuildDockerImage()
			if err != nil {
				fmt.Println("Error building image:", err)
				return
			}
		} else {
			// Existing image selection logic
			imageName, err = cmd.Flags().GetString("image")
			if err != nil || imageName == "" {
				imageName, err = dockerpush.SelectDockerImage()
				if err != nil {
					fmt.Println("Error selecting image:", err)
					return
				}
			}
		}

		err = dockerpush.PushDockerImage(imageName)

		if err != nil {
			fmt.Println("Error pushing Docker image:", err)
		} else {
			fmt.Println("Docker image pushed successfully.")
		}
	},
}

// Initialize flags
func init() {
	dockerfileCmd.Flags().StringP("type", "t", "", "project type (e.g. go, node, python)")
	dockerfileCmd.Flags().BoolP("force", "f", false, "overwrite existing Dockerfile")

	dockerpushCmd.Flags().StringP("image", "i", "", "Docker image name")
	dockerpushCmd.Flags().BoolP("new", "n", false, "Build and Push new Docker image")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Add subcommands to root command
	rootCmd.AddCommand(
		dockerfileCmd,
		dockerpushCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}
