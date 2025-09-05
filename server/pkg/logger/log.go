package logger

import (
	"log"
	"os"
)

var (
	Tip   = log.New(os.Stdout, "TIP: ", log.Ltime|log.Lshortfile)
	Info  = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ltime|log.Lshortfile)
)

