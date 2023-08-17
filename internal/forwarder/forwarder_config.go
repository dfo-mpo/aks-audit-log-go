package forwarder

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type ForwarderConfiguration struct {
	EhubNamespaceConnectionString string
	EventHubName                  string
	BlobStorageConnectionString   string
	BlobContainerName             string
	WebSinkURL                    string

	VerboseLevel int

	PostMaxRetries            int
	PostRetryIncrementalDelay int
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

	verboseLevelENV := os.Getenv("VERBOSELEVEL")
  verboseLevel, err := strconv.Atoi(verboseLevelENV)
  if verboseLevelENV == "" || err != nil {
    f.VerboseLevel = 1
  } else {
    f.VerboseLevel = verboseLevel
  }

	postMaxRetriesENV := os.Getenv("POSTMAXRETRIES")
  postMaxRetries, err := strconv.Atoi(postMaxRetriesENV)
  if postMaxRetriesENV == "" || err != nil {
    f.PostMaxRetries = 5
  } else {
    f.PostMaxRetries = postMaxRetries
  }

	postRetryIncrementalDelayENV := os.Getenv("POSTRETRYINCREMENTALDELAY")
  postRetryIncrementalDelay, err := strconv.Atoi(postRetryIncrementalDelayENV)
  if postRetryIncrementalDelayENV == "" || err != nil {
    f.PostRetryIncrementalDelay = 1000
  } else {
    f.PostRetryIncrementalDelay = postRetryIncrementalDelay
  }

	if f.VerboseLevel > 3 {
		fmt.Println("EventHubName: {0}", f.EventHubName)
		fmt.Println("BlobContainerName: {0}", f.BlobContainerName)
		fmt.Println("WebSinkURL : {0}", f.WebSinkURL)
		fmt.Println("VerboseLevel: {0}", f.VerboseLevel)
    fmt.Println("PostMaxRetries: {0}", f.PostMaxRetries)
    fmt.Println("PostRetryIncrementalDelay: {0}", f.PostRetryIncrementalDelay)

		fmt.Println("EhubNamespaceConnectionString length: {0}", len(f.EhubNamespaceConnectionString))
		fmt.Println("BlobStorageConnectionString length: {0}", len(f.BlobStorageConnectionString))
	}

	return f
}
