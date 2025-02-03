package ecs

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ServiceConfig struct {
	ServiceName string
	Image       string
	Port        int
	CPU         int
	Memory      int
	EnvVars     map[string]pulumi.StringInput
	Secrets     map[string]pulumi.StringInput
}

type ECSArgs struct {
	CustomerName string
	Environment  string
	VpcId        string
	SubnetIds    []string
	Services     []ServiceConfig
}

func NewECSCluster(ctx *pulumi.Context, name string, args *ECSArgs) (*ecs.Cluster, error) {
	// Create ECS Cluster
	cluster, err := ecs.NewCluster(ctx, fmt.Sprintf("%s-cluster", name), &ecs.ClusterArgs{
		Name: pulumi.String(fmt.Sprintf("%s-%s-cluster", args.CustomerName, args.Environment)),
		Settings: ecs.ClusterSettingArray{
			&ecs.ClusterSettingArgs{
				Name:  pulumi.String("containerInsights"),
				Value: pulumi.String("enabled"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create Task Execution Role
	executionRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-task-execution-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
            "Version": "2012-10-17",
            "Statement": [{
                "Action": "sts:AssumeRole",
                "Principal": {
                    "Service": "ecs-tasks.amazonaws.com"
                },
                "Effect": "Allow",
                "Sid": ""
            }]
        }`),
	})
	if err != nil {
		return nil, err
	}

	// Attach policies to Task Execution Role
	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-task-execution-policy", name), &iam.RolePolicyAttachmentArgs{
		Role:      executionRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	})
	if err != nil {
		return nil, err
	}

	// Create services
	for _, svc := range args.Services {
		// Create Task Definition
		containerDef := map[string]interface{}{
			"name":  svc.ServiceName,
			"image": svc.Image,
			"portMappings": []map[string]interface{}{
				{
					"containerPort": svc.Port,
					"protocol":      "tcp",
				},
			},
			"environment": []map[string]string{},
			"secrets":     []map[string]string{},
		}

		// Add environment variables
		for key, value := range svc.EnvVars {
			containerDef["environment"] = append(containerDef["environment"].([]map[string]string), map[string]string{
				"name":  key,
				"value": value.ToStringOutput().ApplyT(func(s string) string { return s }).(string),
			})
		}

		// Add secrets
		for key, value := range svc.Secrets {
			containerDef["secrets"] = append(containerDef["secrets"].([]map[string]string), map[string]string{
				"name":      key,
				"valueFrom": value.ToStringOutput().ApplyT(func(s string) string { return s }).(string),
			})
		}

		taskDef, err := ecs.NewTaskDefinition(ctx, fmt.Sprintf("%s-%s", name, svc.ServiceName), &ecs.TaskDefinitionArgs{
			Family:                  pulumi.String(fmt.Sprintf("%s-%s-%s", args.CustomerName, svc.ServiceName, args.Environment)),
			Cpu:                     pulumi.String(fmt.Sprintf("%d", svc.CPU)),
			Memory:                  pulumi.String(fmt.Sprintf("%d", svc.Memory)),
			NetworkMode:             pulumi.String("awsvpc"),
			RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
			ExecutionRoleArn:        executionRole.Arn,
			ContainerDefinitions:    pulumi.String(containerDef),
		})
		if err != nil {
			return nil, err
		}

		// Create ECS Service
		_, err = ecs.NewService(ctx, fmt.Sprintf("%s-%s", name, svc.ServiceName), &ecs.ServiceArgs{
			Cluster:        cluster.Arn,
			DesiredCount:   pulumi.Int(1),
			LaunchType:     pulumi.String("FARGATE"),
			TaskDefinition: taskDef.Arn,
			NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
				AssignPublicIp: pulumi.Bool(true),
				Subnets:        pulumi.ToStringArray(args.SubnetIds),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return cluster, nil
}
