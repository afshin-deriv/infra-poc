package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		customerName := "customer1"

		// Create KMS key for state encryption
		key, err := kms.NewKey(ctx, fmt.Sprintf("%s-state-key", customerName), &kms.KeyArgs{
			Description:       pulumi.String(fmt.Sprintf("KMS key for %s Pulumi state", customerName)),
			EnableKeyRotation: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Add alias to KMS key for easier identification
		_, err = kms.NewAlias(ctx, fmt.Sprintf("%s-state-key-alias", customerName), &kms.AliasArgs{
			Name:        pulumi.String(fmt.Sprintf("alias/%s-pulumi-state", customerName)),
			TargetKeyId: key.Id,
		})
		if err != nil {
			return err
		}

		// Create S3 bucket for state storage
		bucket, err := s3.NewBucket(ctx, fmt.Sprintf("%s-state-bucket", customerName), &s3.BucketArgs{
			Bucket: pulumi.String(fmt.Sprintf("pulumi-state-%s-111111111111", customerName)), // Replace with actual account ID
			Versioning: &s3.BucketVersioningArgs{
				Enabled: pulumi.Bool(true),
			},
			ServerSideEncryptionConfiguration: &s3.BucketServerSideEncryptionConfigurationArgs{
				Rule: &s3.BucketServerSideEncryptionConfigurationRuleArgs{
					ApplyServerSideEncryptionByDefault: &s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs{
						KmsMasterKeyId: key.Id,
						SseAlgorithm:   pulumi.String("aws:kms"),
					},
				},
			},
			BlockPublicAccess: &s3.BucketPublicAccessBlockArgs{
				BlockPublicAcls:       pulumi.Bool(true),
				BlockPublicPolicy:     pulumi.Bool(true),
				IgnorePublicAcls:      pulumi.Bool(true),
				RestrictPublicBuckets: pulumi.Bool(true),
			},
		})
		if err != nil {
			return err
		}

		// Export values for use in other stacks
		ctx.Export("stateBucketName", bucket.Id)
		ctx.Export("stateKmsKeyId", key.Id)
		ctx.Export("stateKmsKeyArn", key.Arn)

		return nil
	})
}
