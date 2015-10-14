package log

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	if Log == nil {
		Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
}
