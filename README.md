# Logrus Stackdriver Formatter

Formatter for logrus which uses gofiber for the httpRequest instead of the default net/http Request, allowing log
entries to be recognized by the fluentd Stackdriver agent on Google Cloud Platform.

Example:

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	logForm "github.com/justDMNK/logrus-stackdrive-formatter"
	"time"
)

func main() {
	log.SetFormatter(logForm.NewFormatter())
	log.Info("hello world!")

	// log a HTTP request in your handler
	log.WithField("httpRequest", &logForm.HTTPRequest{
		Request:      ctx,
		Status:       fiber.StatusOK,
		ResponseSize: 31337,
		Latency:      123 * time.Millisecond,
	}).Info("additional info")
}
```

## Original

- https://github.com/joonix/log
