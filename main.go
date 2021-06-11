package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const usage = `NAME:
	usurp - Temporary AWS role assumption shim

USAGE:
	usage: usurp <role arn> <command>
`

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func abort(status int, message interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	os.Exit(status)
}

func assumeRole(roleArn string) (Credentials, error) {
	user, exists := os.LookupEnv("USER")
	if !exists {
		user = "usurp-user"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return Credentials{}, fmt.Errorf("cant get aws config: %w", err)
	}

	stsClient := sts.NewFromConfig(cfg)
	o, err := stsClient.AssumeRole(
		context.Background(),
		&sts.AssumeRoleInput{
			RoleArn:         aws.String(roleArn),
			RoleSessionName: aws.String(user),
		},
	)
	if err != nil {
		return Credentials{}, fmt.Errorf("cant assume role %s: %w", roleArn, err)
	}

	return Credentials{
			AccessKeyId:     *o.Credentials.AccessKeyId,
			SecretAccessKey: *o.Credentials.SecretAccessKey,
			SessionToken:    *o.Credentials.SessionToken,
		},
		nil
}

func runCommand(creds Credentials, command []string) {
	commandPath, err := exec.LookPath(command[0])
	if err != nil {
		abort(1, err)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)

	err = syscall.Exec(commandPath, command, os.Environ())
	if err != nil {
		abort(1, err)
	}
}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "show program help")
	flag.Parse()

	if help || flag.NArg() < 2 {
		fmt.Println(usage)
		os.Exit(64)
	}

	roleArn := flag.Arg(0)
	command := flag.Args()[1:]

	fmt.Fprintf(os.Stderr, "💅 Assuming role: %s\n", roleArn)

	creds, err := assumeRole(roleArn)
	if err != nil {
		abort(1, err)
	}

	runCommand(creds, command)
}
