package forwarder

import (
	"math"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type ForwarderConfiguration struct {
	EhubNamespaceConnectionString string
	EventHubName                  string
	BlobStorageConnectionString   string
	BlobContainerName             string
	WebSinkURL                    string

	PostMaxRetries            int
	PostRetryIncrementalDelay int
	RateLimiter               float64
	RateLimiterBurst          int
}

func (f *ForwarderConfiguration) InitConfig() *ForwarderConfiguration {
	log.Debug().Msg("InitConfig")

	f.EhubNamespaceConnectionString = os.Getenv("EHUBNAMESPACECONNECTIONSTRING")
	if f.EhubNamespaceConnectionString == "" {
		log.Fatal().Msg("EhubNamespaceConnectionString is not set")
	}

	f.EventHubName = os.Getenv("EVENTHUBNAME")
	if f.EventHubName == "" {
		log.Fatal().Msg("EventHubName is not set")
	}

	f.BlobStorageConnectionString = os.Getenv("BLOBSTORAGECONNECTIONSTRING")
	if f.BlobStorageConnectionString == "" {
		log.Fatal().Msg("BlobStorageConnectionString is not set")
	}

	f.BlobContainerName = os.Getenv("BLOBCONTAINERNAME")
	if f.BlobContainerName == "" {
		log.Fatal().Msg("BlobContainerName is not set")
	}

	f.WebSinkURL = os.Getenv("WEBSINKURL")
	if f.WebSinkURL == "" {
		log.Fatal().Msg("WebSinkURL is not set")
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

	rateLimiterENV := os.Getenv("RATELIMITEREVENTSPERSECONDS")
	rateLimiter, err := strconv.ParseFloat(rateLimiterENV, 64)
	if rateLimiterENV == "" || err != nil {
		f.RateLimiter = math.MaxFloat64
	} else {
		f.RateLimiter = rateLimiter
	}

	rateLimiterBurstENV := os.Getenv("RATELIMITERBURST")
	rateLimiterBurst, err := strconv.Atoi(rateLimiterBurstENV)
	if rateLimiterBurstENV == "" || err != nil {
		f.RateLimiterBurst = 50
	} else {
		f.RateLimiterBurst = rateLimiterBurst
	}

	log.Info().Msgf("EventHubName: %s", f.EventHubName)

	log.Info().Msgf("BlobContainerName: %s", f.BlobContainerName)

	log.Info().Msgf("WebSinkURL : %s", f.WebSinkURL)

	log.Info().Msgf("PostMaxRetries: %d", f.PostMaxRetries)

	log.Info().Msgf("PostRetryIncrementalDelay: %d", f.PostRetryIncrementalDelay)

	log.Info().Msgf("RateLimiter: %v", f.RateLimiter)

	log.Info().Msgf("RateLimiterBurst: %d", f.RateLimiterBurst)

	return f
}
