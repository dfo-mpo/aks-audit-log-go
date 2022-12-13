package eventhub

import (
  "context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
  "github.com/JuanPabloSGU/aks-audit-log-go/internal/forwarder"
)

func Run() {
  f := forwarder.ForwarderConfiguration{}

  config := f.InitConfig()

  checkpointStore, err := checkpoints.NewBlobStoreFromConnectionString(config.BlobStorageConnectionString, config.BlobContainerName, nil)
  if err != nil {
    panic(err)
  }

  consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(config.EhubNamespaceConnectionString, config.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
  if err != nil {
    panic(err)
  }

  defer consumerClient.Close(context.TODO())

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
        if err := processEvents(partitionClient); err != nil {
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

func processEvents(partitionClient *azeventhubs.ProcessorPartitionClient) error {
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
      fmt.Printf("Event received with body %v\n", string(event.Body))
    }

    if len(events) != 0 {
      if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1]); err != nil {
        return err
      }
    }
  }
}

func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
  defer partitionClient.Close(context.TODO())
}
