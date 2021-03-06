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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"

	"github.com/ArjenSchwarz/aqua/builder"
)

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Display or create IAM roles",
	Long: `Running aqua role will display the names of all the roles you have access to.

To create an IAM role, please use aqua role create.
`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := builder.GetRoles()

		if err != nil {
			printFailure(err.Error())
			return
		}

		if len(resp.Roles) == 0 {
			printSuccess("No roles have been found.")
			return
		}
		var list = "The following roles have been found:"
		for _, role := range resp.Roles {
			list += fmt.Sprintf("\n* %s", aws.StringValue(role.RoleName))
		}

		printSuccess(list)
	},
}

func init() {
	RootCmd.AddCommand(roleCmd)
}
