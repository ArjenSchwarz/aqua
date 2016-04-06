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
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/spf13/cobra"
)

var (
	functionName string
	roleName     string
	region       string
	httpMethod   string
	filePath     string
)

var helloworld64 = "UEsDBBQAAAAIAFqghkjgANcUYwAAAG4AAAAIABwAaW5kZXguanNVVAkAA8veBFfN3gRXdXgLAAEE9QEAAAQUAAAALYxBCsJAEATveUWTU4KyDzDkITnG3dYIZkZ2ZiVB/HsWsW4FRXF7aXYLyyzpyYwRtyLRHyod3xQ/I6o4N+/xaVD5a7ASI5m6dtICqyV8oRFzvpe1ql1anPB7hKumvR+a73AAUEsBAh4DFAAAAAgAWqCGSOAA1xRjAAAAbgAAAAgAGAAAAAAAAQAAAKSBAAAAAGluZGV4LmpzVVQFAAPL3gRXdXgLAAEE9QEAAAQUAAAAUEsFBgAAAAABAAEATgAAAKUAAAAAAA=="

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "aqua",
	Short: "Create a basic gateway for a Lambda function",
	Long: `Running aqua will create a gateway for the provided function.
If the function doesn't exist yet, it will first create it using the provided
file or a basic example that echoes back your parameters.
`,
	Run: buildGateway,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().StringVarP(&functionName, "name", "n", "", "The name of the Lambda function")
	RootCmd.Flags().StringVarP(&roleName, "role", "r", "-", "The name of the Role for the lambda function. (Required when making a new Lambda function)")
	RootCmd.Flags().StringVar(&region, "region", "us-east-1", "The region for the lambda function and API Gateway")
	RootCmd.Flags().StringVarP(&httpMethod, "method", "m", "POST", "The HTTP method for the created endpoint API Gateway")
	RootCmd.Flags().StringVarP(&filePath, "file", "f", "helloworld", "The zip file for your Lambda function")
}

func buildGateway(cmd *cobra.Command, args []string) {
	builder := Builder{}
	err := builder.ensureLambdaFunction()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = builder.createApiGateway()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	builder.addResources()
	err = builder.configureResources()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = builder.deployAPI()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = builder.addPermissions()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Your endpoint is available at %s", builder.Endpoint())
}

type Builder struct {
	Lambda       *lambda.FunctionConfiguration
	ApiGateway   *apigateway.RestApi
	RootResource *apigateway.Resource
	Resource     *apigateway.Resource
}

func FormatName() string {
	return strings.ToLower(functionName)
}

func (*Builder) getRole(name string) (*iam.GetRoleOutput, error) {
	svc := iam.New(session.New())

	params := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	return svc.GetRole(params)
}

func (builder *Builder) ensureLambdaFunction() error {
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

func (builder *Builder) createLambdaFunction(svc *lambda.Lambda) error {
	if roleName == "-" {
		return errors.New("When creating a Lambda function you have to provide a Role for it using the --role flag")
	}
	role, err := builder.getRole(roleName)

	if err != nil {
		return err
	}

	var functionData []byte

	if filePath == "helloworld" {
		functionData, err = base64.StdEncoding.DecodeString(helloworld64)
		if err != nil {
			return err
		}
	} else {
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
		Runtime:      aws.String("nodejs"),
	}
	newLambda, err := svc.CreateFunction(params)

	if err != nil {
		return err
	}

	builder.Lambda = newLambda

	return nil
}

func (builder *Builder) createApiGateway() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.CreateRestApiInput{
		Name:        aws.String(fmt.Sprintf("%sAPI", FormatName())),
		Description: aws.String(fmt.Sprintf("API for Lambda function %s", functionName)),
	}
	gateway, err := svc.CreateRestApi(params)
	if err != nil {
		return err
	}

	builder.ApiGateway = gateway

	return nil
}

