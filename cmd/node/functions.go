package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func purgeFunctions(store blockless.FunctionStore, workspace string) error {

	functions, err := store.RetrieveFunctions()
	if err != nil {
		return fmt.Errorf("could not retrieve functions: %w", err)
	}

	for _, function := range functions {
		err = store.RemoveFunction(function.CID)
		if err != nil {
			return fmt.Errorf("could not remove function: %w", err)
		}

		fdir := filepath.Join(workspace, function.Files)
		err = os.RemoveAll(fdir)
		if err != nil {
			return fmt.Errorf("could not remove directory: %w", err)
		}
	}

	return nil
}
