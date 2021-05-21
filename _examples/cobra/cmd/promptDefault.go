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
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os/user"

	"github.com/spf13/cobra"
)

/*
➜  cobra git:(master) ✗ go run main.go promptDefault
Username: promptui-name
Your username is "promptui-name"
promptDefault called
 */

// promptDefaultCmd represents the promptDefault command
var promptDefaultCmd = &cobra.Command{
	Use:   "promptDefault",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		validate := func(input string) error {
			if len(input) < 3 {
				return errors.New("Username must have more than 3 characters")
			}
			return nil
		}

		var username string
		u, err := user.Current()
		if err == nil {
			username = u.Username
		}

		prompt := promptui.Prompt{
			Label:    "Username",
			Validate: validate,
			Default:  username,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("Your username is %q\n", result)

		fmt.Println("promptDefault called")
	},
}

func init() {
	rootCmd.AddCommand(promptDefaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promptDefaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promptDefaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
