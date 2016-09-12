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
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"

	"github.com/ArjenSchwarz/aqua/builder"
)

var (
	roleType string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a role",
	Long: `Using aqua create role you can create a new IAM role.

You can only choose from the following types:
* Basic (Basic Lambda execution)
* S3 (S3 support)
* Aqua (Everything required for running Aqua)

If you require a different role, please create it manually.

Example: aqua role create --name basic_execution_role --type basic
`,
	Run: createRole,
}

func init() {
	roleCmd.AddCommand(createCmd)

	settings.RoleType = createCmd.Flags().StringP("type", "t", "basic", "The type of the role")
	settings.RoleFilename = createCmd.Flags().String("filename", "", "The role description you want to add")
}

func createRole(cmd *cobra.Command, args []string) {
	roleType = strings.ToLower(aws.StringValue(settings.RoleType))
	var roleTemplate string
	switch roleType {
	case "basic":
		roleTemplate = builder.BasicRole
	case "s3":
		roleTemplate = builder.S3Role
	case "aqua":
		roleTemplate = builder.AquaRole
	case "custom":
		if _, err := os.Stat(*settings.RoleFilename); err != nil {
			printFailure(err.Error())
			return
		}
		roleContents, err := ioutil.ReadFile(*settings.RoleFilename)
		if err != nil {
			printFailure(err.Error())
			return
		}
		roleTemplate = string(roleContents)
	default:
		printFailure("I'm sorry, but I can't create that role for you.")
		return
	}
	err := builder.CreateIAMRole(roleTemplate, settings.RoleName)
	if err != nil {
		printFailure(err.Error())
		return
	}
	printSuccess(fmt.Sprintf("Role %s of type %s has been created",
		aws.StringValue(settings.RoleName),
		aws.StringValue(settings.RoleType)))
}
