package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"testing"
)

func TestFormatter(t *testing.T) {
	for _, tt := range formatterTests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			logger := logrus.New()
			logger.Out = &out
			logger.SetFormatter(NewFormatter(PrettyPrintFormat, StackdriverFormat, DisableTimestampFormat))

			tt.run(logger)
			m := map[string]interface{}{}
			if err := json.Unmarshal(out.Bytes(), &m); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.out, m) {
				correct, _ := json.MarshalIndent(&tt.out, "", "  ")
				t.Log(out.String())
				t.Log("expected:")
				t.Log(string(correct))
				t.Error("invalid format")
			}
		})
	}
}

var formatterTests = []struct {
	run  func(*logrus.Logger)
	out  map[string]interface{}
	name string
}{
	{
		name: "With Field",
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Info("my log entry")
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"message":  "my log entry",
			"foo":      "bar",
		},
	},
	{
		name: "With Field, skip empty message",
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Info()
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"foo":      "bar",
		},
	},
	{
		name: "WithField, HTTPRequest and Error",
		run: func(logger *logrus.Logger) {
			app := fiber.New()
			app.Get("/test", func(ctx *fiber.Ctx) error {
				logger.WithFields(logrus.Fields{"foo": "bar", "httpRequest": &HTTPRequest{Request: ctx}}).Error("my log entry")
				return nil
			})
			req, _ := http.NewRequest("GET", "http://foo.bar/test", nil)
			app.Test(req)
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"foo":      "bar",
			"httpRequest": map[string]interface{}{
				"requestMethod": "GET",
				"requestUrl":    "/test",
			},
		},
	},
	{
		name: "WithField and WithError",
		run: func(logger *logrus.Logger) {
			logger.
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Info("my log entry")
		},
		out: map[string]interface{}{
			"severity": "INFO",
			"message":  "my log entry",
			"foo":      "bar",
			"error":    "test error",
		},
	},
	{
		name: "WithField and Error",
		run: func(logger *logrus.Logger) {
			logger.WithField("foo", "bar").Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"foo":      "bar",
		},
	},
	{
		name: "WithField, WithError and Error",
		run: func(logger *logrus.Logger) {
			logger.
				WithField("foo", "bar").
				WithError(errors.New("test error")).
				Error("my log entry")
		},
		out: map[string]interface{}{
			"severity": "ERROR",
			"message":  "my log entry",
			"foo":      "bar",
			"error":    "test error",
		},
	},
}
