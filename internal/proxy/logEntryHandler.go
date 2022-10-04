package proxy

import (
	"sync"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy/logger"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy/websocket"
)

var lock = &sync.Mutex{}

var instance *logEntryHandler

type logEntryHandler struct {
}

func GetLogEntryHandler(host string) *logEntryHandler {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &logEntryHandler{}
		}
	}

	return instance
}

func (logEntryHandler *logEntryHandler) Submit(logEntry *domain.LogEntry) {
	go func() {
		zapLogger := logger.GetLogger(logEntry.Service)
		zapLogger.Infow(
			logEntry.Message,
			"service", logEntry.Service,
			"method", logEntry.Method,
			"path", logEntry.Path,
			"query", logEntry.Query,
			"request", logEntry.Request,
			"response", logEntry.Response,
			"status", logEntry.Status,
			"requestHeaders", logEntry.RequestHeaders,
			"responseHeaders", logEntry.ResponseHeaders,
		)
	}()

	go func() {
		websocket.GetHubInstance().Broadcast <- logEntry
	}()
}
