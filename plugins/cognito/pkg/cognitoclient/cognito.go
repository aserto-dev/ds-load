package cognitoclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoClient struct {
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID    string
}

func NewCognitoClient(accessKey, secretKey, userPoolID, region string) (*CognitoClient, error) {
	c := &CognitoClient{}

	// Create a new AWS session with access key and secret key
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		},
	})
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return nil, err
	}

	c.cognitoClient = cognitoidentityprovider.New(sess)
	c.userPoolID = userPoolID

	return c, nil
}

func (c *CognitoClient) ListUsers(ctx context.Context) ([]*cognitoidentityprovider.UserType, error) {
	var paginationToken *string

	users := make([]*cognitoidentityprovider.UserType, 0)

	for {
		listUsersInput := &cognitoidentityprovider.ListUsersInput{
			UserPoolId:      aws.String(c.userPoolID),
			PaginationToken: paginationToken,
			Limit:           aws.Int64(60),
		}

		listUsersOutput, err := c.cognitoClient.ListUsersWithContext(ctx, listUsersInput)
		if err != nil {
			fmt.Println("Failed to list users:", err)
			return nil, err
		}

		users = append(users, listUsersOutput.Users...)

		if listUsersOutput.PaginationToken == nil {
			break
		}

		paginationToken = listUsersOutput.PaginationToken
	}

	return users, nil
}

func (c *CognitoClient) ListGroups() ([]*cognitoidentityprovider.GroupType, error) {
	listGroupsInput := &cognitoidentityprovider.ListGroupsInput{
		UserPoolId: aws.String(c.userPoolID),
	}

	resp, err := c.cognitoClient.ListGroups(listGroupsInput)
	if err != nil {
		fmt.Println("Failed to list groups:", err)
		return nil, err
	}

	return resp.Groups, nil
}

func (c *CognitoClient) GetGroupsForUser(user string) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error) {
	listUsersInGroupInput := &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(user),
	}

	return c.cognitoClient.AdminListGroupsForUser(listUsersInGroupInput)
}

func (c *CognitoClient) GetUserByID(id string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	return c.listUser(id)
}

func (c *CognitoClient) GetUserByEmail(email string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	return c.listUser(email)
}

func (c *CognitoClient) listUser(user string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	adminGetUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(user),
	}

	return c.cognitoClient.AdminGetUser(adminGetUserInput)
}
