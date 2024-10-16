package hostsfile

import (
	"barrier/pkg/fsutil"
	"fmt"
	"os"
	"runtime"
)

// File is a hosts file
type File struct {
	file *os.File

	fileLocation   string
	backupLocation string
}

// New returns a new hostsfile wrapper
func New() (*File, error) {
	location := location()
	backupLocation := location + ".backup"

	osFile, err := os.OpenFile(location, os.O_WRONLY, 0665)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	file := &File{
		file:           osFile,
		fileLocation:   location,
		backupLocation: backupLocation,
	}

	return file, nil
}

// Backup creates a copy of hosts file with .backup suffix
func (h *File) Backup() error {
	return fsutil.CopyFile(h.fileLocation, h.backupLocation)
}

// Write writes data to file
func (h *File) Write() {
	// h.file.Name
}

func location() string {
	if runtime.GOOS == "windows" {
		// TODO: add for windows as well
	}

	return "/home/wittyjudge/projects/golang/src/barrier/test/hosts"
	// return "/etc/hosts"
}