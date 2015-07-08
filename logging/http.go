package logging

import (
	"net"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Handler returns a http.Handler that logs all request results for a given handler
func Handler(logger log.Logger, handler http.Handler) http.Handler {
	return loggingHandler{
		logger:  logger,
		handler: handler,
	}
}

// loggingHandler is a http.Handler that logs request results
type loggingHandler struct {
	logger  log.Logger
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	writer := &responseLogger{w: w}
	url := *req.URL
	h.handler.ServeHTTP(writer, req)

	status := writer.Status()
	size := writer.Size()

	host, _, err := net.SplitHostPort(req.RemoteAddr)

	if err != nil {
		host = req.RemoteAddr
	}

	uri := url.RequestURI()

	h.logger.Log("status", status, "size", size, "host", host, "uri", uri)
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP status
// code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}
