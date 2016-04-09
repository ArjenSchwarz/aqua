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

import "github.com/spf13/cobra"

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Aqua as a Lambda function",
	Long: `Use this command to install Aqua as a Lambda function.

It will download the latest version from GitHub, and install it as a Lambda
function with the name you specify.

Due to potential security risks, Aqua will always be installed with API keys
being required.

Before you do this, you first have to ensure you have an IAM role with all the
required permissions. You can use the following command to create such a role:

aqua role create --role role_name --type aqua

Example:
aqua install --name aqua -role aquarole
`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath = "https://github.com/ArjenSchwarz/aqua/releases/download/latest/aqua.zip"
		apikeyRequired = true
		buildGateway(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
