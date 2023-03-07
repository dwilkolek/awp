package domain

import (
	"strings"
)

type Environment string

const (
	DEV  Environment = "dev"
	DEMO Environment = "demo"
	PROD Environment = "prod"
	LAB  Environment = "lab"
)

func (env Environment) String() string {
	return string(env)
}

func ParseEnvironment(s string) (e Environment, err error) {
	environments := map[string]Environment{
		LAB.String():  LAB,
		DEV.String():  DEV,
		DEMO.String(): DEMO,
		PROD.String(): PROD,
	}

	e = environments[strings.ToLower(s)]

	return e, nil
}
