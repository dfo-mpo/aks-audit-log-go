# AKS audit log integration with Falco

This program is a golang version of the following program: https://github.com/sysdiglabs/aks-audit-log. Falco runtime security tool can also detect events for Kubernetes commands. To do so, access to the Kubernetee audit logs is required to get visibility into events in the cluster. The purpose of **aks-audit-log-go** is to receive the audit logs and forward them to Falco in order to produce alerts.

## Architecture

![Architecture](https://github.com/opencost/opencost-helm-chart/assets/20731423/d0272650-1336-46c1-9600-4dbb76ab29d2)

## How program works

There are four packages: main, httpclient, forwarder and eventhub.

### forwarder package

This package does two things: (1) it sets the configurations using the environment variables and (2) starts a server to maintain statistics.

### eventhub package

This package is responsible for receiving the events from the event log and unmarshalling the event for it to be sent into a post request.

### webhook package

This package ensures that POST request to the Falco pod (with k8s-audit plugin) is properly sent.

### httpclient package

This package sends the http POST request to the Falco pod (with k8s-audit plugin) pod.

## How to run locally

### Configuration

There is a `.envrc.example` file that contains the environment variables to be configured. You can save a copy as `.envrc` and then source it using `source .envrc` to load the environment variables into your shell session.

**Note:** `POSTMAXRETRIES`, `POSTRETRYINCREMENTALDELAY` and `LOGLEVEL` are optional variables. They will default to 5 for `POSTMAXRETRIES`, 1000 for `POSTRETRYINCREMENTALDELAY` and "info" for `LOGLEVEL`.

### Log level

The `LOGLEVEL` environment variable sets what is sent to its log. The following log levels are allowed (from higest to lowest): `panic`, `fatal`, `error`, `warn`, `info`, `debug` and `trace`.

### Running

Open a terminal and move to the directory for this application.

1. Build the code:
   `CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -a -installsuffix cgo .`

2. Run the binary (after building you will have an executable binary file in the current directory):
   `./aks-audit-log-go`
