package model

import (
	"mnc-be-tech-test/util"

	"github.com/go-stomp/stomp"
)

// Artemis is model for artemis conn
type Artemis struct {
	Conn   *stomp.Conn
	Prefix string
	Skip   map[string]struct{}
}

// ArtemisSendData query changelog for table name and row id
type ArtemisSendData struct {
	Name string
	ID   uint
}

// ArtemisMessage messsage send or receive by artemis
type ArtemisMessage struct {
	Data []byte
	File []byte
}

type AuditTrail struct {
	RequestID   string
	CreatedBy   string
	CreatedAt   *util.LocalTime
	ServiceName string
	ActionType  string
	Description string
	Mode        string
}
