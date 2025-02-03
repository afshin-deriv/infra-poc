# Multi-Account Infrastructure Management
Infrastructure as Code (IaC) repository for managing multiple customer environments across separate AWS accounts using Pulumi.

# Overview
This repository demonstrates a secure and scalable approach to managing infrastructure for multiple customers, each with their own AWS account and environments (dev/prod).

## Repository Structure
```
├── components/             # Reusable infrastructure components
│   ├── ecs/                # ECS Cluster and Service definitions
│   └── secrets/            # Secrets management utilities
├── customers/              # Customer-specific implementations
│   ├── customer1/          # Customer 1 (AWS Account: 111111111111)
│   │   ├── dev/            # Development environment
│   │   └── prod/           # Production environment
│   └── customer2/          # Customer 2 (AWS Account: 222222222222)
│       ├── dev/
│       └── prod/
└── .github/
    └── workflows/          # Automated deployment workflows
```

## Setup Guide

### 1. AWS Account Setup
Each customer requires their own AWS account with proper IAM configuration:
```bash
# For each customer AWS account (e.g., 111111111111, 222222222222)
aws iam create-role --role-name pulumi-deployment-role
aws iam attach-role-policy \
    --role-name pulumi-deployment-role \
    --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
# Configure OIDC trust relationship for GitHub Actions
# (Add through AWS Console or CLI - see documentation for OIDC setup)
```

### 2. GitHub Repository Configuration
Add these secrets to your GitHub repository:
```
PULUMI_ACCESS_TOKEN=pul-xxxx                                                       # Pulumi access token
AWS_OIDC_ROLE_ARN=arn:aws:iam::111111111111:role/pulumi-deployment-role            # Customer 1
AWS_OIDC_ROLE_ARN_CUSTOMER2=arn:aws:iam::222222222222:role/pulumi-deployment-role  # Customer 2
```

### 3. Pulumi Stack Configuration
For each customer environment:
```bash
# Initialize stack
cd customers/customer1/dev
pulumi stack init dev

# Configure AWS credentials
pulumi config set aws:region us-west-2
pulumi config set --secret aws:assumeRoleARN \
    arn:aws:iam::ACCOUNT_ID:role/pulumi-deployment-role
```

## Secret Management

### Secrets Structure
Secrets are stored hierarchically in AWS Secrets Manager:
```bash
/customer/environment/service/key
Example: /customer1/dev/api/DB_PASSWORD
```

### Accessing Secrets
In your Pulumi code:
```go
// Retrieve a secret
secretValue, err := secrets.GetSecret(ctx, "api/DB_PASSWORD")

// Create a new secret
_, err := secrets.CreateSecret(ctx, &secrets.SecretArgs{
    CustomerName: "customer1",
    Environment: "dev",
    ServiceName: "api",
    SecretKey:   "DB_PASSWORD",
    SecretValue: pulumi.String("mypassword"),
})
```

## Deployment Workflows

### Initial Deployment (Day 1)
Manual deployment for initial setup:
```bash
# Deploy development environment
cd customers/customer1/dev
pulumi up

# Deploy production environment
cd customers/customer1/prod
pulumi up
```

### Ongoing Deployments (Day 2)
Automated via GitHub Actions:
- Push to main branch → Deploys to development
- Create GitHub Release → Deploys to production

## Security Best Practices
1. AWS Account Isolation
- Separate AWS account per customer
- Environment segregation within accounts
- Least privilege IAM policies

2. Secrets Management
- AWS Secrets Manager for sensitive data
- SSM Parameter Store for configuration
- All secrets encrypted at rest and in transit

3. Access Control
- GitHub Actions OIDC for AWS authentication
- No long-lived AWS credentials
- Role-based access with minimal permissions

## Testing Instructions
1. Clone the repository
2. Update AWS account IDs in configuration files:
  - GitHub Actions workflows
  - Pulumi stack configurations
3. Push changes to trigger development deployment
4. Create a release to trigger production deployment
5. Monitor deployments in GitHub Actions tab

## Troubleshooting
Common issues and solutions:
  - OIDC authentication failures: Verify trust relationship configuration
  - Pulumi access errors: Check role permissions and stack configurations
  - Secret access denied: Verify IAM permissions for Secrets Manager

## Contributing
External contributions are not open. However, if you have a fix or improvement, feel free to open a pull request. I may review and merge PRs.
