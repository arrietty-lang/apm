package apm

import "os"

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
