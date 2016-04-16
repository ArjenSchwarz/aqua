// Copyright © 2016 Arjen Schwarz <developer@arjen.eu>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/ArjenSchwarz/aqua/builder"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

// createapikeyCmd represents the createapikey command
var createapikeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API key",
	Long:  `Creates an API key in the region specified (defaults to us-east-1)`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := builder.CreateAPIKey(keyname, keydescription, keyenabled, apiid, settings.Region)
		if err != nil {
			printFailure(err.Error())
			return
		}
		printSuccess(aws.StringValue(key.Id))
	},
}

var (
	keyname        string
	keydescription string
	keyenabled     bool
	apiid          string
)

func init() {
	apikeyCmd.AddCommand(createapikeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createapikeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createapikeyCmd.Flags().BoolVar(&keyenabled, "enabled", true, "Should the key be enabled")
	createapikeyCmd.Flags().StringVar(&keyname, "keyname", "", "The name for the key")
	createapikeyCmd.Flags().StringVar(&keydescription, "description", "", "The description for the key")
	createapikeyCmd.Flags().StringVar(&apiid, "apiid", "", "The ID of the API you wish to attach the key to")
}
