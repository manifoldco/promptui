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

	"github.com/spf13/cobra"
)

/*
➜  cobra git:(master) ✗ go run main.go confirm
continue confirm command?: y
confirm called
 */

// confirmCmd represents the confirm command
var confirmCmd = &cobra.Command{
	Use:   "confirm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := promptui.Prompt{
			Label:     "continue confirm command?",
			IsConfirm: true,
		}

		result, err := prompt.Run()

		if err != nil && err != promptui.ErrAbort {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "y" {
			fmt.Println("confirm called")
		} else {
			fmt.Println("confirm cancel")
		}
	},
}

func init() {
	rootCmd.AddCommand(confirmCmd)

	// Here you will define your flags and confirmuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confirmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// confirmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
