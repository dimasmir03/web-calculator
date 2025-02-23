package logging

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.size += size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.responseData.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(h http.Handler) http.Handler {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseData := responseData{
			status: 0,
			size:   0,
		}
		lwr := loggingResponseWriter{
			rw,
			&responseData,
		}
		h.ServeHTTP(&lwr, r)
		if responseData.status < 300 {
			Logger.WithFields(logrus.Fields{
				"uri":    r.RequestURI,
				"method": r.Method,
				"status": responseData.status,
				"size":   responseData.size,
			}).Info()
		} else {
			Logger.WithFields(logrus.Fields{
				"uri":    r.RequestURI,
				"method": r.Method,
				"status": responseData.status,
				"size":   responseData.size,
			}).Error()
		}

	})
}
