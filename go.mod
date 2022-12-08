module example/aks-audit-log-go

go 1.19

require internal/eventhub v1.0.0

replace internal/eventhub => ./internal/eventhub/

require internal/forwarder v1.0.0

replace internal/forwarder => ./internal/forwarder/

require internal/webhook v1.0.0

require github.com/joho/godotenv v1.4.0 // indirect

replace internal/webhook => ./internal/webhook/
