package domain

import "fmt"

type Environment string

const (
	DEV  Environment = "dev"
	DEMO Environment = "demo"
	PROD Environment = "prod"
)

func (env Environment) String() string {
	return string(env)
}

func ParseEnvironment(s string) (e Environment, err error) {
	environments := map[Environment]struct{}{
		DEV:  {},
		DEMO: {},
		PROD: {},
	}

	env := Environment(s)
	_, ok := environments[env]
	if !ok {
		return e, fmt.Errorf(`cannot parse:[%s] as environment`, s)
	}
	return e, nil
}
