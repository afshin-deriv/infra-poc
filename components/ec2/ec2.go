package ec2

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// InstanceArgs defines the input parameters for creating an EC2 instance.
type InstanceArgs struct {
	// The instance type (e.g., "t2.micro").
	InstanceType pulumi.StringInput
	// The Amazon Machine Image (AMI) ID.
	AmiId pulumi.StringInput
	// The subnet in which to launch the instance.
	SubnetId pulumi.StringInput
	// The security group IDs to associate with the instance.
	SecurityGroupIds pulumi.StringArrayInput
	// Optional key name for SSH access.
	KeyName pulumi.StringPtrInput
	// Optional user data script.
	UserData pulumi.StringPtrInput
}

// Instance is a wrapper around the underlying EC2 instance resource.
type Instance struct {
	pulumi.ResourceState
	// Underlying EC2 instance.
	Instance *ec2.Instance
}

// NewInstance creates a new EC2 instance using the provided configuration.
func NewInstance(ctx *pulumi.Context, name string, args *InstanceArgs, opts ...pulumi.ResourceOption) (*Instance, error) {
	ec2Instance, err := ec2.NewInstance(ctx, name, &ec2.InstanceArgs{
		InstanceType:        args.InstanceType,
		Ami:                 args.AmiId,
		SubnetId:            args.SubnetId,
		VpcSecurityGroupIds: args.SecurityGroupIds,
		KeyName:             args.KeyName,
		UserData:            args.UserData,
		Tags: pulumi.StringMap{
			"Name": pulumi.String(name),
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &Instance{
		Instance: ec2Instance,
	}, nil
}
