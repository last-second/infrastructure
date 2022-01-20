package main

import (
	"os"

	"github.com/last-second/infrastructure/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
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
	scopeName := &stacks.ScopeName{Value: "LastSecond"}

	logrus.WithFields(logrus.Fields{
		"stage":     stage,
		"scopeName": scopeName.Value,
	}).Info()

	stacks.NewMemoryStack(app, stacks.MemoryStackProps{
		CommonStackProps: stacks.CommonStackProps{
			Version:   "v1",
			Stage:     stage,
			ScopeName: scopeName,
		},
	})

	stacks.NewRecallStack(app, stacks.RecallStackProps{
		CommonStackProps: stacks.CommonStackProps{
			Version:   "v1",
			Stage:     stage,
			ScopeName: scopeName,
		},
	})

	app.Synth(nil)
}
