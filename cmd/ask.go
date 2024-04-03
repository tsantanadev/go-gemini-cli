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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	query        string
	shouldStream bool
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question to Gemini",
	Long: `The ask command allows you to ask a question to Gemini, a generative AI model. 
	You can provide a query using the --query flag and receive a response from the model. 
	If you want the response to be streamed to the console, you can use the --stream flag. 
	Please make sure to configure your API Key before using this command by running 'go-gemini-cli config'.`,
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
		if shouldStream {
			iter := model.GenerateContentStream(ctx, genai.Text(query))
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
		} else {
			resp, err := model.GenerateContent(ctx, genai.Text(query))
			if err != nil {
				fmt.Print("Error getting response: ", err)
				os.Exit(1)
			}

			for _, part := range resp.Candidates[0].Content.Parts {
				fmt.Println(part)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().StringVarP(&query, "query", "q", "", "The query to ask Gemini")
	askCmd.Flags().BoolVarP(&shouldStream, "stream", "s", false, "The response will be streamed to the console")
}
