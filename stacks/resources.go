package stacks

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoFuncProps struct {
	// Note: path relative to project root
	Path        string
	Environment *map[string]*string
}

func createGoFunc(
	scope constructs.Construct,
	name *ScopeName,
	props GoFuncProps,
) awslambda.Function {
	code := awslambda.AssetCode_FromAsset(
		// this path should include the vendor folder
		// TODO add `go mod vendor` back in
		jsii.String("../services"),
		&awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image: awslambda.Runtime_GO_1_X().BundlingImage(),
				User:  jsii.String("root"),
				Environment: &map[string]*string{
					"CGO_ENABLED": jsii.String("0"),
					"GOOS":        jsii.String("linux"),
					"GOARCH":      jsii.String("amd64"),
				},
				Command: stringList(
					"bash",
					"-c",
					"go build -o /asset-output/main "+props.Path,
				),
			},
		},
	)

	return awslambda.NewFunction(
		scope,
		name.Get(),
		&awslambda.FunctionProps{
			FunctionName: name.Get(),
			Runtime:      awslambda.Runtime_GO_1_X(),
			Handler:      jsii.String("main"),
			Code:         code,
			Environment:  props.Environment,
		},
	)
}

func crudTablePolicyStatement(region, tableName string) awsiam.PolicyStatement {
	return awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: &[]*string{
			aws.String("dynamodb:*"),
		},
		Resources: &[]*string{
			aws.String(fmt.Sprintf("arn:aws:dynamodb:%s:*:table/%s", region, tableName)),
		},
	})
}
