package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MemoryStackProps struct {
	CommonStackProps
}

func NewMemoryStack(scope constructs.Construct, props MemoryStackProps) awscdk.Stack {
	id := getStackId("Memory")
	stackName := getStackName(id, props.CommonStackProps.Stage)
	stack := awscdk.NewStack(scope, &id, &awscdk.StackProps{
		StackName: jsii.String(stackName),
	})

	awsdynamodb.NewTable(stack, jsii.String("UserTable"), &awsdynamodb.TableProps{
		TableName: jsii.String("UserTable"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// TODO TaskTable

	return stack
}
