package utils

import "os"

func CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}
