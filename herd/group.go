package herd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

// Gets the region from the metadata service
// fallback parameter specifies the region to use if there is an error
func getMetadataRegion(fallback string) string {
	metaClient := ec2metadata.New(&ec2metadata.Config{})
	region, err := metaClient.Region()
	if err != nil {
		return fallback
	}
	return region
}

// Get a default AWS config
func defaultAwsConfig() aws.Config {
	region := getMetadataRegion("us-east-1")
	return aws.Config{Region: &region}
}

// Get instance id from the metadata service
func getMyInstanceId() string {
	mdConfig := ec2metadata.New(&ec2metadata.Config{})
	id, err := mdConfig.GetMetadata("instance-id")
	if err != nil {
		return ""
	}
	return id
}

// Gets the other hostnames in the autoscale group
func GetOtherHosts() (hostnames []string) {
	instances := []string{}
	config := defaultAwsConfig()
	service := autoscaling.New(&config)

	resp, err := service.DescribeAutoScalingInstances(
		&autoscaling.DescribeAutoScalingInstancesInput{},
	)
	if len(resp.AutoScalingInstances) == 0 || err != nil {
		return instances
	}
	instanceId := getMyInstanceId()
	groupName := groupNameForInstance(instanceId, resp)
	if groupName == "" {
		return instances
	}

	otherIds := otherInstanceIds(
		instanceId,
		getInstanceIdsForGroup(groupName, resp),
	)
	return formatDnsNames(getInstanceDnsNames(otherIds))
}

// Strips the suffixes off of the other hostnames
func formatDnsNames(instanceNames []string) []string {
	strippedNames := []string{}
	for _, name := range instanceNames {
		parts := strings.Split(name, ".")
		if len(parts) >= 1 {
			strippedNames = append(strippedNames, parts[0])
		}
	}
	return strippedNames
}

// Get private dns names for a given list of instance ids
// returns an empty list on error
func getInstanceDnsNames(instanceIds []string) []string {
	dnsNames := []string{}

	ids := []*string{}
	for _, id := range instanceIds {
		ids = append(ids, aws.String(id))
	}

	config := defaultAwsConfig()
	service := ec2.New(&config)
	params := &ec2.DescribeInstancesInput{InstanceIds: ids}
	resp, err := service.DescribeInstances(params)
	if err != nil {
		return dnsNames
	}

	for _, resv := range resp.Reservations {
		for _, instance := range resv.Instances {
			dnsNames = append(dnsNames, *instance.PrivateDnsName)
		}
	}
	return dnsNames
}

// Get other instance ids
func otherInstanceIds(myId string, allIds []string) []string {
	others := []string{}
	for _, id := range allIds {
		if id != myId {
			others = append(others, id)
		}
	}
	return others
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
