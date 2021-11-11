## Kubernetes Cluster Autoscaler

[![Go Report Card](https://goreportcard.com/badge/github.com/Chathuru/kubernetes-cluster-autoscaler)](https://goreportcard.com/report/github.com/Chathuru/kubernetes-cluster-autoscaler)

### Developer Requirements
- [Kubernetes](https://kubernetes.io/) version 1.19.x
- [Go](https://golang.org/doc/install) version 1.15.x (to build the source and develop plugins)

### Build docker image
- compile: `CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/autoscaler cmd/main.go`
- build docker image: `docker build -t autoscaler -f dockerfiles/scratch/Dockerfile .`
  - test image : `docker run -it --mount type=bind,source="$(pwd)"/conf.yml,target=/conf.yml --mount type=bind,source=<path to /.kube/config>,target=/.kube/config --network=host autoscaler`
- push docker image to registry

### Building from source
Build Kubernetes Cluster Autoscaler from the source
```
go build -o bin/autoscaler cmd/main.go
```

Take a copy of `conf.yml-sample` as `conf.yml`. Fill all the required data fields.

Configure the Kubectl. This look for the kube config file in the default path `~/.kube/config` or you can run as a pon in a Kubernetes cluster.

```
.
├── autoscaler
├── conf.yml
└── plugin
    └── AWS.so
```

Copy the build plugin to `./bin/plugin` directory. Change the `CloudType` in conf.yml to `AWS`. This will load and configure the plugin.
