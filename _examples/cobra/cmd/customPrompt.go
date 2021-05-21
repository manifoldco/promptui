/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strconv"

	"github.com/spf13/cobra"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

/*
➜  cobra git:(master) ✗ go run main.go customPrompt
Spicy Level 1
You answered 1
customPrompt called
 */


// customPromptCmd represents the customPrompt command
var customPromptCmd = &cobra.Command{
	Use:   "customPrompt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		validate := func(input string) error {
			_, err := strconv.ParseFloat(input, 64)
			return err
		}

		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		}

		prompt := promptui.Prompt{
			Label:     "Spicy Level",
			Templates: templates,
			Validate:  validate,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You answered %s\n", result)

		fmt.Println("customPrompt called")
	},
}

func init() {
	rootCmd.AddCommand(customPromptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// customPromptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// customPromptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
