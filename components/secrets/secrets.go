package secrets

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SecretArgs struct {
	CustomerName string
	Environment  string
	ServiceName  string
	SecretKey    string
	SecretValue  pulumi.StringInput
}

// CreateSecret creates a new secret in AWS Secrets Manager
func CreateSecret(ctx *pulumi.Context, args *SecretArgs) (*secretsmanager.Secret, error) {
	secretPath := fmt.Sprintf("/%s/%s/%s/%s",
		args.CustomerName,
		args.Environment,
		args.ServiceName,
		args.SecretKey,
	)

	secret, err := secretsmanager.NewSecret(ctx, fmt.Sprintf("%s-%s-%s",
		args.CustomerName, args.ServiceName, args.SecretKey), &secretsmanager.SecretArgs{
		Name:        pulumi.String(secretPath),
		Description: pulumi.String(fmt.Sprintf("Secret for %s %s", args.CustomerName, args.ServiceName)),
	})
	if err != nil {
		return nil, err
	}

	_, err = secretsmanager.NewSecretVersion(ctx, fmt.Sprintf("%s-%s-%s-value",
		args.CustomerName, args.ServiceName, args.SecretKey), &secretsmanager.SecretVersionArgs{
		SecretId:     secret.ID(),
		SecretString: args.SecretValue,
	})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GetSecret retrieves a secret from AWS Secrets Manager
func GetSecret(ctx *pulumi.Context, customerName string, environment string, serviceName string, key string) (pulumi.StringOutput, error) {
	secretPath := fmt.Sprintf("/%s/%s/%s/%s",
		customerName,
		environment,
		serviceName,
		key,
	)

	secret, err := secretsmanager.LookupSecret(ctx, &secretsmanager.LookupSecretArgs{
		Name: secretPath,
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	version, err := secretsmanager.LookupSecretVersion(ctx, &secretsmanager.LookupSecretVersionArgs{
		SecretId: secret.Id,
	})
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	return pulumi.String(version.SecretString).ToStringOutput(), nil
}
