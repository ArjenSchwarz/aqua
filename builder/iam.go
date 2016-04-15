package builder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// GetRole retrieves the IAM role with the provided name
func GetRole(name *string) (*iam.GetRoleOutput, error) {
	svc := iam.New(session.New())

	params := &iam.GetRoleInput{
		RoleName: name,
	}
	return svc.GetRole(params)
}

// GetRoles returns all the roles the caller has access to
func GetRoles() (*iam.ListRolesOutput, error) {
	svc := iam.New(session.New())

	params := &iam.ListRolesInput{}

	return svc.ListRoles(params)
}

// CreateIAMRole creates an IAM Role based on the provided template
func CreateIAMRole(roleTemplate string, roleName *string) error {
	svc := iam.New(session.New())

	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(TrustDocument),
		RoleName:                 roleName,
	}
	_, err := svc.CreateRole(params)

	if err != nil {
		return err
	}

	putParams := &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(roleTemplate),
		PolicyName:     aws.String("policyNameType"),
		RoleName:       roleName,
	}
	_, err = svc.PutRolePolicy(putParams)

	if err != nil {
		return err
	}
	return nil
}
