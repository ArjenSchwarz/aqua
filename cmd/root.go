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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/spf13/cobra"
)

var (
	functionName   string
	roleName       string
	region         string
	filePath       string
	authentication string
	jsonOutput     bool
	apikeyRequired bool
	runtime        string
	httpMethod     = "POST"
)

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
	RootCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "The name of the Lambda function")
	RootCmd.PersistentFlags().StringVarP(&roleName, "role", "r", "", "The name of the IAM Role)")
	RootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "The region for the lambda function and API Gateway")
	RootCmd.Flags().StringVarP(&authentication, "authentication", "a", "NONE", "The Authentication method to be used")
	RootCmd.Flags().StringVarP(&filePath, "file", "f", "", "The zip file for your Lambda function, either locally or http(s). The file will first be downloaded locally.")
	RootCmd.Flags().BoolVarP(&apikeyRequired, "apikey", "k", false, "Endpoint can only be accessed with an API key")
	RootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Set to true to print output in JSON format")
	RootCmd.Flags().StringVar(&runtime, "runtime", "nodejs4.3", "The runtime of the Lambda function.")
}

func buildGateway(cmd *cobra.Command, args []string) {
	builder := builder{}
	err := builder.ensureLambdaFunction()

	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.createAPIGateway()

	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.addResources()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.configureResources()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.deployAPI()
	if err != nil {
		printFailure(err.Error())
		return
	}

	err = builder.addPermissions()
	if err != nil {
		printFailure(err.Error())
		return
	}

	msg := fmt.Sprintf("Your endpoint is available at %s", builder.Endpoint())
	if apikeyRequired {
		msg += "\nRemember to configure your API keys before you can use this endpoint."
	}
	printSuccess(msg)
}

type builder struct {
	Lambda       *lambda.FunctionConfiguration
	APIGateway   *apigateway.RestApi
	RootResource *apigateway.Resource
	Resource     *apigateway.Resource
}

func formatName() string {
	return strings.ToLower(functionName)
}

func (*builder) getRole(name string) (*iam.GetRoleOutput, error) {
	svc := iam.New(session.New())

	params := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	return svc.GetRole(params)
}

func (builder *builder) ensureLambdaFunction() error {
	svc := lambda.New(session.New(), &aws.Config{Region: aws.String(region)})

	searchParams := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}
	lambda, err := svc.GetFunctionConfiguration(searchParams)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// If it didn't find the function, we can create it
			if awsErr.Code() == "ResourceNotFoundException" {
				return builder.createLambdaFunction(svc)
			}
		}
		return err
	}
	builder.Lambda = lambda

	return nil
}

func (builder *builder) createLambdaFunction(svc *lambda.Lambda) error {
	if roleName == "" {
		return errors.New("When creating a Lambda function you have to provide a Role for it using the --role flag")
	}
	role, err := builder.getRole(roleName)

	if err != nil {
		return err
	}

	var functionData []byte

	if filePath == "" {
		functionData, err = base64.StdEncoding.DecodeString(helloworld64)
		if err != nil {
			return err
		}
	} else {
		if filePath[0:5] == "http:" || filePath[0:6] == "https:" {
			filePath, err = downloadFile(filePath)
			if err != nil {
				return err
			}
		}
		functionData, err = ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
	}

	params := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: functionData,
		},
		FunctionName: aws.String(functionName),
		Handler:      aws.String("index.handler"),
		Role:         role.Role.Arn,
		Runtime:      aws.String(runtime),
	}
	newLambda, err := svc.CreateFunction(params)

	if err != nil {
		return err
	}

	builder.Lambda = newLambda

	return nil
}

func (builder *builder) createAPIGateway() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.CreateRestApiInput{
		Name:        aws.String(fmt.Sprintf("%sAPI", formatName())),
		Description: aws.String(fmt.Sprintf("API for Lambda function %s", functionName)),
	}
	gateway, err := svc.CreateRestApi(params)
	if err != nil {
		return err
	}

	builder.APIGateway = gateway

	return nil
}

