package cmd

import (
	"deploy-cli/features/dockergen"
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

// Initialize flags
func init() {
	dockerfileCmd.Flags().StringP("type", "t", "", "project type (e.g. go, node, python)")
	dockerfileCmd.Flags().BoolP("force", "f", false, "overwrite existing Dockerfile")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Add subcommands to root command
	rootCmd.AddCommand(dockerfileCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}
