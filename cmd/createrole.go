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
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
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
`,
	Run: createRole,
}

func init() {
	roleCmd.AddCommand(createCmd)

	// createCmd.Flags().StringVarP(&roleName, "role", "r", "", "The name you wish to give the role")
	createCmd.Flags().StringVarP(&roleType, "type", "t", "basic", "The type of the role")
}

func createRole(cmd *cobra.Command, args []string) {
	roleType = strings.ToLower(roleType)
	var roleTemplate string
	switch roleType {
	case "basic":
		roleTemplate = basicRole
	case "s3":
		roleTemplate = s3Role
	case "aqua":
		roleTemplate = aquaRole
	default:
		printFailure("I'm sorry, but I can't create that role for you.")
		return
	}
	err := doRoleCreation(roleTemplate)
	if err != nil {
		printFailure(err.Error())
		return
	}
	printSuccess(fmt.Sprintf("Role %s of type %s has been created", roleName, roleType))
}

func doRoleCreation(roleTemplate string) error {
	svc := iam.New(session.New())

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustDocument), // Required
		RoleName:                 aws.String(roleName),      // Required
	}
	_, err := svc.CreateRole(params)

	if err != nil {
		return err
	}

	putParams := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(roleTemplate),     // Required
		PolicyName:     aws.String("policyNameType"), // Required
		RoleName:       aws.String(roleName),         // Required
	}
	_, err = svc.PutRolePolicy(putParams)

	if err != nil {
		return err
	}
	return nil
}
