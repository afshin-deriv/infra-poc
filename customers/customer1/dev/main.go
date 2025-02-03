package main

import (
	"infrastructure/components/ecs"
	"infrastructure/components/secrets"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create secrets first
		dbPassword, err := secrets.CreateSecret(ctx, &secrets.SecretArgs{
			CustomerName: "customer1",
			Environment:  "dev",
			ServiceName:  "api",
			SecretKey:    "DB_PASSWORD",
			SecretValue:  pulumi.String("dev-password-123"), // In real world, use Pulumi config
		})
		if err != nil {
			return err
		}

		apiKey, err := secrets.CreateSecret(ctx, &secrets.SecretArgs{
			CustomerName: "customer1",
			Environment:  "dev",
			ServiceName:  "api",
			SecretKey:    "API_KEY",
			SecretValue:  pulumi.String("dev-api-key-123"),
		})
		if err != nil {
			return err
		}

		// Create ECS Infrastructure
		_, err = ecs.NewECSCluster(ctx, "customer1", &ecs.ECSArgs{
			CustomerName: "customer1",
			Environment:  "dev",
			VpcId:        "vpc-12345", // Use Pulumi config in real world
			SubnetIds:    []string{"subnet-1", "subnet-2"},
			Services: []ecs.ServiceConfig{
				{
					ServiceName: "api",
					Image:       "111111111111.dkr.ecr.us-west-2.amazonaws.com/customer1/api:latest",
					Port:        8080,
					CPU:         256,
					Memory:      512,
					EnvVars: map[string]pulumi.StringInput{
						"ENVIRONMENT": pulumi.String("dev"),
						"LOG_LEVEL":   pulumi.String("debug"),
					},
					Secrets: map[string]pulumi.StringInput{
						"DB_PASSWORD": dbPassword.Arn,
						"API_KEY":     apiKey.Arn,
					},
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
