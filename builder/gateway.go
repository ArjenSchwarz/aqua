package builder

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

// CreateAPIGateway creates an API Gateway and attaches it to the GatewayBuilder
func (builder *GatewayBuilder) CreateAPIGateway() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: builder.Settings.Region})

	params := &apigateway.CreateRestApiInput{
		Name: aws.String(fmt.Sprintf("%sLambdaI", builder.Settings.CleanName())),
		Description: aws.String(fmt.Sprintf("API for Lambda function %s",
			aws.StringValue(builder.Settings.FunctionName))),
	}
	gateway, err := svc.CreateRestApi(params)
	if err != nil {
		return err
	}

	builder.APIGateway = gateway

	return nil
}

// AddResources adds a Resource (endpoint) to the API Gateway and attaches it to the GatewayBuilder
func (builder *GatewayBuilder) AddResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: builder.Settings.Region})

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
		PathPart:  aws.String(builder.Settings.CleanName()),
		RestApiId: builder.APIGateway.Id,
	}
	resource, err := svc.CreateResource(resourceParams)

	if err != nil {
		return err
	}

	builder.Resource = resource

	return nil
}

// ConfigureResources configures the Resource in the GatewayBuilder to be set up
// for receiving POST messages and translate them into simple JSON messages
func (builder *GatewayBuilder) ConfigureResources() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: builder.Settings.Region})

	uriString := fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations",
		aws.StringValue(builder.Settings.Region),
		aws.StringValue(builder.Lambda.FunctionArn))

	methodParams := &apigateway.PutMethodInput{
		AuthorizationType: builder.Settings.Authentication,
		HttpMethod:        builder.Settings.HTTPMethod,
		ResourceId:        builder.Resource.Id,
		RestApiId:         builder.APIGateway.Id,
		ApiKeyRequired:    builder.Settings.ApikeyRequired,
	}
	_, err := svc.PutMethod(methodParams)

	if err != nil {
		return err
	}

	params := &apigateway.PutIntegrationInput{
		HttpMethod: builder.Settings.HTTPMethod,
		ResourceId: builder.Resource.Id,
		RestApiId:  builder.APIGateway.Id,
		Type:       aws.String("AWS"),
		IntegrationHttpMethod: builder.Settings.HTTPMethod,
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
		HttpMethod:       builder.Settings.HTTPMethod,
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
		HttpMethod:     builder.Settings.HTTPMethod,
		ResourceId:     builder.Resource.Id,
		RestApiId:      builder.APIGateway.Id,
		StatusCode:     aws.String("200"),
		ResponseModels: map[string]*string{},
	}
	_, err = svc.PutMethodResponse(methodResponsParams)

	return err
}

// DeployAPI deploys the API attached to the GatewayBuilder
func (builder *GatewayBuilder) DeployAPI() error {
	svc := apigateway.New(session.New(), &aws.Config{Region: builder.Settings.Region})

	params := &apigateway.CreateDeploymentInput{
		RestApiId: builder.APIGateway.Id,
		StageName: aws.String("prod"),
	}
	_, err := svc.CreateDeployment(params)

	return err
}
