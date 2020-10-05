package main

import (
	"flag"
	"fmt"

	"github.com/marco-ostaska/twsgomon/internal/eventlog"
)

func main() {

	f, configFile := getFlags()

	if f {
		startMeUp(configFile)
	}

}

func getFlags() (bool, string) {

	version := flag.Bool("version", false, "display twsgomon version")
	fl := flag.String("config", "", "Json configuration file for twsgom⌕n")
	help := flag.Bool("help", false, "display this help and exit")
	flag.Parse()

	if *version {
		fmt.Println("twsgom⌕n version: 0.2.0")
		return false, ""
	}

	if *help {
		flag.Usage()
		return false, ""
	}

	if len(*fl) == 0 {
		flag.Usage()
		return false, ""
	}

	return true, *fl
}

func startMeUp(f string) {
	// Parsing Config json

	fmt.Println("Starting twsgom⌕n")

	eventlog.ConfigFile.UnmarshalIt(f)
	eventlog.LogEvent(eventlog.Info, "Starting twsgom⌕n")

	// following log
	var e eventlog.LogFile
	e.EvenlogPath = eventlog.ConfigFile.EvenlogPath
	e.StartFollow()

}
