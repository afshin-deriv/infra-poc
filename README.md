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
cd customers/customer1/staging
pulumi login s3://customer1-pulumi-bucket
pulumi stack init staging
pulumi config set env1-name env1-value
pulumi config set --secret secret1 S3cr37
pulumi up
```

## Best Practices
Each customer would have their own:
- S3 bucket for state
- KMS key for encryption
- IAM roles and permissions
- All in their/our own AWS account
- Code reusability in IoC

This separation ensures:
- Complete isolation between customers
- Customer-specific access controls
- Clear resource ownership
- Independent state management
- Not combining app code with IoC


## Contributing
External contributions are not open. However, if you have a fix or improvement, feel free to open a pull request. I may review and merge PRs.
