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

		// Export values for use in CI/CD
		ctx.Export("secretArn", secret.Arn)
		return nil
	})
}
