package main

import (
	ec2comp "github.com/afshin-deriv/infra-poc/components/ec2"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Configure the parameters for the EC2 instance.
		instanceArgs := &ec2comp.InstanceArgs{
			InstanceType: pulumi.String("t2.micro"),
			AmiId:        pulumi.String("ami-04b4f1a9cf54c11d0"),
			// Replace with a valid AMI ID for your region.
			SubnetId: pulumi.String("subnet-f8e7bad6"),
			// Replace with a valid Subnet ID.
			SecurityGroupIds: pulumi.StringArray{pulumi.String("sg-0224a2f2333a31820")},
			// Replace with valid Security Group IDs.
			KeyName: pulumi.StringPtr("access-server"),
			// Optional: Replace with your SSH key name.
			UserData: pulumi.StringPtr("#!/bin/bash\necho Hello, World > /var/tmp/hello.txt"),
		}

		// Create the EC2 instance.
		instance, err := ec2comp.NewInstance(ctx, "afshin-test-ec2-instance", instanceArgs)
		if err != nil {
			return err
		}
		// Export useful outputs.
		ctx.Export("instanceId", instance.Instance.ID())
		ctx.Export("publicIp", instance.Instance.PublicIp)

		return nil
	})
}
