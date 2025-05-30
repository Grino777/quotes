package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckStorageFolder() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error retrieving executable path: %w", err)
	}

	dir := filepath.Dir(exePath)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf("failed to find project root with `go.mod`")
		}
		dir = parent
	}

	storageDir := filepath.Join(dir, "storage")

	_, err = os.Stat(storageDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(storageDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating storage directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking storage directory: %w", err)
	}

	err = os.Chmod(storageDir, 0755)
	if err != nil {
		return fmt.Errorf("error setting permissions for storage directory: %w", err)
	}

	return nil
}
