package cmd

import (
	"fmt"
	"os"
)

func FatalOnError(error error) {
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
}

func waitOrError(tc chan bool, ec chan error, numProcPtr *int) error {
	numProcs := *numProcPtr
	for numProcs > 0 {
		select {
		case err := <-ec:
			return err
		case <-tc:
			numProcs--
		}
	}

	return nil
}
