package builder

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

// Config contains all the provided settings
type Config struct {
	FunctionName   *string
	RoleName       *string
	Region         *string
	FilePath       *string
	Authentication *string
	JSONOutput     *bool
	ApikeyRequired *bool
	Runtime        *string
	RoleType       *string
	HTTPMethod     *string
	RoleFilename   *string
	NoGateway      *bool
}

// IsWebPath checks if the provided filepath is a web address
func (config Config) IsWebPath() bool {
	value := aws.StringValue(config.FilePath)
	return value[0:5] == "http:" || value[0:6] == "https:"
}

// CleanName returns a cleaned up version of the FunctionName
func (config Config) CleanName() string {
	value := aws.StringValue(config.FunctionName)
	return strings.ToLower(value)
}
