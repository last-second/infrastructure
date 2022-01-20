package stacks

import (
	"github.com/aws/jsii-runtime-go"
)

type ScopeName struct {
	Value string
}

func (s *ScopeName) Append(value string) *ScopeName {
	return &ScopeName{Value: s.Value + "-" + value}
}

func (s *ScopeName) Get() *string {
	return jsii.String(s.Value)
}

type CommonStackProps struct {
	Version, Stage string
	*ScopeName
}

func stringList(values ...string) *[]*string {
	result := []*string{}
	for _, value := range values {
		result = append(result, jsii.String(value))
	}

	return &result
}
