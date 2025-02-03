# Multi-Account Infrastructure Management

## Structure

```
.
├── components/
│   ├── ecs/        # ECS Cluster and Service component
│   └── secrets/    # Secrets management component
├── customers/
│   ├── customer1/  # Customer 1 infrastructure (Account: 111111111111)
│   │   ├── dev/
│   │   └── prod/
│   └── customer2/  # Customer 2 infrastructure (Account: 222222222222)
│       ├── dev/
│       └── prod/
└── .github/
    └── workflows/  # GitHub Actions workflows
```

## Prerequisites

1. AWS Accounts Setup:

```
# Customer1 Account (111111111111)
aws iam create-role --role-name pulumi-deployment-role
aws iam attach-role-policy --role-name pulumi-deployment-role --policy-arn arn:aws:iam::aws:policy/AdministratorAccess

# Customer2 Account (222222222222)
# Repeat the same steps
```

3. Configure GitHub repository:

```
# Add secrets
PULUMI_ACCESS_TOKEN=<your-token>
AWS_OIDC_ROLE_ARN=arn:aws:iam::111111111111:role/pulumi-deployment-role
AWS_OIDC_ROLE_ARN_CUSTOMER2=arn:aws:iam::222222222222:role/pulumi-deployment-role
```

2. IAM Roles:

```
# In each AWS account
Role: pulumi-deployment-role
Policy: AdministratorAccess (limit as needed)
Trust Relationship: GitHub OIDC
```

3. Pulumi Configuration:

```
# Per customer/environment
pulumi config set aws:region us-west-2
pulumi config set --secret aws:assumeRoleARN arn:aws:iam::ACCOUNT_ID:role/pulumi-deployment-role
```

5. Test deployments:

```
# Manual test
cd customers/customer1/dev
pulumi up

# Automated test
git push origin main  # Triggers dev deployment
gh release create v1.0.0  # Triggers prod deployment
```

## Secret Management

1. Store secrets in AWS Secrets Manager:

```
# Format: /customer/environment/service/key
/customer1/dev/api/DB_PASSWORD
/customer1/prod/api/DB_PASSWORD
```

2. Access in IoC code:

```
secrets.GetSecret(ctx, "api/DB_PASSWORD")
```

## Deployment Process

1. Manual Deployment (Day 1: initial deployment phase):

```
# Deploy Customer1 Dev
cd customers/customer1/dev
pulumi up

# Deploy Customer1 Prod
cd customers/customer1/prod
pulumi up
```

2. Automated Deployment (GitHub Actions, Day2: Ongoing operations):

```
Push to main -> Deploy to Dev
Create Release -> Deploy to Prod by Github action
```

## GitHub Actions Setup

1. Configure Secrets in GitHub:

```
PULUMI_ACCESS_TOKEN=pul-xxxx
AWS_OIDC_ROLE_ARN=arn:aws:iam::111111111111:role/pulumi-deployment-role
```

2. Trigger deployments:

```
Push to main branch
Create new Release
```

## Security Best Practices

1. AWS Account Isolation:

   - Each customer has separate AWS account
   - Minimal IAM permissions per account

2. Secret Management:

   - Use AWS Secrets Manager for sensitive data
   - Use SSM Parameter Store for configuration
   - Encrypt secrets at rest and in transit

3. Access Control:
   - GitHub Actions OIDC for AWS authentication
   - No long-lived credentials
   - Separate roles per customer account

## Testing the Setup

1. Clone repository
2. Update AWS account IDs in configs
3. Push to trigger GitHub Actions
4. Monitor deployment in Actions tab