func (builder *Builder) addResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.GetResourcesInput{
		RestApiId: builder.ApiGateway.Id,
		Limit:     aws.Int64(1),
	}
	resp, err := svc.GetResources(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	builder.RootResource = resp.Items[0]

	resourceParams := &apigateway.CreateResourceInput{
		ParentId:  builder.RootResource.Id,
		PathPart:  aws.String(FormatName()),
		RestApiId: builder.ApiGateway.Id,
	}
	resource, err := svc.CreateResource(resourceParams)

	if err != nil {
		return err
	}
	builder.Resource = resource
	methodParams := &apigateway.PutMethodInput{
		AuthorizationType: aws.String("NONE"),     // Required
		HttpMethod:        aws.String(httpMethod), // Required
		ResourceId:        builder.Resource.Id,    // Required
		RestApiId:         builder.ApiGateway.Id,  // Required
	}
	_, err = svc.PutMethod(methodParams)

	if err != nil {
		return err
	}
	return nil
}

func (builder *Builder) configureResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	uriString := fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations", region, aws.StringValue(builder.Lambda.FunctionArn))

	params := &apigateway.PutIntegrationInput{
		HttpMethod: aws.String(httpMethod), // Required
		ResourceId: builder.Resource.Id,    // Required
		RestApiId:  builder.ApiGateway.Id,  // Required
		Type:       aws.String("AWS"),      // Required
		IntegrationHttpMethod: aws.String(httpMethod),
		RequestTemplates: map[string]*string{
			"application/x-www-form-urlencoded": aws.String(`{"body": $input.json("$")}`), // Required
		},
		Uri: aws.String(uriString),
	}
	_, err := svc.PutIntegration(params)

	if err != nil {
		return err
	}

	integrationResponseParams := &apigateway.PutIntegrationResponseInput{
		HttpMethod:       aws.String(httpMethod), // Required
		ResourceId:       builder.Resource.Id,    // Required
		RestApiId:        builder.ApiGateway.Id,  // Required
		StatusCode:       aws.String("200"),      // Required
		SelectionPattern: aws.String(".*"),
	}
	_, err = svc.PutIntegrationResponse(integrationResponseParams)

	if err != nil {
		return err
	}

	methodResponsParams := &apigateway.PutMethodResponseInput{
		HttpMethod:     aws.String(httpMethod), // Required
		ResourceId:     builder.Resource.Id,    // Required
		RestApiId:      builder.ApiGateway.Id,  // Required
		StatusCode:     aws.String("200"),      // Required
		ResponseModels: map[string]*string{},
	}
	_, err = svc.PutMethodResponse(methodResponsParams)

	return err
}

func (builder *Builder) deployAPI() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &apigateway.CreateDeploymentInput{
		RestApiId: builder.ApiGateway.Id, // Required
		StageName: aws.String("prod"),    // Required
	}
	_, err := svc.CreateDeployment(params)

	return err
}

func (builder *Builder) addPermissions() error {
	svc := lambda.New(session.New(), &aws.Config{Region: aws.String(region)})

	params := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),                                                 // Required
		FunctionName: aws.String(functionName),                                                            // Required
		Principal:    aws.String("apigateway.amazonaws.com"),                                              // Required
		StatementId:  aws.String(fmt.Sprintf("apigateway-%s-test", aws.StringValue(builder.Resource.Id))), // Required
		SourceArn:    aws.String(fmt.Sprintf("%s/*/%s/%s", builder.ApiARN(), httpMethod, functionName)),
	}
	_, err := svc.AddPermission(params)

	if err != nil {
		return err
	}

	params.SourceArn = aws.String(fmt.Sprintf("%s/prod/%s/%s", builder.ApiARN(), httpMethod, functionName))
	params.StatementId = aws.String(fmt.Sprintf("apigateway-%s-prod", aws.StringValue(builder.Resource.Id)))

	_, err = svc.AddPermission(params)

	if err != nil {
		return err
	}

	return err
}

// ApiARN returns the ARN of the API
func (builder *Builder) ApiARN() string {
	apiArn := strings.Replace(aws.StringValue(builder.Lambda.FunctionArn), "lambda", "execute-api", 1)
	return strings.Replace(apiArn, fmt.Sprintf("function:%s", functionName), aws.StringValue(builder.ApiGateway.Id), 1)
}

func (builder *Builder) Endpoint() string {
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
		aws.StringValue(builder.ApiGateway.Id),
		region,
		FormatName())
}
