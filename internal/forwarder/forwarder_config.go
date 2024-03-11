package forwarder

import (
	"fmt"
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

	VerboseLevel int

	PostMaxRetries            int
	PostRetryIncrementalDelay int
	RateLimiter               float64
	RateLimiterBurst          int
}

func (f *ForwarderConfiguration) InitConfig() *ForwarderConfiguration {
	msg := fmt.Sprint("InitConfig")
	log.Debug().Msg(msg)

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

	rateLimiterENV := os.Getenv("RATELIMITEREVENTSPERSECONDS")
	rateLimiter, err := strconv.ParseFloat(rateLimiterENV, 64)
	if rateLimiterENV == "" || err != nil {
		f.RateLimiter = 10
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

	if f.VerboseLevel > 3 {
		msg := fmt.Sprintf("EventHubName: %s", f.EventHubName)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("BlobContainerName: %s", f.BlobContainerName)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("WebSinkURL : %s", f.WebSinkURL)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("VerboseLevel: %d", f.VerboseLevel)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("PostMaxRetries: %d", f.PostMaxRetries)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("PostRetryIncrementalDelay: %d", f.PostRetryIncrementalDelay)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("RateLimiter: %v", f.RateLimiter)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("RateLimiterBurst: %d", f.RateLimiterBurst)
		log.Info().Msg(msg)

		msg = fmt.Sprintf("EhubNamespaceConnectionString length: %d", len(f.EhubNamespaceConnectionString))
		log.Info().Msg(msg)

		msg = fmt.Sprintf("BlobStorageConnectionString length: %d", len(f.BlobStorageConnectionString))
		log.Info().Msg(msg)
	}

	return f
}
