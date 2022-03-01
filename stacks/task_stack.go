package stacks

import (
	"net/http"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type TaskStackProps struct {
	CommonStackProps
}

func TaskStack(scope constructs.Construct, props TaskStackProps) awscdk.Stack {
	stackName := props.ScopeName.Append("TaskStack")
	stack := awscdk.NewStack(scope, stackName.Get(), &awscdk.StackProps{
		StackName: stackName.Get(),
	})

	taskTableName := stackName.Append("TaskTable")
	awsdynamodb.NewTable(stack, taskTableName.Get(), &awsdynamodb.TableProps{
		TableName:     taskTableName.Get(),
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_RETAIN,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	crudTaskTablePolicy := awsiam.NewPolicy(stack, taskTableName.Append("Role").Get(), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			crudTablePolicyStatement(*stack.Region(), taskTableName.Value),
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
	task := versioned.AddResource(jsii.String("task"), nil)

	// TODO CRUD task

	environment := map[string]*string{
		"LOGLEVEL":       aws.String("debug"),
		"TASKTABLE_NAME": taskTableName.Get(),
	}

	createTaskName := stackName.Append("CreateTask")
	createTaskFunc := createGoFunc(stack, createTaskName, GoFuncProps{
		Path:        "./lambda/create_task/main.go",
		Environment: &environment,
	})
	createTaskFunc.Role().AttachInlinePolicy(crudTaskTablePolicy)
	task.AddMethod(jsii.String(http.MethodPost), awsapigateway.NewLambdaIntegration(createTaskFunc, nil), nil)

	getTaskName := stackName.Append("GetTask")
	getTaskFunc := createGoFunc(stack, getTaskName, GoFuncProps{
		Path:        "./lambda/get_task/main.go",
		Environment: &environment,
	})
	getTaskFunc.Role().AttachInlinePolicy(crudTaskTablePolicy)
	task.AddMethod(jsii.String(http.MethodGet), awsapigateway.NewLambdaIntegration(getTaskFunc, nil), nil)

	updateTaskName := stackName.Append("UpdateTask")
	updateTaskFunc := createGoFunc(stack, updateTaskName, GoFuncProps{
		Path:        "./lambda/update_task/main.go",
		Environment: &environment,
	})
	updateTaskFunc.Role().AttachInlinePolicy(crudTaskTablePolicy)
	task.AddMethod(jsii.String(http.MethodPost), awsapigateway.NewLambdaIntegration(updateTaskFunc, nil), nil)

	return stack
}
