package bach

import (
	log "github.com/Sirupsen/logrus"
)

func RunScript(script string) error {
	log.Debug("Running Script: ", script)

	if script != "" {
		cmd := NewCommand(script)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
