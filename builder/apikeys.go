package builder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

// ListAPIKeys shows the available API keys
func ListAPIKeys(settings *Config) (*apigateway.GetApiKeysOutput, error) {
	svc := apigateway.New(session.New(), &aws.Config{Region: settings.Region})

	params := &apigateway.GetApiKeysInput{}
	resp, err := svc.GetApiKeys(params)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateAPIKey creates a new API key
func CreateAPIKey(name string, description string, enabled bool, apikey string, region *string) (*apigateway.ApiKey, error) {
	svc := apigateway.New(session.New(), &aws.Config{Region: region})

	params := &apigateway.CreateApiKeyInput{
		Description: aws.String(description),
		Enabled:     aws.Bool(enabled),
		Name:        aws.String(name),
	}
	if apikey != "" {
		stagekey := &apigateway.StageKey{
			RestApiId: aws.String(apikey),
			StageName: aws.String("prod"),
		}
		params.StageKeys = append(params.StageKeys, stagekey)
	}
	resp, err := svc.CreateApiKey(params)

	if err != nil {
		return nil, err
	}
	return resp, nil
}
