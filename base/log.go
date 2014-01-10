package core

import (
	"io/ioutil"
	"log"
	"os"
)

var VLog = log.New(ioutil.Discard, "", 0)

func EnableVerboseLog() {
	VLog = log.New(os.Stderr, "DEBUG: ", 0)
}
