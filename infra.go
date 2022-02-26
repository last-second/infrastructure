package main

import (
	"fmt"
	"os"

	"github.com/last-second/infrastructure/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func getConfigValue(app awscdk.App, contextKey, envKey, defaultValue string) string {
	contextValue := app.Node().TryGetContext(jsii.String(contextKey))

	if value, ok := (contextValue).(string); ok {
		return value
	}

	if value, ok := os.LookupEnv(envKey); ok {
		return value
	}

	return defaultValue
}

func main() {
	app := awscdk.NewApp(nil)

	stage := getConfigValue(app, "stage", "AWS_STAGE", "local")
	scopeName := &stacks.ScopeName{Value: fmt.Sprintf("LastSecond-%s", stage)}

	stacks.UserStack(app, stacks.UserStackProps{
		CommonStackProps: stacks.CommonStackProps{
			Version:   "v1",
			Stage:     stage,
			ScopeName: scopeName,
		},
	})

	app.Synth(nil)
}
