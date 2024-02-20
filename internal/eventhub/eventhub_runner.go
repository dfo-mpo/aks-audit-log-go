package eventhub

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/jemag/aks-audit-log-go/internal/forwarder"
)

func Run() {
	forwarder := forwarder.ForwarderConfiguration{}
	eventhub := HubEventUnpacker{}
	config := forwarder.InitConfig()
	eventhub.InitConfig(config)

	checkpointStore, err := checkpoints.NewBlobStoreFromConnectionString(
		config.BlobStorageConnectionString,
		config.BlobContainerName,
		nil,
	)
	if err != nil {
		panic(err)
	}

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(
		config.EhubNamespaceConnectionString,
		config.EventHubName,
		azeventhubs.DefaultConsumerGroup,
		nil,
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if cerr := consumerClient.Close(context.TODO()); cerr != nil {
			// Handle the error, you can log it or take appropriate action
			fmt.Printf("Error closing consumer client: %v\n", cerr)
		}
	}()

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		panic(err)
	}

	dispatchPartitionClients := func() {
		for {
			partitionClient := processor.NextPartitionClient(context.TODO())

			if partitionClient == nil {
				break
			}

			go func() {
				randomName, err := generate(8)
				if err != nil {
					// Handle the error, you can log it or take appropriate action
					fmt.Printf("Error generating random name: %v\n", err)
					return
				}

				if config.VerboseLevel > 1 {
					fmt.Printf("{%q} > Recieved event pack\n", randomName)
				}

				if err := processEvents(eventhub, partitionClient, randomName); err != nil {
					panic(err)
				}
			}()
		}
	}

	go dispatchPartitionClients()

	processorCtx, processorCancel := context.WithCancel(context.TODO())
	defer processorCancel()

	if err := processor.Run(processorCtx); err != nil {
		panic(err)
	}
}

func processEvents(eventhub HubEventUnpacker, partitionClient *azeventhubs.ProcessorPartitionClient, randomName string) error {
	defer closePartitionResources(partitionClient)

	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		fmt.Printf("Processing %d event(s)\n", len(events))

		for _, event := range events {
			err := eventhub.Process(event.Body, randomName)
			if err != nil {
				return err
			}

			if err := partitionClient.UpdateCheckpoint(context.TODO(), event); err != nil {
				return err
			}
		}
	}
}

func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	defer func() {
		if err := partitionClient.Close(context.TODO()); err != nil {
			// Handle the error, you can log it or take appropriate action
			fmt.Printf("Error closing partition client: %v\n", err)
		}
	}()
}

func generate(size int) (string, error) {
	alphabet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}

	return *(*string)(unsafe.Pointer(&b)), nil
}
