/*
rabbit-herder is a tool for clustering RabbitMQ hosts in an AWS Autoscale group.

There are a number of assumptions and requirements for using this tool:
	* All hosts have the same erlang cookie
	* You are NOT using the environment variable RABBITMQ_USE_LONGNAME
	* The autoscale group is in a VPC with DNS resolution and hostnames enabled
	* The instance running this tool has an IAM role with permissions to:
		* DescribeAutoscalingInstances (autoscaling)
		* DescribeInstances (ec2)

For argument help, run:
	rabbit-herder -h

*/
package main

import (
	"fmt"
	"github.com/micahhausler/rabbit-herder/herd"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

// The binary version
const Version = "0.0.1"

var apiPtr = flag.StringP("api", "a", "http://localhost:15672", "The rabbitmq API to connect to.")
var userPtr = flag.StringP("user", "u", "guest", "The user account for the API")
var passwordPtr = flag.StringP("password", "p", "guest", "The password for the API")
var dryRunP = flag.BoolP("dry-run", "d", false, "Print commands, but don't run them")
var version = flag.BoolP("version", "v", false, "Print version and exit")

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("rabbit-herder %s\n", Version)
		os.Exit(0)
	}

	apiHosts := herd.GetApiHosts(*apiPtr, *userPtr, *passwordPtr)
	if len(apiHosts) > 1 {
		fmt.Printf(
			"API responded len(%d) with hosts: %s\n",
			len(apiHosts),
			strings.Trim(fmt.Sprint(apiHosts), "[]"),
		)
		fmt.Println("Already in a cluster!")
		os.Exit(0)
	}

	ec2Hosts := herd.GetOtherHosts()
	if len(ec2Hosts) == 0 {
		fmt.Println("No hosts to join!")
	} else {
		fmt.Printf("Joining hosts: %s\n", ec2Hosts)
		herd.JoinCluster(ec2Hosts, *dryRunP)
	}
}
