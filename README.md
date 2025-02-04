# Multi-Account Infrastructure Management

Infrastructure as Code (IaC) repository for managing multiple customer environments across separate AWS accounts using Pulumi.

# Overview

This repository demonstrates a secure and scalable approach to managing infrastructure for multiple customers, each with their own AWS account and environments (dev/prod).

## Repository Structure

```
infrastructure/
├── components/
│   ├── ecs/
│   └── secrets/
└── customers/
    ├── customer1/
    │   │
    │   ├── dev/
    │   │   ├── main.go     # Dev environment (references backend outputs)
    │   │   └── Pulumi.yaml # Uses backend S3/KMS
    │   └── prod/
    │       ├── main.go     # Prod environment (references backend outputs)
    │       └── Pulumi.yaml # Uses backend S3/KMS
    └── customer2/
        ├── dev/
        └── prod/
```

## Setup

Deploy Backend Infrastructure
For each customer AWS account:

```
# Deploy backend infrastructure (S3 bucket and KMS key)
cd customers/customer1/project1/staging
aws s3api create-bucket --bucket customer1-staging-pulumi-bucket --region us-east-1
aws s3api put-bucket-versioning --bucket customer1-staging-pulumi-bucket --versioning-configuration Status=Enabled

aws s3api put-bucket-encryption --bucket customer1-staging-pulumi-bucket --server-side-encryption-configuration '{
  "Rules": [
    {
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "aws:kms"
      }
    }
  ]
}'

pulumi login s3://customer1-staging-pulumi-bucket

pulumi new aws-go \
  --name project1 \
  --description "A minimal AWS Go Pulumi program for customer1" \
  --stack staging \
  --yes

pulumi config set aws:region us-east-1
pulumi stack init organization/project1/staging
pulumi config set env1-name env1-value
pulumi config set --secret dbSecret S3cr37-db-password
pulumi up
pulumi stack output secretArn

cd customers/customer1/project1/prod
aws s3api create-bucket --bucket customer1-prod-pulumi-bucket --region us-east-1
aws s3api put-bucket-versioning --bucket customer1-prod-pulumi-bucket --versioning-configuration Status=Enabled

aws s3api put-bucket-encryption --bucket customer1-prod-pulumi-bucket --server-side-encryption-configuration '{
  "Rules": [
    {
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "aws:kms"
      }
    }
  ]
}'

pulumi login s3://customer1-prod-pulumi-bucket

pulumi new aws-go \
  --name project1 \
  --description "A minimal AWS Go Pulumi program for customer1" \
  --stack prod \
  --yes
..
pulumi up
pulumi stack output secretArn

```

## Best Practices

Each customer would have their own:

- S3 bucket for state
- KMS key for pulumi encryption
- AWS Secrets Manager Secret to store sensitive data
- IAM roles and permissions
- All in their/our own AWS account
- Code reusability in IoC

This separation ensures:

- Complete isolation between customers and thier environments
- Customer-specific access controls
- Clear resource ownership
- Independent state management
- Not combining app code with IoC

## Contributing

External contributions are not open. However, if you have a fix or improvement, feel free to open a pull request. I may review and merge PRs.
