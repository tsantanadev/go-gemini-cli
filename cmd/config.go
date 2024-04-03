/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	apiKey string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config command is used to configure the Gemini API Key",
	Long: `The config command is used to configure the Gemini API Key. It allows you to set the API key that will be used for making requests to the Gemini API. 
	You can provide the API key using the --api-key flag. If the API key is provided, it will be stored in the configuration file. 
	Subsequent requests to the Gemini API will use this stored API key. If no API key is provided, an error message will be displayed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey != "" {
			viper.Set("gemini.api-key", apiKey)

			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error writing config file:", err)
			}

		} else {
			fmt.Println("No API Key provided")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "The API key to use for Gemini")
}
