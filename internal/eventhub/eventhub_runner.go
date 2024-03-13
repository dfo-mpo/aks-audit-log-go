package eventhub

import (
	"context"
	"errors"
	"time"

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
				log.Debug().Str("partition_id", partitionClient.PartitionID()).Msg("Recieved event pack")

				if err := processEvents(eventhub, partitionClient, rateLimiter); err != nil {
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

func processEvents(eventhub HubEventUnpacker, partitionClient *azeventhubs.ProcessorPartitionClient, rateLimiter *rate.Limiter) error {
	defer closePartitionResources(partitionClient)

	for {
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		log.Debug().Int("event_count", len(events)).Msg("Processing event(s)")

		for _, event := range events {
			err := eventhub.Process(event.Body, partitionClient.PartitionID(), event.SequenceNumber, rateLimiter)
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
