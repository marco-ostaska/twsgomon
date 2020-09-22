package eventlog

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// Alerts provides an abstraction over Alerts.json
type Alerts struct {
	Alerts []EventAlerts `json:"alerts"`
}

// EventAlerts provides an abstraction with Alrts configuration
type EventAlerts struct {
	EventNumber      int      `json:"eventNumber"`
	PositionalFields []string `json:"PositionalFields,omitempty"`
	MessageBus       bool     `json:"MessageBus,omitempty"`
	AlertKey         string   `json:"AlertKey,omitempty"`
	Severity         string   `json:"Severity,omitempty"`
	AlertGroup       string   `json:"AlertGroup,omitempty"`
	Summary          string   `json:"Summary,omitempty"`
	Node             string   `json:"Node,omitempty"`
	MessageBusURL    string   `json:"MessageBusURL,omitempty"`
}

// UnmarshalIt is responsible to open the event and parse it
func (a *Alerts) UnmarshalIt(s string) {
	b := GetConfig(s)
	err := json.Unmarshal(b, &a)

	if err != nil {
		log.Fatalln(err)
	}
}

func (e *EventAlerts) checkEventNum(eventNum int) bool {
	if e.EventNumber == eventNum {
		return true
	}

	return false
}

func (e *EventAlerts) checkPosFields(posFields []string) bool {
	for _, ps := range e.PositionalFields {
		ss := strings.Split(ps, ":")
		i, err := strconv.Atoi(ss[0])

		if err != nil {
			LogEvent(Debug, err)
			log.Println(err)
		}

		if len(posFields) < i {
			return false
		}

		if posFields[i-1] != ss[1] {
			return false
		}

	}
	return true
}

func (a *Alerts) parseAlert(eventNum int, EventClass string, PosF []string) {

	for _, v := range a.Alerts {
		alert := v.checkEventNum(eventNum)

		if alert == true && len(v.PositionalFields) > 0 {
			alert = v.checkPosFields(PosF)
			LogEvent(Info, "This will generate an alert for:", EventClass, eventNum, strings.Join(PosF, " "))

			if v.MessageBus {
				LogEvent(Debug, "Building Payload for MessageBus")
				v.parsePayload(EventClass, strings.Join(PosF, " "))
				v.hookMessageBus()
			}

		}

	}

}
