package eventlog

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (e *EventAlerts) fmtPayload() string {

	a := fmt.Sprintf("{ ")
	a += fmt.Sprintf("%s : %s", strconv.Quote("AlertKey"), strconv.Quote(e.AlertKey))
	a += fmt.Sprintf(", %s : %s", strconv.Quote("Severity"), e.Severity)
	a += fmt.Sprintf(", %s : %s", strconv.Quote("AlertGroup"), strconv.Quote(e.AlertGroup))
	a += fmt.Sprintf(", %s : %s", strconv.Quote("Node"), strconv.Quote(e.Node))
	a += fmt.Sprintf(", %s : %s", strconv.Quote("Summary"), strconv.Quote(e.Summary))
	a += fmt.Sprintf("}")

	LogEvent(Debug, "Payload:", a)

	return a

}

func (e *EventAlerts) parsePayload(EventClass string, Msg string) {
	e.AlertGroup = EventClass
	e.Summary = EventClass + ": " + Msg

	if len(e.Node) == 0 {
		node, err := os.Hostname()

		if err != nil {
			log.Fatalln(err)
		}

		e.Node = node
		return
	}

	LogEvent(Info, "Node found in configuration file, will not use machine hostname, will use", e.Node, "instead")

}

func (e *EventAlerts) hookMessageBus() {
	p := e.fmtPayload()
	url := "http://localhost:31311"
	method := "POST"

	payload := strings.NewReader(p)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		LogEvent(Error, err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		LogEvent(Error, err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		LogEvent(Error, err)
		return
	}

	LogEvent(Info, "WebHook Status:", string(body))

}