/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a shell command",
	Long: `Execute a shell command and display the output. 
	If an error occurs during command execution, it will be analyzed using the Gemini API 
	to provide additional insights and suggestions for troubleshooting.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellCmd := exec.Command(args[0], args[1:]...)

		output, cmdErr := shellCmd.CombinedOutput()
		if cmdErr != nil {
			fmt.Println("Error executing command: ", cmdErr)

			// Analyze the error with the Gemini API
			client, err := genai.NewClient(cmd.Context(), option.WithAPIKey(viper.GetString("gemini.api-key")))
			if err != nil {
				fmt.Println("Error creating client: ", err)
				os.Exit(1)
			}
			defer client.Close()

			model := client.GenerativeModel("gemini-pro")
			iter := model.GenerateContentStream(context.Background(), genai.Text(cmdErr.Error()))
			fmt.Println("\nGemini analysis of the error:\n")

			for {
				resp, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Print("Error getting response: ", err)
					break
				}

				for _, part := range resp.Candidates[0].Content.Parts {
					fmt.Println(part)
				}
			}

			return
		}

		fmt.Println("Command output: ", string(output))
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
