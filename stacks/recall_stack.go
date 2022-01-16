package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RecallStackProps struct {
	CommonStackProps
}

func NewRecallStack(scope constructs.Construct, props RecallStackProps) awscdk.Stack {
	id := getStackId("Recall")
	stackName := getStackName(id, props.CommonStackProps.Stage)
	stack := awscdk.NewStack(scope, &id, &awscdk.StackProps{
		StackName: jsii.String(stackName),
	})

	restApi := awsapigateway.NewRestApi(stack, jsii.String(stackName+"-RestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(stackName + "-RestApi"),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_METHODS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
		},
	})

	api := restApi.Root().AddResource(jsii.String("api"), nil)

	// TODO authorizer

	userEndpoints(stack, api, props)

	// TODO taskEndpoints(scope, api)

	return stack
}

func userEndpoints(scope constructs.Construct, api awsapigateway.Resource, props RecallStackProps) {
	versioned := api.AddResource(jsii.String(props.Version), nil)
	user := versioned.AddResource(jsii.String("user"), nil)

	// TODO CRUD user

	getUser := createGoFunc(scope, "GetUser", GoFuncProps{Path: "./lambda/get_user/main.go"})
	user.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getUser, nil), nil)

}
