package bach

import (
	"fmt"
)

func RunScriptBefore(script string) error {
	fmt.Printf("Running Script: %s\n", script)

	if script != "" {
		fmt.Println(script)
		cmd := NewCommand(script)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
