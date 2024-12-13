# AKS audit log integration with Falco

This program is a Go version of the following program: https://github.com/sysdiglabs/aks-audit-log. The purpose of **aks-audit-log-go** is to receive Kubernetes audit logs and forward them to Falco runtime security tool that can do detections based on runtime security rules for Kubernetes API calls (using the [k8s-audit plugin](https://github.com/falcosecurity/plugins/tree/master/plugins/k8saudit)).

## Architecture

![Architecture](https://github.com/dfo-mpo/aks-audit-log-go/assets/20731423/3f1fdd1f-525e-435a-9fbd-0675134c17bd)

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

**Note:** `POSTMAXRETRIES`, `POSTRETRYINCREMENTALDELAY`, `LOGLEVEL` and `KEEPALIVE` are optional variables. They will default to 5 for `POSTMAXRETRIES`, 1000 for `POSTRETRYINCREMENTALDELAY`, "info" for `LOGLEVEL` and false for `KEEPALIVE`.

#### EventHub connection string

The `EHUBNAMESPACECONNECTIONSTRING` environment variable sets the EventHub connection string.

#### EventHub name

The `EVENTHUBNAME` environment variable sets the EventHub name.

#### Blob storage connection string

The `BLOBSTORAGECONNECTIONSTRING` environment variable sets the Blob storage connection string.

#### Blob container name

The `BLOBCONTAINERNAME` environment variable sets the Blob container name.

#### Webhook URL

The `WEBSINKURL` environment variable sets the webhook URL for the Falco pod (with k8s-audit plugin) pod. For example, `http://localhost:8765/k8s-audit`.

#### Maximum retries

The `POSTMAXRETRIES` environment variable sets the maximum number of retries for the POST request to the Falco pod (with k8s-audit plugin) pod.

#### Retry incremental delay

The `POSTRETRYINCREMENTALDELAY` environment variable sets the incremental delay between retries for the POST request to the Falco pod (with k8s-audit plugin) pod.

#### Rate limiter events per seconds

The `RATELIMITEREVENTSPERSECONDS` environment variable sets the rate limiter for the number of events per second.

#### Rate limiter burst

The `RATELIMITERBURST` environment variable sets the rate limiter for the burst.

#### Log level

The `LOGLEVEL` environment variable sets what is sent to its log. The following log levels are allowed (from highest to lowest): `panic`, `fatal`, `error`, `warn`, `info`, `debug` and `trace`.

#### Keep alive

The `KEEPALIVE` environment variable sets whether the connection to the webhook URL should be kept alive or not.

### Running

Open a terminal and move to the directory for this application.

1. Build the code:
   `CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -a -installsuffix cgo .`

2. Run the binary (after building you will have an executable binary file in the current directory):
   `./aks-audit-log-go`
