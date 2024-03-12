package eventhub

import (
	"context"
	"crypto/rand"
	"errors"
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/jemag/aks-audit-log-go/internal/forwarder"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func Run() {
	forwarder := forwarder.ForwarderConfiguration{}
	eventhub := HubEventUnpacker{}
	config := forwarder.InitConfig()
	eventhub.InitConfig(config)

	checkClient, err := container.NewClientFromConnectionString(config.BlobStorageConnectionString, config.BlobContainerName, nil)
	if err != nil {
		panic(err)
	}

	checkpointStore, err := checkpoints.NewBlobStore(checkClient, nil)
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

	defer consumerClient.Close(context.TODO())

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		panic(err)
	}

	rateLimiter := rate.NewLimiter(rate.Limit(config.RateLimiter), config.RateLimiterBurst)

	dispatchPartitionClients := func() {
		for {
			partitionClient := processor.NextPartitionClient(context.TODO())

			if partitionClient == nil {
				break
			}

			go func() {
				randomName, err := generate(8)
				if err != nil {
					log.Error().Msgf("Error generating random name: %v", err)
					return
				}

				log.Debug().Msgf("{%v} > Recieved event pack", randomName)

				if err := processEvents(eventhub, partitionClient, randomName, rateLimiter); err != nil {
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

func processEvents(eventhub HubEventUnpacker, partitionClient *azeventhubs.ProcessorPartitionClient, randomName string, rateLimiter *rate.Limiter) error {
	defer closePartitionResources(partitionClient)

	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		log.Debug().Msgf("Processing %d event(s)", len(events))

		for _, event := range events {
			err := eventhub.Process(event.Body, randomName, rateLimiter)
			if err != nil {
				return err
			}

			if err := partitionClient.UpdateCheckpoint(context.TODO(), event, nil); err != nil {
				return err
			}
		}
	}
}

func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	defer partitionClient.Close(context.TODO())
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
