package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Response holds the response details
type Response struct {
	URL     string
	BodyLen *int64
	Err     error
}

func main() {
	startedAt := time.Now()

	// Define flags for the list of URLs, number of concurrent requests and verbose
	var urlsFlag = flag.String("urls", "", "a list of URLs separated by a comma")
	var verboseFlag = flag.Bool("verbose", false, "print more logs to monitor the process")
	var errorFlag = flag.Bool("error", false, "print error for each url at the end")

	// parse provided flags
	flag.Parse()

	// A new logger will be created based on the verboseFlag.
	// If a user does not need log info, nothing will be printed
	var err error
	var logger *zap.Logger = zap.NewNop()
	if *verboseFlag {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level.SetLevel(zapcore.DebugLevel)
		logger, err = cfg.Build()
		panicIfError(err)
	}
	defer logger.Sync()

	// Split the list of URLs
	urls := strings.Split(*urlsFlag, ",")
	logger.Debug("received urls",
		zap.Int("number", len(urls)),
		zap.Bool("verbose", *verboseFlag),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var requestHanlder = NewRequestHandler(logger)
	responseCh := requestHanlder.HandleUrls(ctx, urls)

	// the result will be stored in the slice called responses
	responses := make([]Response, 0, len(urls))
	for i := 0; i < len(urls); i++ {
		responses = append(responses, <-responseCh)
	}

	// sort http responses by the length of the body
	sort.Slice(responses, func(i, j int) bool {
		switch {
		case responses[j].BodyLen == nil:
			return true
		case responses[i].BodyLen == nil:
			return false
		default:
			return *responses[i].BodyLen > *responses[j].BodyLen
		}
	})

	for _, resp := range responses {
		if *errorFlag && resp.Err != nil {
			fmt.Printf("%s - %s - %s\n", resp.URL, printInt(resp.BodyLen), resp.Err)
		} else {
			fmt.Printf("%s - %s\n", resp.URL, printInt(resp.BodyLen))
		}
	}

	logger.Debug("the whole process took", zap.Duration("duration", time.Since(startedAt)))
}
