package eventlog

// This file is based on https://github.com/google/mtail
// Adapted for the needs of twsgomon

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogFile provides an abstraction over LogToRead
type LogFile struct {
	logFile *os.File
	TwsGoMonConfig
}

// StartFollow starts read config and trigger all the rest
func (f *LogFile) StartFollow() {
	f.TwsGoMonConfig = ConfigFile
	f.NewRead()
}

func (f *LogFile) openFile() {
	LogToRead, err := filepath.Abs(f.EvenlogPath)

	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filepath.Clean(LogToRead))
	if err != nil {
		LogEvent(Fatal, err.Error())
		log.Fatalln(err)
	}

	f.logFile = file

}

// NewRead starts the log reading cycle
// initializing the *LogFile
func (f *LogFile) NewRead() {
	f.openFile()
	defer CloseFile(f.logFile)

	LogEvent(Debug, "Opening:", f.TwsgomonlogPath)
	f.ReadLog()
}

// ReadLog keeps reading the file looking up for changes
// It tracks for log rotations
func (f *LogFile) ReadLog() {
	b := make([]byte, 0, 4096)

	for {
		n, err := f.logFile.Read(b[:cap(b)])
		if err != nil && err != io.EOF {
			LogEvent(Fatal, err.Error())
			log.Fatalln(err)
		}

		if err == io.EOF {
			LogEvent(Debug, f.EvenlogPath, err.Error())
			time.Sleep(time.Second * f.Interval)
		}
		b = b[:n]

		truncated, terr := f.checkForTruncate()

		if terr != nil {
			LogEvent(Fatal, terr.Error())
			log.Fatalln("FileTruncatedError: ", err)
		}

		if truncated {
			// Try again: offset was greater than filesize and now we've seeked to start.
			continue
		}

		moved := f.isLogMoved()

		if moved {
			LogEvent(Info, f.EvenlogPath, "is deleted or moved")
			err := f.logFile.Close()
			if err != nil {
				LogEvent(Fatal, err)
				log.Fatalln(err)
			}
			LogEvent(Debug, "closing ", f.EvenlogPath, " file")
			f.NewRead()
			continue
		}
		if f.ReadFromBeginning {
			f.parseLine(b, n)
		}
		f.ReadFromBeginning = true
	}
}

// checkForTruncate checks if log file is truncated
func (f *LogFile) checkForTruncate() (bool, error) {

	currentOffset, err := f.logFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}

	fi, err := f.logFile.Stat()

	if err != nil {
		return false, err
	}

	if currentOffset > 0 && fi.Size() == 0 {
		return false, nil
	}

	if currentOffset == 0 || fi.Size() >= currentOffset {
		return false, nil
	}

	_, serr := f.logFile.Seek(0, io.SeekStart)
	LogEvent(Info, f.EvenlogPath, "is truncated")
	return true, serr

}

// isLogMoved check for log rotation when file is moved
// and a new one is created
func (f *LogFile) isLogMoved() bool {
	s1, err := f.logFile.Stat()
	if err != nil {
		// We have a fd but it's invalid, handle as a rotation (delete/create)
		LogEvent(Info, "We have a fd but it's invalid, handle as a rotation (delete/create)")
		return true
	}

	s2, err := os.Stat(f.EvenlogPath)
	if err != nil {
		return false
	}

	if !os.SameFile(s1, s2) {
		// new inode detected
		LogEvent(Debug, "new inode detected")
		return true
	}

	return false

}

// parseLine this print lines when file is changed
// this is for dev purpose only
// it will be replaced by parser in final version
func (f *LogFile) parseLine(b []byte, n int) {

	if n > 0 {

		sn := strings.Split(string(b[:n]), "\n")
		for _, v := range sn {
			if len(v) > 0 {
				var eventParser EventLog
				LogEvent(Debug, "Parsing line:", v)
				eventParser.ParseIt(v)
			}
		}
	}
}
