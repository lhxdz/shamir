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
			return true
		}

		log.Errorf("get file %q	stat failed: %v", file, err)
		return false
	}

	return stat.Size() <= size
}
