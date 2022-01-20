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
	id := props.ScopeName.Append("RecallStack").Append(props.Stage)
	stack := awscdk.NewStack(scope, id.Get(), &awscdk.StackProps{
		StackName: id.Get(),
	})

	restApiName := id.Append("RestApi").Get()
	restApi := awsapigateway.NewRestApi(stack, restApiName, &awsapigateway.RestApiProps{
		RestApiName: restApiName,
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_METHODS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
		},
	})

	api := restApi.Root().AddResource(jsii.String("api"), nil)

	// TODO authorizer

	userEndpoints(stack, api, id, props.Version)

	// TODO taskEndpoints(scope, api)

	return stack
}

func userEndpoints(scope constructs.Construct, api awsapigateway.Resource, id *ScopeName, version string) {
	versioned := api.AddResource(jsii.String(version), nil)
	user := versioned.AddResource(jsii.String("user"), nil)

	// TODO CRUD user

	getUser := createGoFunc(scope, id.Append("GetUser"), GoFuncProps{Path: "./lambda/get_user/main.go"})
	user.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getUser, nil), nil)

}
