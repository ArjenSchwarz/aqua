package builder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var lambdases *lambda.Lambda

func lambdaSession(settings *Config) *lambda.Lambda {
	if lambdases == nil {
		lambdases = lambda.New(session.New(), &aws.Config{Region: settings.Region})
	}
	return lambdases
}

// AddPermissions adds test and production permissions for the Gateway to the Lambda function
func (builder *GatewayBuilder) AddPermissions() error {
	svc := lambdaSession(builder.Settings)

	params := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: builder.Settings.FunctionName,
		Principal:    aws.String("apigateway.amazonaws.com"),
		StatementId:  aws.String(fmt.Sprintf("apigateway-%s-test", aws.StringValue(builder.Resource.Id))),
		SourceArn: aws.String(fmt.Sprintf("%s/*/%s/%s",
			builder.APIARN(),
			aws.StringValue(builder.Settings.HTTPMethod),
			aws.StringValue(builder.Settings.FunctionName))),
	}
	_, err := svc.AddPermission(params)

	if err != nil {
		return err
	}

	params.SourceArn = aws.String(fmt.Sprintf("%s/prod/%s/%s",
		builder.APIARN(),
		aws.StringValue(builder.Settings.HTTPMethod),
		aws.StringValue(builder.Settings.FunctionName)))
	params.StatementId = aws.String(fmt.Sprintf("apigateway-%s-prod",
		aws.StringValue(builder.Resource.Id)))

	_, err = svc.AddPermission(params)

	if err != nil {
		return err
	}

	return err
}

// EnsureLambdaFunction retrieves an existing Lambda function or creates a new one
func (builder *GatewayBuilder) EnsureLambdaFunction() error {
	svc := lambdaSession(builder.Settings)
	searchParams := &lambda.GetFunctionConfigurationInput{
		FunctionName: builder.Settings.FunctionName,
	}
	lambda, err := svc.GetFunctionConfiguration(searchParams)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// If it didn't find the function, we can create it
			if awsErr.Code() == "ResourceNotFoundException" {
				lambda, err = createLambdaFunction(builder.Settings)
				if err == nil {
					builder.Lambda = lambda
				}
				return err
			}
		}
		return err
	}
	builder.Lambda = lambda
	return nil
}

func createLambdaFunction(settings *Config) (*lambda.FunctionConfiguration, error) {
	if aws.StringValue(settings.RoleName) == "" {
		return nil, errors.New("When creating a Lambda function you have to provide a Role for it using the --role flag")
	}
	role, err := GetRole(settings.RoleName)

	if err != nil {
		return nil, err
	}

	var functionData []byte

	if aws.StringValue(settings.FilePath) == "" {
		functionData, err = base64.StdEncoding.DecodeString(Helloworld64)
		if err != nil {
			return nil, err
		}
	} else {
		if settings.IsWebPath() {
			settings.FilePath, err = downloadFile(aws.StringValue(settings.FilePath))
			if err != nil {
				return nil, err
			}
		}
		functionData, err = ioutil.ReadFile(aws.StringValue(settings.FilePath))
		if err != nil {
			return nil, err
		}
	}

	svc := lambdaSession(settings)

	params := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: functionData,
		},
		FunctionName: settings.FunctionName,
		Handler:      aws.String("index.handler"),
		Role:         role.Role.Arn,
		Runtime:      settings.Runtime,
	}
	lambda, err := svc.CreateFunction(params)

	if err != nil {
		return nil, err
	}

	return lambda, nil
}

func downloadFile(rawURL string) (*string, error) {
	fileName := os.TempDir() + strconv.FormatInt(time.Now().Unix(), 10) + "aqua.zip"
	file, err := os.Create(fileName)

	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return nil, err
	}

	return &fileName, nil
}

// CreateSchedule creates a schedule for a Lambda function
func CreateSchedule(settings *Config, schedule string) error {
	svc := lambdaSession(settings)

	searchParams := &lambda.GetFunctionConfigurationInput{
		FunctionName: settings.FunctionName,
	}
	lambdaInst, err := svc.GetFunctionConfiguration(searchParams)

	if err != nil {
		return err
	}
	params := &lambda.AddPermissionInput{
		Action:       aws.String("lambda:InvokeFunction"),
		FunctionName: settings.FunctionName,
		Principal:    aws.String("events.amazonaws.com"),
		StatementId:  aws.String(fmt.Sprintf("scheduler-%s", *settings.FunctionName)),
		SourceArn:    aws.String(createEventARN(lambdaInst)),
	}
	_, err = svc.AddPermission(params)
	if err != nil {
		return err
	}

	eventssvc := cloudwatchevents.New(session.New(), &aws.Config{Region: settings.Region})

	cleanedName := cleanName(schedule)

	putruleparams := &cloudwatchevents.PutRuleInput{
		Name:               aws.String(cleanedName),
		ScheduleExpression: aws.String(schedule),
	}
	_, err = eventssvc.PutRule(putruleparams)
	if err != nil {
		return err
	}

	puttargetparams := &cloudwatchevents.PutTargetsInput{
		Rule: aws.String(cleanedName),
		Targets: []*cloudwatchevents.Target{
			{
				Arn: lambdaInst.FunctionArn,
				Id:  aws.String("1"),
			},
		},
	}
	_, err = eventssvc.PutTargets(puttargetparams)

	if err != nil {
		return err
	}
	return nil
}

func createEventARN(lambdaInst *lambda.FunctionConfiguration) string {
	eventArn := strings.Replace(aws.StringValue(lambdaInst.FunctionArn), "lambda", "events", 1)
	return strings.Replace(eventArn, "function:", "rule/", 1)
}

func cleanName(toClean string) string {
	r, _ := regexp.Compile("[^A-Za-z0-9]+")
	return r.ReplaceAllString(toClean, "")
}
