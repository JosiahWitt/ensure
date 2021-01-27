module github.com/JosiahWitt/ensure/cmd/ensure

go 1.15

require (
	bursavich.dev/fs-shim v1.0.1
	github.com/JosiahWitt/ensure v0.2.0
	github.com/JosiahWitt/erk v0.5.6
	github.com/golang/mock v1.4.4-0.20201210203420-1fe605df5e5f
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/mod v0.4.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/JosiahWitt/ensure => ../../
