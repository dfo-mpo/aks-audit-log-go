package forwarder

import (
  "fmt"
  "log"
  "os"
  "strconv"
)

type ForwarderConfiguration struct {
  EhubNamespaceConnectionString   string
  EventHubName                    string
  BlobStorageConnectionString     string
  BlobContainerName               string
  WebSinkURL                      string

  VerboseLevel                    int

  PostMaxRetries                  int
  PostRetryIncrementalDelay       int
}

func NewForwarderConfiguration() *ForwarderConfiguration {
  return &ForwarderConfiguration{}
}

func (f *ForwarderConfiguration) InitConfig() *ForwarderConfiguration {
  fmt.Println("InitConfig")

  f.EhubNamespaceConnectionString = os.Getenv("EHUBNAMESPACECONNECTIONSTRING")
  if f.EhubNamespaceConnectionString == "" {
    log.Fatal("EhubNamespaceConnectionString is not set")
  }

  f.EventHubName = os.Getenv("EVENTHUBNAME")
  if f.EventHubName == "" {
    log.Fatal("EventHubName is not set")
  }

  f.BlobStorageConnectionString = os.Getenv("BLOBSTORAGECONNECTIONSTRING")
  if f.BlobStorageConnectionString == "" {
    log.Fatal("BlobStorageConnectionString is not set")
  }

  f.BlobContainerName = os.Getenv("BLOBCONTAINERNAME")
  if f.BlobContainerName == "" {
    log.Fatal("BlobContainerName is not set")
  }

  f.WebSinkURL = os.Getenv("WEBSINKURL")
  if f.WebSinkURL == "" {
    log.Fatal("WebSinkURL is not set")
  }

  verboseLevelFromENV := os.Getenv("VERBOSELEVEL")
  if verboseLevelFromENV != "" {
    verboseLevel, err := strconv.Atoi(verboseLevelFromENV)
    if err != nil {
      log.Fatal("VerboseLevel is not set")
    }

    f.VerboseLevel = verboseLevel
  }

  f.PostMaxRetries = 5
  f.PostRetryIncrementalDelay = 1000

  if f.VerboseLevel > 3 {
    fmt.Println("EventHubName: {0}", f.EventHubName)
    fmt.Println("BlobContainerName: {0}", f.BlobContainerName)
    fmt.Println("WebSinkURL : {0}", f.WebSinkURL)
    fmt.Println("VerboseLevel: {0}", f.VerboseLevel)

    fmt.Println("EhubNamespaceConnectionString length: {0}", len(f.EhubNamespaceConnectionString))
    fmt.Println("BlobStorageConnectionString length: {0}", len(f.BlobStorageConnectionString))
  }

  return f
}

func (f *ForwarderConfiguration) isValid() bool {
  return (f.EhubNamespaceConnectionString != "" &&
    f.BlobStorageConnectionString != "" &&
    f.WebSinkURL != "" &&
    f.EventHubName != "" &&
    f.BlobContainerName != "")
}
