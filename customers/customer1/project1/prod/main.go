package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		secret, err := secretsmanager.NewSecret(ctx, "dbSecret", &secretsmanager.SecretArgs{})
		if err != nil {
			return err
		}

		// ARN of the AWS Secrets Manager secret
		// to be used in CI/CD and App code
		ctx.Export("secretArn", secret.Arn)
		return nil
	})
}
