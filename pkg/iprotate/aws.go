package iprotate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
)

type ApiEndpoint struct {
	Url          *url.URL
	ApiId        string
	DeploymentId string
	ProxyUrl     string
	Config       aws.Config
	Context      context.Context
}

func CreateApi(awsProfile string, url *url.URL) (*ApiEndpoint, error) {
	endpoint := ApiEndpoint{
		Url:     url,
		Context: context.TODO(),
	}

	cfg, err := loadProfileConfig(endpoint.Context, awsProfile)
	if err != nil {
		return nil, err
	}

	client := apigatewayv2.NewFromConfig(*cfg)

	jsonDoc := fmt.Sprintf(`
	{
		"Name": "go-ip-rotate",
		"Description": "go-ip-rotate",
		"ProtocolType": "HTTP",
		"Target": "%s"
	}
	`, url.String())

	input := &apigatewayv2.CreateApiInput{}
	json.Unmarshal([]byte(jsonDoc), input)

	response, err := client.CreateApi(endpoint.Context, input)

	if err != nil {
		return nil, errors.New("CreateApi: " + err.Error())
	}

	deploymentOutput, err := client.CreateDeployment(endpoint.Context, &apigatewayv2.CreateDeploymentInput{
		ApiId: response.ApiId,
	})

	if err != nil {
		return nil, errors.New("CreateDeployment: " + err.Error())
	}

	endpoint.ApiId = *response.ApiId
	endpoint.ProxyUrl = *response.ApiEndpoint
	endpoint.DeploymentId = *deploymentOutput.DeploymentId

	return &endpoint, nil
}

func (ep *ApiEndpoint) Delete() error {
	cfg, err := loadProfileConfig(ep.Context, "default")
	if err != nil {
		return err
	}

	client := apigatewayv2.NewFromConfig(*cfg)

	client.DeleteApi(ep.Context, &apigatewayv2.DeleteApiInput{
		ApiId: &ep.ApiId,
	})

	return nil
}

func loadProfileConfig(ctx context.Context, profile string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
	)

	if err != nil {
		return nil, errors.New("LoadDefaultConfig: " + err.Error())
	}

	return &cfg, nil
}
