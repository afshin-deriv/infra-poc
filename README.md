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
    │   ├── backend/        # Customer1's backend infrastructure
    │   │   ├── main.go     # S3/KMS setup
    │   │   └── Pulumi.yaml
    │   ├── dev/
    │   │   ├── main.go     # Dev environment (references backend outputs)
    │   │   └── Pulumi.yaml # Uses backend S3/KMS
    │   └── prod/
    │       ├── main.go     # Prod environment (references backend outputs)
    │       └── Pulumi.yaml # Uses backend S3/KMS
    └── customer2/
        ├── backend/        # Customer2's backend infrastructure
        │   ├── main.go     # Different S3/KMS in different account
        │   └── Pulumi.yaml
        ├── dev/
        └── prod/
```

## Setup
1. Deploy Backend Infrastructure
For each customer AWS account:
```
# Deploy backend infrastructure (S3 bucket and KMS key)
customers/customer1/backend
pulumi login s3://customer1-bucket-name
pulumi stack init backend-customer1
pulumi config set --secret aws:kmsKeyId "arn:aws:kms:region:account-id:key/key-id"
pulumi stack ls
pulumi up
```

2. Update Stack Configuration
Update each environment's Pulumi.yaml:
```yaml
name: customer1-dev
runtime: go
description: Customer1 Development Infrastructure
backend:
  url: s3://<bucket-name>/customer1/dev
encryption:
  provider: aws-kms
  key-id: <kms-key-arn>
```

3. Initialize Stacks
For each environment:
```bash
cd customers/customer1/dev
pulumi stack init customer1-dev --secrets-provider="aws-kms"
pulumi config set aws:region us-west-2
```

4. Required IAM Permissions
The deployment role needs these permissions:
- S3 access for state management
- KMS access for encryption
- Permissions to deploy resources
_See iam/deployment-role-policy.json for the complete policy._

## Deployment Workflows
1. First, deploy backend infrastructure:
```
# Deploy backend for Customer1
cd customers/customer1/backend
pulumi stack init backend
pulumi up

# Get outputs for environment configuration
export BACKEND_BUCKET=$(pulumi stack output stateBucketName)
export BACKEND_KMS_KEY=$(pulumi stack output stateKmsKeyArn)

# Configure dev environment
cd ../dev
pulumi stack init dev --secrets-provider="aws-kms"
pulumi config set aws:region us-west-2
# Set backend configuration using outputs from backend stack
```
2. Then deploy environment infrastructure:
```
# Deploy dev environment
cd customers/customer1/dev
pulumi up
```

## Security Best Practices
Each customer would have their own:
- S3 bucket for state
- KMS key for encryption
- IAM roles and permissions
- All in their own AWS account

- State files are stored in your AWS account
- Encryption is handled by AWS KMS
- No Pulumi login required
- Access controlled via IAM

This separation ensures:
- Complete isolation between customers
- Customer-specific access controls
- Clear resource ownership
- Independent state management


## Contributing
External contributions are not open. However, if you have a fix or improvement, feel free to open a pull request. I may review and merge PRs.
