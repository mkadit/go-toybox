package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	logfile "github.com/mkadit/go-toybox/internal/logger"
)

// LogMiddleware will log incoming requests, process the body efficiently, and recover from panics.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Clone the request to not modify the original request
		req := r.Clone(r.Context())

		// Create the message log
		msgLog := logfile.CreateMessageLog(
			fmt.Sprintf("%s-%s", req.Method, req.URL.Path),
			"", "IN", "CLIENT", "TOPIC", req.URL.Path)

		// // Efficiently read the body using a JSON Decoder
		// decoder := json.NewDecoder(req.Body)
		// var m any
		//
		// // Limit the body size to a reasonable limit (e.g., 10MB). Modify this as needed.
		// const maxBodySize = 10 * 1024 * 1024 // 10MB
		// bodyBuffer := bytes.NewBuffer([]byte{})
		// // Create a TeeReader to copy the body to bodyBuffer while allowing it to be processed later
		// tee := io.TeeReader(req.Body, bodyBuffer)
		//
		// // Read the body into bodyBuffer
		// if _, err := io.Copy(io.Discard, tee); err != nil {
		// 	logfile.LogErrorHTTP(msgLog.InternalID, err, "Error reading request body")
		// 	http.Error(rw, "Request body too large", http.StatusRequestEntityTooLarge)
		// 	return
		// }
		//
		// // Restore the body into the request so it can be accessed by the next handler
		// req.Body = io.NopCloser(bytes.NewReader(bodyBuffer.Bytes()))
		//
		// // Unmarshal the JSON body to the struct
		// fmt.Println(bodyBuffer)
		// if err := decoder.Decode(&m); err != nil {
		// 	logfile.LogErrorHTTP(msgLog.InternalID, err, models.ErrBadRequest.Error())
		// 	http.Error(rw, "Invalid JSON", http.StatusBadRequest)
		// 	return
		// }
		//
		// // Log the message with the request details (async logging)
		// go logfile.LogMsgInterfaceHTTP(*msgLog, &m)

		// Create a context with the message log
		ctx := context.WithValue(r.Context(), logfile.MessageLog{}, msgLog)

		// Call the next handler in the chain
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Recover from any panics in the handler chain
		defer func() {
			if r := recover(); r != nil {
				logfile.LogErrorHTTP("SYSTEM", fmt.Errorf("recovered from panic: %v", r), "panic recovery")
				log.Println("Recovered from panic:", r)
			}

			// Example of skipping logging for a specific route (e.g., login)
			if r.URL.Path == "/api/login" {
				return
			}
			// Additional processing (logging, etc.) can go here if necessary
			fmt.Println("Processed request:", r.URL.Path)
		}()
	})
}
