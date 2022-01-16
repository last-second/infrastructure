package main

import (
	"os"

	"github.com/last-second/infrastructure/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

func getConfigValue(app awscdk.App, key, defaultValue string) string {
	contextValue := app.Node().TryGetContext(jsii.String(key))

	if value, ok := (contextValue).(string); ok {
		return value
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue
	}

	return ""
}

func main() {
	app := awscdk.NewApp(nil)

	var (
		stage = getConfigValue(app, "AWS_STAGE", "local")
	)

	logrus.Info()

	stacks.NewMemoryStack(app, stacks.MemoryStackProps{
		CommonStackProps: stacks.CommonStackProps{
			Version: "v1",
			Stage:   stage,
		},
	})
	stacks.NewRecallStack(app, stacks.RecallStackProps{
		CommonStackProps: stacks.CommonStackProps{
			Version: "v1",
			Stage:   stage,
		},
	})

	app.Synth(nil)
}
