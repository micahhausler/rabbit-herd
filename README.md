[![Build Status](https://travis-ci.org/micahhausler/rabbit-herder.svg)](https://travis-ci.org/micahhausler/rabbit-herder)
[![https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/micahhausler/rabbit-herder/)

# rabbit-herder
rabbit-herder is a tool for clustering RabbitMQ hosts in an AWS Autoscale group.

There are a number of assumptions and requirements for using this tool:

* All hosts have the same erlang cookie
* You are NOT using the environment variable `RABBITMQ_USE_LONGNAME`
* The autoscale group is in a VPC with DNS resolution and hostnames enabled
* The instance running this tool has an IAM role with permissions to:
	* DescribeAutoscalingInstances (autoscaling)
	* DescribeInstances (ec2)


## Usage

```
Usage of rabbit-herder:
  -a, --api="http://localhost:15672": The rabbitmq API to connect to.
  -d, --dry-run[=false]: Print commands, but don't run them
  -p, --password="guest": The password for the API
  -u, --user="guest": The user account for the API
  -v, --version[=false]: Print version and exit
```

## Build the Docker image
To build the docker images, simply run `make`

## Wishlist

- [ ] Use Consul/etcd2/zookeeper for locks
- [ ] Add tests/CI
- [ ] Automate docker builds

## License
MIT License. See [License](/LICENSE) for full text
