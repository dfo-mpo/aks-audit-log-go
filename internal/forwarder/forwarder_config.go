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

	envVars := []struct {
		value    *string
		name   string
		errorMsg string
	}{
		{&f.EhubNamespaceConnectionString, "EHUBNAMESPACECONNECTIONSTRING", "EhubNamespaceConnectionString is not set"},
		{&f.EventHubName, "EVENTHUBNAME", "EventHubName is not set"},
		{&f.BlobStorageConnectionString, "BLOBSTORAGECONNECTIONSTRING", "BlobStorageConnectionString is not set"},
		{&f.BlobContainerName, "BLOBCONTAINERNAME", "BlobContainerName is not set"},
		{&f.WebSinkURL, "WEBSINKURL", "WebSinkURL is not set"},
	}

	for _, envVar := range envVars {
		*envVar.value = os.Getenv(envVar.name)
		if *envVar.value == "" {
			log.Fatal(envVar.errorMsg)
		}
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
