package artemis

import (
	"mnc-be-tech-test/config"
	"mnc-be-tech-test/logger"

	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/frame"
)

// these are the default options send message
var sendOptions []func(*frame.Frame) error = []func(*frame.Frame) error{
	// stomp.SendOpt.Header("AMQ_SCHEDULED_DELAY", "1"),
}

// log file to use
var log = logger.LogType.DefaultLog

func openArtemisConn() (*stomp.Conn, error) {
	conn, err := stomp.Dial("tcp", config.Env.ArtemisConn, config.Env.ArtemisOptions...)
	if err != nil {
		log.Println("Could not connect artemis")
		log.Println(err.Error())
	}
	return conn, err
}
