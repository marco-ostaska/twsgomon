package eventlog

import (
	"strconv"
	"strings"
)

// EventLog provides the fields to parse the event.log
// as mentioned in http://publib.boulder.ibm.com/tividd/td/TWS/SC32-1276-02/en_US/HTML/plusmst70.htm
type EventLog struct {
	EventNumber      int
	EventClass       string
	PositionalFields []string
}

// ParseIt parses the lines of ehte log
func (e EventLog) ParseIt(s string) {
	ss := strings.Fields(s)
	eventNumber, err := strconv.Atoi(ss[0])
	if err != nil {
		LogEvent(Warn, "Unable to parse line:", s)
		return
	}
	e.EventNumber = eventNumber
	e.EventClass = eventMap[e.EventNumber]
	e.PositionalFields = ss[1:]
	e.logIt()

}

func (e *EventLog) logIt() {

	LogEvent(Info, "[Parser Start]")
	LogEvent(Info, "EventNumber:          ", e.EventNumber)
	LogEvent(Info, "EventClass:           ", e.EventClass)
	LogEvent(Info, "PositionalFields:")

	for i, v := range e.PositionalFields {
		LogEvent(Info, "         Field", i+1, ":     ", v)
	}
	LogEvent(Info, "[Parser End]")

	e.alerter()

}

func (e *EventLog) alerter() {
	var alert Alerts
	LogEvent(Debug, "Parsing", ConfigFile.AlertConfigFile)
	alert.UnmarshalIt(ConfigFile.AlertConfigFile)
	alert.parseAlert(e.EventNumber, e.EventClass, e.PositionalFields)
}

var eventMap = map[int]string{
	51:  "TWS_Process_Reset",
	101: "TWS_Job_Abend",
	102: "TWS_Job_Failed",
	103: "TWS_Job_Launched",
	104: "TWS_Job_Done",
	105: "TWS_Job_Suspended",
	106: "TWS_Job_Submitted",
	107: "TWS_Job_Cancel",
	108: "TWS_Job_Ready",
	109: "TWS_Job_Hold",
	110: "TWS_Job_Restart",
	111: "TWS_Job_Failed",
	112: "TWS_Job_SuccP",
	113: "TWS_Job_Extern",
	114: "TWS_Job_INTRO",
	115: "TWS_Job_Stuck",
	116: "TWS_Job_Wait",
	117: "TWS_Job_Waitd",
	118: "TWS_Job_Sched",
	120: "TWS_Job_Late",
	121: "TWS_Job_Until_Cont",
	122: "TWS_Job_Until_Canc",
	204: "TWS_Job_Recovery_Prompt",
	119: "TWS_Job",
	151: "TWS_Schedule_Abend",
	152: "TWS_Schedule_Stuck",
	153: "TWS_Schedule_Started",
	154: "TWS_Schedule_Done",
	155: "TWS_Schedule_Susp",
	156: "TWS_Schedule_Submit",
	157: "TWS_Schedule_Cancel",
	158: "TWS_Schedule_Ready",
	159: "TWS_Schedule_Hold",
	160: "TWS_Schedule_Extern",
	161: "TWS_Schedule_CnPend",
	163: "TWS_Schedule_Late",
	164: "TWS_Schedule_Until_Cont",
	165: "TWS_Schedule_Until_Canc",
	201: "TWS_Global_Prompt",
	202: "TWS_Schedule_Prompt",
	203: "TWS_Job_Prompt",
	251: "TWS_Link_Dropped",
	252: "TWS_Link_Failed",
	301: "TWS_Domain_Manager_Switch",
}
