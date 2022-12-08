package forwarder

import (
  "fmt"
  "log"
  "os"

  "github.com/joho/godotenv"
)

func InitConfig() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  reponse := os.Getenv("TEST")

  fmt.Printf(reponse)
}
