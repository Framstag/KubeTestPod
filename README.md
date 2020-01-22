# Why?

I needed a simple, fast to compile, install and modify Container/POD for Kubernetes
with some logging, health checks and a way to get it to die by curl call for
simple tests of Kubernetes resilience logic.

I have tested it under ArchLinux using minikube.

# How?

## Build application

Call the `buildGo.sh` to compile the Go binary. The script set `CGO_ENABLED=0`
to avoid differences between the build system and the Docker system (like having
different libc implementations).

## Build docker image

Call `buildDocker.sh` to build the Docker file and move into the registry. For minikube
call 

```shell script
eval $(minikube docker-env)
```

to place the Docker image directly into the minikube registry.

## Get it into Kubernetes

* (Optionally) Patch the `deployment.yaml` for your needs
* (Optionally) Create a namespace: `kubectl create namespace minitestpod`
* Import the yaml: `kubectl create -f deployment.yaml`


# Endpoints

* `/` - a nice hello world response
* `/liveness` - a HTTP 200 OK
* `/readiness` - a HTTP 200 OK
* `/die` - Immediate OS.exit(), Kubernetes should detect this and should restart the POD

# Environment

* `NAME` - Change the name in the log messages
* `PORT` - Change the port of the internl HTTP server

# Kill it

```shell script
curl --header "Host:kubetestpod.example.com" 192.168.99.100/die
```

(You need the header to get the ingress working in cases where DNS resolving is not
activated)

# Kubernetes helper

Get the pods

```shell script
kubectl get pods --namespace kubetestpod
```

Get the ingress (access form the outside)

```shell script
kubectl get ingress --namespace kubetestpod
```

Get the events

```shell script
kubectl get events --namespace kubetestpod
```

Get the logs (pod name and namespace may be different)

```shell script
kubectl logs -f kubetestpod-6877f6945b-tgx77 --namespace kubetestpod
```

# Words of warning

Go is not my first language.

# Help

Patches welcome!