package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type InstanceInfo struct {
	ID    string
	Name  string
	State string
}

func GetEC2Instances(profile string) ([]InstanceInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	ec2Client := ec2.NewFromConfig(cfg)

	result, err := ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []InstanceInfo
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceName := "(no name)"
			instanceState := string(instance.State.Name)

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					instanceName = *tag.Value
					break
				}
			}

			instances = append(instances, InstanceInfo{
				ID:    *instance.InstanceId,
				Name:  instanceName,
				State: instanceState,
			})
		}
	}
	return instances, nil
}
