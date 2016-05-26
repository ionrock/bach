package bach

import (
	log "github.com/Sirupsen/logrus"
)

func InitLogging(debug bool) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	if debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debug("Logging debug configured")
}
