package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	//"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	//"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/micahhausler/rabbit-herd/herd"
	"os"
)

const Version = "0.0.1"

var version = flag.Bool("version", false, "print version and exit")
var pathPtr = flag.StringP("path", "p", "/", "The ec2metadata `path`")
var instanceIdPtr = flag.StringP("instance", "i", "i-1d1bd4bf", "The ec2 instance id to find the group for")

func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Printf("utmptail %s\n", Version)
		os.Exit(0)
	}
	/*
		p := &ec2rolecreds.EC2RoleProvider{}
		v, err := p.Retrieve()
		if err != nil {
			panic(err)
		}
		mdConfig := ec2metadata.New(&ec2metadata.Config{})
		if mdConfig.Available() {
			data, err := mdConfig.GetMetadata(*pathPtr)
			if err != nil {
				panic(err)
			}

			fmt.Println(data)
		}
	*/

	self, err := herd.Self()
	if err != nil {
		panic(err)
	}
	j, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println("Current Node: ")
	fmt.Println(string(j))
	fmt.Println("")

	instances := herd.GetGroupInstances(*instanceIdPtr)
	fmt.Printf("Group for %s: %s\n", *instanceIdPtr, instances)

	/*
		credentialMap := herd.GetCredentials()
		j, err := json.MarshalIndent(credentialMap, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Credentials: %s", string(j))
	*/

}
