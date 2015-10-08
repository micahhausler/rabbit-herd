package herd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func GetCredentials() map[string]string {
	credentials := map[string]string{}
	p := &ec2rolecreds.EC2RoleProvider{}

	v, err := p.Retrieve()
	if err != nil {
		panic(err)
	}
	credentials["AccessKeyID"] = v.AccessKeyID
	credentials["SecretAccessKey"] = v.SecretAccessKey
	credentials["SessionToken"] = v.SessionToken
	return credentials
}

func GetMetadata(path string) (string, error) {
	mdConfig := ec2metadata.New(&ec2metadata.Config{})
	return mdConfig.GetMetadata(path)
}

type Node struct {
	Ip             string            `json:"ip"`
	InstanceId     string            `json:"instance-id"`
	AutoScaleGroup string            `json:"as-group"`
	Tags           map[string]string `json:"tags"`
}

// Gets the autoscale group name for an instance-id. If the instance is not
// part of a group or there is an error, an empty string is returned.
func GetGroupInstances(instanceId string) []string {
	instances := []string{}
	metaClient := ec2metadata.New(&ec2metadata.Config{})
	region, err := metaClient.Region()
	if err != nil {
		region = "us-east-1"
	}
	service := autoscaling.New(&aws.Config{Region: &region})

	resp, err := service.DescribeAutoScalingInstances(
		&autoscaling.DescribeAutoScalingInstancesInput{},
	)
	if len(resp.AutoScalingInstances) == 0 || err != nil {
		return instances
	}
	groupName := groupNameForInstance(instanceId, resp)
	if groupName == "" {
		return instances
	}

	return getInstanceIdsForGroup(groupName, resp)
}

// Gets the group name for the instance
func groupNameForInstance(instanceId string, resp *autoscaling.DescribeAutoScalingInstancesOutput) string {
	for _, instanceDetail := range resp.AutoScalingInstances {
		if instanceId == *instanceDetail.InstanceId {
			return *instanceDetail.AutoScalingGroupName
		}
	}
	return ""
}

func getInstanceIdsForGroup(groupName string, resp *autoscaling.DescribeAutoScalingInstancesOutput) []string {
	instanceIds := []string{}
	for _, instanceDetail := range resp.AutoScalingInstances {
		if groupName == *instanceDetail.AutoScalingGroupName {
			instanceIds = append(instanceIds, *instanceDetail.InstanceId)
		}
	}
	return instanceIds
}

// Get this node
func Self() (Node, error) {
	node := Node{}
	mdConfig := ec2metadata.New(&ec2metadata.Config{})
	ip, err := mdConfig.GetMetadata("local-ipv4")
	if err != nil {
		return node, err
	}
	node.Ip = ip
	id, err := mdConfig.GetMetadata("instance-id")
	if err != nil {
		return node, err
	}
	node.InstanceId = id

	return node, nil
}
