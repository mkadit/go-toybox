package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
)

func ValidateEcho(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		msgLog := logfile.CreateMessageLog("ECHO", "", "IN", "FRONT END", "TOPIC", "/echo")
		logfile.LogMsgInterfaceHTTP(*msgLog, &models.GenericRequest{})
		ctx := context.WithValue(r.Context(), logfile.MessageLog{}, msgLog)
		m := map[string]string{}
		c := r.Clone(r.Context())
		err := json.NewDecoder(c.Body).Decode(&m)
		if err != nil {
			fmt.Println(err)
		}
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
