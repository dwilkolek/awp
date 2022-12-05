module github.com/tfmcdigital/aws-web-proxy

go 1.18

require (
	github.com/aws/aws-sdk-go-v2 v1.17.1
	github.com/aws/aws-sdk-go-v2/config v1.15.16
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.63.1
	github.com/aws/aws-sdk-go-v2/service/ecs v1.18.13
	github.com/aws/aws-sdk-go-v2/service/ssm v1.30.0
	github.com/elliotchance/sshtunnel v1.3.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/txn2/txeh v1.3.0
	go.uber.org/zap v1.23.0
	golang.org/x/exp v0.0.0-20220927162542-c76eaa363f9d
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.12.11 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.11 // indirect
	github.com/aws/smithy-go v1.13.4 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20191122220453-ac88ee75c92c // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)
