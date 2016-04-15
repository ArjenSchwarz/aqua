package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// GatewayBuilder contains the resources needed for creating the Gateway
type GatewayBuilder struct {
	Lambda       *lambda.FunctionConfiguration
	APIGateway   *apigateway.RestApi
	RootResource *apigateway.Resource
	Resource     *apigateway.Resource
	Settings     *Config
}

// APIARN returns the ARN of the API
func (builder *GatewayBuilder) APIARN() string {
	apiArn := strings.Replace(aws.StringValue(builder.Lambda.FunctionArn), "lambda", "execute-api", 1)
	return strings.Replace(apiArn,
		fmt.Sprintf("function:%s", aws.StringValue(builder.Settings.FunctionName)),
		aws.StringValue(builder.APIGateway.Id), 1)
}

// Endpoint returns the endpoint of the API Gateway
func (builder *GatewayBuilder) Endpoint() string {
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/%s",
		aws.StringValue(builder.APIGateway.Id),
		aws.StringValue(builder.Settings.Region),
		builder.Settings.CleanName())
}
