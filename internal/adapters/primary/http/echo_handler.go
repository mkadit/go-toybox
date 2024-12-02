package http

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
	"github.com/mkadit/go-toybox/internal/utils"
)

func (ad Adapter) GetAddition(ctx context.Context) {
}

func (ad Adapter) Echo(w http.ResponseWriter, r *http.Request) {
	var reqBody models.GenericRequest
	ctx := r.Context()
	msgLog := ctx.Value(logfile.MessageLog{}).(*logfile.MessageLog)

	err := utils.DecodeJSONBody(w, r, &reqBody)
	if err != nil {
		go logfile.LogErrorHTTP(msgLog.InternalID, err, models.ErrParseBody.Error())
		return
	}

	go logfile.LogMsgInterfaceHTTP(*msgLog, &reqBody)

	resp := &models.GenericResponse{Message: "Success"}

	msgLog.UpdateMessageLog("OUT", "FRONTEND", "TOPIC")
	logfile.LogMsgInterfaceHTTP(*msgLog, resp)

	render.JSON(w, r, resp)
	return
}
