package secure

import (
	"os"
	"path/filepath"

	"shamir/pkg/utils/log"
)

func ValidateFileSize(file string, size int64) bool {
	newPath, err := filepath.Abs(file)
	if err != nil {
		log.Errorf("abs file path %q failed: %v", file, err)
		return false
	}

	stat, err := os.Stat(newPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		log.Errorf("get file %q	stat failed: %v", file, err)
		return false
	}

	if stat.Size() > size {
		log.Errorf("Error: file%q size more than %d", newPath, size)
		return false
	}

	return true
}
