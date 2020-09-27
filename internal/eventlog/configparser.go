package eventlog

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ConfigFile will be need in multiple ocasion
var ConfigFile TwsGoMonConfig

// TwsGoMonConfig is the struct to for twsgomon json configuration file
// provided by -c flag for twsgomon cmd
type TwsGoMonConfig struct {
	EvenlogPath       string        `json:"evenlog.path"`      // path where event.log is
	TwsgomonlogPath   string        `json:"twsgomonlog.path"`  // path to save the logs generated by twsgomon
	Interval          time.Duration `json:"interval"`          // Inteval in seconds to check for log changes
	ReadFromBeginning bool          `json:"readFromBeginning"` // Set if it should read it from begning
	DebugLevel        string        `json:"debug.level"`       // LogDebug level
	AlertConfigFile   string        `json:"alert.config"`      // Alert config File

}

// GetConfig opens the configuration file and reads its data
func GetConfig(s string) []byte {
	jsonFile, err := os.Open(filepath.Clean(s))
	if err != nil {
		log.Fatalln(err)
	}
	defer CloseFile(jsonFile)

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	return b
}

// CloseFile function to defer, because of G307 (CWE-703): Deferring unsafe method "Close"
func CloseFile(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Fatalln(err)
	}
}

// UnmarshalIt is responsible to open the event and parse it
func (t *TwsGoMonConfig) UnmarshalIt(s string) {
	b := GetConfig(s)
	err := json.Unmarshal(b, &t)

	if err != nil {
		log.Fatalln(err)
	}
}
