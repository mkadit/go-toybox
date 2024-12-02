package logfile

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mkadit/go-toybox/common/config"
)

type ErrorMessage struct {
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

type CriticalMessage struct {
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

type HttpMessage struct {
	UUID        string      `json:"uuid"`
	Command     string      `json:"command"`
	Flow        string      `json:"flow"`
	TrxID       int         `json:"trx_id"`
	CBID        string      `json:"cbid"`
	FullMessage interface{} `json:"full_message"`
}

type MessageNewWebSocket struct {
	Time        string      `json:"time"`
	UUID        string      `json:"uuid"`
	CBID        string      `json:"cbid"`
	FullMessage interface{} `json:"full_message"`
}

type MessageWebSocket struct {
	Time        string `json:"time"`
	UUID        string `json:"uuid"`
	Command     string `json:"command"`
	Flow        string `json:"flow"`
	CBID        string `json:"cbid"`
	ConnectorID int    `json:"connector_id"`
	TrxID       int64  `json:"trx_id"`
}

type MessageWebSocketResponse struct {
	MessageNewWebSocket
	Command     string `json:"command"`
	TypeMessage string `json:"type_message"`
}

type MessageLog struct {
	Time       string
	SystemName string
	InternalID string
	ReffTrx    string
	Step       int
	Flow       string
	Entity     string
	RC         string
	TypeTrx    string
	Header     string
	URL        string
	Msg        string
	Data       any
}

func CreateMessageLog(action string, reffTrx string, flow string, entity string, typeTrx string, url string) *MessageLog {
	var httpString string
	ctxReqID := fmt.Sprintf("%s-%s", action, uuid.New().String())
	if flow == "IN" {
		httpString = "request from"
	} else {
		httpString = "response to"
	}
	msg := fmt.Sprintf("%s %s: (%s)", httpString, entity, url)
	return &MessageLog{
		SystemName: config.AppName,
		ReffTrx:    reffTrx,
		InternalID: ctxReqID,
		Flow:       flow,
		Step:       1,
		Entity:     entity,
		TypeTrx:    typeTrx,
		URL:        url,
		Msg:        msg,
	}
}

func (ml *MessageLog) UpdateMessageLog(flow string, entity string, typeTrx string) {
	var httpString string
	ml.Entity = entity
	ml.TypeTrx = typeTrx

	if flow == "IN" {
		ml.Step += 1
		httpString = "request from"
	} else {
		ml.Step -= 1
		httpString = "response to"
	}

	msg := fmt.Sprintf("%s %s: (%s)", httpString, entity, ml.URL)
	ml.Flow = flow
	ml.Msg = msg
}
