/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	query string
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question to Gemini",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if !viper.InConfig("gemini.api-key") {
			fmt.Println("No API Key configured. Please run `go-gemini config` to configure your API Key.")
			os.Exit(1)
		}

		client, err := genai.NewClient(ctx, option.WithAPIKey(viper.GetString("gemini.api-key")))
		if err != nil {
			fmt.Print("Error creating client: ", err)
			os.Exit(1)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-pro")
		resp, err := model.GenerateContent(ctx, genai.Text(query))
		if err != nil {
			fmt.Print("Error generating content: ", err)
			os.Exit(1)
		}

		fmt.Println(resp.Candidates[0].Content)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().StringVarP(&query, "query", "q", "", "The query to ask Gemini")
}
