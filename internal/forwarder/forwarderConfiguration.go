package forwarder

import (
  "log"
  "os"

  "github.com/joho/godotenv"
)

type ForwarderConfiguration interface {
  GetEventHubConnection()       string
  GetEventHub()                 string
  GetBlobStorageConnection()    string
  GetBlobStorage()              string
  GetWebSinkURL()               string
}

func GetEventHubConnection() string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  eventHubConnection := os.Getenv("EVENTHUBCONNECTION")
  if eventHubConnection == "" {
    log.Fatal("EventHub Connection unknown")
  }

  return eventHubConnection
}

func GetEventHub() string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  eventHub := os.Getenv("EVENTHUB")
  if eventHub == "" {
    log.Fatal("EventHub unknown")
  }

  return eventHub
}

func GetBlobStorageConnection() string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  blobStorageConnection:= os.Getenv("BLOBSTORAGECONNECTION")
  if blobStorageConnection == "" {
    log.Fatal("Blob Storage Connection unknown")
  }

  return blobStorageConnection
}

func GetBlobStorage() string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  blobStorage := os.Getenv("BLOBSTORAGE")
  if blobStorage == "" {
    log.Fatal("Blob Storage unknown")
  }

  return blobStorage
}

func GetWebSinkURL() string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  webSinkURL := os.Getenv("WEBSINKURL")
  if webSinkURL == "" {
    log.Fatal("Websink URL unknown")
  }

  return webSinkURL
}
