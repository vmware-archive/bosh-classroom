package templates

var DefaultTemplate *Template = &Template{
	AWSTemplateFormatVersion: "2010-09-09",
	Description:              "BOSH 101 Classroom CloudFormation Template",
	Parameters: map[string]Parameter{
		"InstanceType": Parameter{
			Description:           "EC2 Instance Type for Classroom VMs",
			Type:                  "String",
			Default:               "m3.medium",
			ConstraintDescription: "must be a valid EC2 instance type.",
		},

		"AMI": Parameter{
			Description: "AMI to boot",
			Type:        "String",
			ConstraintDescription: "should be an AMI from the latest cloudfoundry/bosh-lite Vagrant box",
		},

		"KeyName": Parameter{
			Description: "The EC2 Key Pair to allow SSH access to the instances",
			Type:        "AWS::EC2::KeyPair::KeyName",
			ConstraintDescription: "must be the name of an existing EC2 KeyPair",
		},

		"InstanceCount": Parameter{
			Description: "Number of EC2 instances to boot for the classroom",
			Type:        "Number",
			Default:     "1",
		},

		"SSHLocation": Parameter{
			Description:           "The IP address range that can be used to SSH to the EC2 instances",
			Type:                  "String",
			Default:               "0.0.0.0/0",
			ConstraintDescription: "must be a valid IP CIDR range of the form x.x.x.x/x.",
		},
	},
	Resources: map[string]Resource{
		"AutoScalingGroup": Resource{
			Type: "AWS::AutoScaling::AutoScalingGroup",
			Properties: map[string]interface{}{
				"AvailabilityZones":       Fn("GetAZs", ""),
				"LaunchConfigurationName": Ref{"LaunchConfig"},
				"MinSize":                 Ref{"InstanceCount"},
				"MaxSize":                 Ref{"InstanceCount"},
			},
		},

		"LaunchConfig": Resource{
			Type: "AWS::AutoScaling::LaunchConfiguration",
			Properties: map[string]interface{}{
				"KeyName":        Ref{"KeyName"},
				"ImageId":        Ref{"AMI"},
				"SecurityGroups": []Ref{{"InstanceSecurityGroup"}},
				"InstanceType":   Ref{"InstanceType"},
				"UserData": Fn("Base64", FnJoin("",
					"#!/bin/bash -xe\n",
					"curl -L -o /home/ubuntu/stemcell.tgz https://s3.amazonaws.com/bosh-warden-stemcells/bosh-stemcell-2776-warden-boshlite-ubuntu-trusty-go_agent.tgz\n",
					"\n")),
			},
		},

		"InstanceSecurityGroup": Resource{
			Type: "AWS::EC2::SecurityGroup",
			Properties: map[string]interface{}{
				"GroupDescription": FnJoin("", "SSH Access for Classroom ", Ref{"AWS::StackName"}),
				"SecurityGroupIngress": []interface{}{
					map[string]interface{}{
						"IpProtocol": "tcp",
						"FromPort":   "22",
						"ToPort":     "22",
						"CidrIp":     Ref{"SSHLocation"},
					},
				},
			},
		},
	},
}
