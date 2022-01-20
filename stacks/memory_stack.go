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
	id := props.ScopeName.Append("MemoryStack").Append(props.Stage)
	stack := awscdk.NewStack(scope, id.Get(), &awscdk.StackProps{
		StackName: id.Get(),
	})

	userTableName := id.Append("UserTable").Get()
	awsdynamodb.NewTable(stack, userTableName, &awsdynamodb.TableProps{
		TableName: userTableName,
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// TODO TaskTable

	return stack
}