func (builder *builder) addResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.GetResourcesInput{
		RestApiId: builder.APIGateway.Id,
		Limit:     aws.Int64(1),
	}
	resp, err := svc.GetResources(params)

	if err != nil {
		return err
	}

	builder.RootResource = resp.Items[0]

	resourceParams := &apigateway.CreateResourceInput{
		ParentId:  builder.RootResource.Id,
		PathPart:  aws.String(formatName()),
		RestApiId: builder.APIGateway.Id,
	}
	resource, err := svc.CreateResource(resourceParams)

	if err != nil {
		return err
	}

	builder.Resource = resource

	return nil
}

func (builder *builder) configureResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	uriString := fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations", region, aws.StringValue(builder.Lambda.FunctionArn))

	methodParams := &apigateway.PutMethodInput{
		AuthorizationType: aws.String(authentication),
		HttpMethod:        aws.String(httpMethod),
		ResourceId:        builder.Resource.Id,
		RestApiId:         builder.APIGateway.Id,
		ApiKeyRequired:    aws.Bool(apikeyRequired),
	}
	_, err := svc.PutMethod(methodParams)

	if err != nil {
		return err
	}

	params := &apigateway.PutIntegrationInput{
		HttpMethod: aws.String(httpMethod),
		ResourceId: builder.Resource.Id,
		RestApiId:  builder.APIGateway.Id,
		Type:       aws.String("AWS"),
		IntegrationHttpMethod: aws.String(httpMethod),
		RequestTemplates: map[string]*string{
			"application/x-www-form-urlencoded": aws.String(`{"body": $input.json("$")}`),
		},
		Uri: aws.String(uriString),
	}
	_, err = svc.PutIntegration(params)

	if err != nil {
		return err
	}

	integrationResponseParams := &apigateway.PutIntegrationResponseInput{
		HttpMethod:       aws.String(httpMethod),
		ResourceId:       builder.Resource.Id,
		RestApiId:        builder.APIGateway.Id,
		StatusCode:       aws.String("200"),
		SelectionPattern: aws.String(".*"),
	}
	_, err = svc.PutIntegrationResponse(integrationResponseParams)

	if err != nil {
		return err
	}

	methodResponsParams := &apigateway.PutMethodResponseInput{
		HttpMethod:     aws.String(httpMethod),
		ResourceId:     builder.Resource.Id,
		RestApiId:      builder.APIGateway.Id,
		StatusCode:     aws.String("200"),
		ResponseModels: map[string]*string{},
	}
	_, err = svc.PutMethodResponse(methodResponsParams)

	return err
}

func (builder *builder) deployAPI() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.CreateDeploymentInput{
		RestApiId: builder.APIGateway.Id,
		StageName: aws.String("prod"),
	}
	_, err := svc.CreateDeployment(params)

	return err
}

func (builder *builder) addPermissions() error {
	svc := lambda.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: aws.String(functionName),
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String(fmt.Sprintf("apigateway-%s-test", aws.StringValue(builder.Resource.Id))),
		SourceArn:    aws.String(fmt.Sprintf("%s/*/%s/%s", builder.APIARN(), httpMethod, functionName)),
	}
	_, err := svc.AddPermission(params)

	if err != nil {
		return err
	}

	params.SourceArn = aws.String(fmt.Sprintf("%s/prod/%s/%s", builder.APIARN(), httpMethod, functionName))
	params.StatementId = aws.String(fmt.Sprintf("apigateway-%s-prod", aws.StringValue(builder.Resource.Id)))

	_, err = svc.AddPermission(params)

	if err != nil {
		return err
	}

	return err
}

func downloadFile(rawURL string) (string, error) {
	fmt.Println("Downloading file...")

	fileName := os.TempDir() + strconv.FormatInt(time.Now().Unix(), 10) + "aqua.zip"
	file, err := os.Create(fileName)

	if err != nil {
		return "", err
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(rawURL)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

// APIARN returns the ARN of the API
func (builder *builder) APIARN() string {
	apiArn := strings.Replace(aws.StringValue(builder.Lambda.FunctionArn), "lambda", "execute-api", 1)
	return strings.Replace(apiArn, fmt.Sprintf("function:%s", functionName), aws.StringValue(builder.APIGateway.Id), 1)
}

// Endpoint returns the endpoint of the API Gateway
func (builder *builder) Endpoint() string {
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
		aws.StringValue(builder.APIGateway.Id),
		region,
		formatName())
}
