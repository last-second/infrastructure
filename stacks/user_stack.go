package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type UserStackProps struct {
	CommonStackProps
}

func UserStack(scope constructs.Construct, props UserStackProps) awscdk.Stack {
	stackName := props.ScopeName.Append("UserStack").Append(props.Stage)
	stack := awscdk.NewStack(scope, stackName.Get(), &awscdk.StackProps{
		StackName: stackName.Get(),
	})

	userTableName := stackName.Append("UserTable")
	awsdynamodb.NewTable(stack, userTableName.Get(), &awsdynamodb.TableProps{
		TableName:   userTableName.Get(),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	restApiName := stackName.Append("RestApi")
	restApi := awsapigateway.NewRestApi(stack, restApiName.Get(), &awsapigateway.RestApiProps{
		RestApiName: restApiName.Get(),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowOrigins: awsapigateway.Cors_ALL_METHODS(),
			AllowMethods: awsapigateway.Cors_ALL_METHODS(),
			AllowHeaders: awsapigateway.Cors_DEFAULT_HEADERS(),
		},
	})

	api := restApi.Root().AddResource(jsii.String("api"), nil)

	// TODO authorizer

	versioned := api.AddResource(jsii.String(props.CommonStackProps.Version), nil)
	user := versioned.AddResource(jsii.String("user"), nil)

	// TODO CRUD user

	getUserName := stackName.Append("GetUser")
	getUser := createGoFunc(stack, getUserName, GoFuncProps{
		Path: "./lambda/get_user/main.go",
		Environment: &map[string]*string{
			"LOGLEVEL":       aws.String("debug"),
			"USERTABLE_NAME": userTableName.Get(),
		},
	})

	crudUserTablePolicy := awsiam.NewPolicy(stack, getUserName.Append("Role").Get(), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			crudTablePolicyStatement(*stack.Region(), userTableName.Value),
		},
	})
	getUser.Role().AttachInlinePolicy(crudUserTablePolicy)

	user.AddMethod(jsii.String("GET"), awsapigateway.NewLambdaIntegration(getUser, nil), nil)

	return stack
}
