package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
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
	environment := map[string]*string{
		"CGO_ENABLED": jsii.String("0"),
		"GOOS":        jsii.String("linux"),
		"GOARCH":      jsii.String("amd64"),
	}

	if props.Environment != nil {
		for key, value := range *props.Environment {
			environment[key] = value
		}
	}

	code := awslambda.AssetCode_FromAsset(
		// this path should include the vendor folder
		jsii.String("../services"),
		&awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image:       awslambda.Runtime_GO_1_X().BundlingImage(),
				User:        jsii.String("root"),
				Environment: &environment,
				Command: stringList(
					"bash",
					"-c",
					"go build -mod=vendor -o /asset-output/main "+props.Path,
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
		},
	)
}
