package sentries

import (
	"github.com/getsentry/sentry-go"
	"nftshopping-store-api/pkg/log"
	"time"
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f251c947765c4c13bd98089459bb47ce@o572642.ingest.sentry.io/5923293",
	})
	if err != nil {
		log.Log.Error(err)
		panic(err)
	}
	// Flush buffered event before the program terminates.
	defer sentry.Flush(2 * time.Second)
}
