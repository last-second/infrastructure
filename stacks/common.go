package stacks

import (
	"github.com/aws/jsii-runtime-go"
)

type CommonStackProps struct {
	Version, Stage string
}

func getStackId(name string) string {
	return "LastSecond-" + name + "Stack"
}

func getStackName(id, stage string) string {
	return id + "-" + stage
}

func stringList(values ...string) *[]*string {
	result := []*string{}
	for _, value := range values {
		result = append(result, jsii.String(value))
	}

	return &result
}
