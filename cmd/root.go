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
	"os"

	"github.com/ArjenSchwarz/aqua/builder"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/spf13/cobra"
)

var (
	httpMethod = "POST"
)

var settings = new(builder.Config)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "aqua",
	Short: "Create a basic gateway for a Lambda function",
	Long: `Create a gateway for the provided Lambda function.

If the function doesn't exist yet, it will first create it using the provided
file or a basic example that echoes back your parameters.

For function code located online, the file will first be downloaded locally.

Example (only create Gateway):
aqua --name functionName --region us-west-1

Example (create Lambda function from local file):
aqua --name functionName --role basic_execution_role --file path/to/function.zip

Example (create Lambda function from web file):
aqua --name functionName --role basic_execution_role --file https://github.com/ArjenSchwarz/aqua/releases/download/latest/igor.zip
`,
	Run: buildGateway,
}

// Execute is the main execution command as created by Cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		printFailure(err.Error())
		os.Exit(-1)
	}
}

func init() {
	settings.FunctionName = RootCmd.PersistentFlags().StringP("name", "n", "", "The name of the Lambda function")
	settings.RoleName = RootCmd.PersistentFlags().StringP("role", "r", "", "The name of the IAM Role")
	settings.Region = RootCmd.PersistentFlags().String("region", "us-east-1", "The region for the lambda function and API Gateway")
	settings.Authentication = RootCmd.Flags().StringP("authentication", "a", "NONE", "The Authentication method to be used")
	settings.FilePath = RootCmd.PersistentFlags().StringP("file", "f", "", "The zip file for your Lambda function, either locally or http(s). The file will first be downloaded locally.")
	settings.ApikeyRequired = RootCmd.Flags().BoolP("apikey", "k", false, "Endpoint can only be accessed with an API key")
	settings.JSONOutput = RootCmd.PersistentFlags().Bool("json", false, "Set to true to print output in JSON format")
	settings.Runtime = RootCmd.Flags().String("runtime", "nodejs4.3", "The runtime of the Lambda function.")
	settings.NoGateway = RootCmd.Flags().Bool("nogateway", false, "Disable the creation of a Gateway. Only create the Lambda function.")
}

func buildGateway(cmd *cobra.Command, args []string) {
	settings.HTTPMethod = aws.String("POST") //For now we force the HTTP method to POST
	builder := builder.GatewayBuilder{Settings: settings}
	err := builder.EnsureLambdaFunction()

	if err != nil {
		printFailure(err.Error())
		return
	}

	if *settings.NoGateway {
		return
	}

	err = builder.CreateAPIGateway()

	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.AddResources()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.ConfigureResources()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.DeployAPI()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.AddPermissions()
	if err != nil {
		printFailure(err.Error())
		return
	}

	messages := make(map[string]string)
	messages["endpoint"] = builder.Endpoint()
	messages["api"] = aws.StringValue(builder.APIGateway.Id)
	if aws.BoolValue(settings.ApikeyRequired) {
		messages["note"] = "Remember to configure your API keys before you can use this endpoint."
	}
	printMap(messages)
}
