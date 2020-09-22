package eventlog

import (
	"fmt"
	"log"
	"os"
)

const (
	// Debug const used for information that is diagnostically helpful
	Debug string = "debug"
	// Info const used for useful information to log
	Info string = "info"
	// Warn const used for Anything that can potentially cause application oddities
	Warn string = "warn"
	// Error const used for Any error which is fatal to the operation, but not the service or application
	Error string = "error"
	// Fatal const used for Any error that is forcing a shutdown of the service or application to prevent data loss
	Fatal string = "fatal"
)

var dbgLevel = map[string]int{
	"debug": 5,
	"info":  4,
	"warn":  3,
	"error": 2,
	"fatal": 1,
}

// LogEvent function to check it it should or not be saved to log
// based on debugLevel and debugType levels
func LogEvent(debugType string, logMessage ...interface{}) {
	if dbgLevel[ConfigFile.DebugLevel] >= dbgLevel[debugType] {
		f, err := os.OpenFile(ConfigFile.TwsgomonlogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)

		logMsg := fmt.Sprint(logMessage)
		logMsg = logMsg[1 : len(logMsg)-1]
		log.Printf("[%s] %s\n", debugType, logMsg)
	}

}
