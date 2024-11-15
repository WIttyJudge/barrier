package hostsfile

import (
	"barrier/pkg/fsutil"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	StartTag           = "###### START barrier\n"
	EndTag             = "###### END barrier"
	DescriptionComment = "# Generated by the barrier CLI tool. DO NOT EDIT!\n"
)

var (
	ErrStartTagNotFound = errors.New("start tag not found")
	ErrEndTagNotFound   = errors.New("end tag not found")
)

type Status int

const (
	Enabled Status = iota
	Disabled
)

// File is a hosts file.
type File struct {
	file *os.File

	fileLocation   string
	backupLocation string
}

// New returns a new hostsfile wrapper.
func New() (*File, error) {
	location := location()
	backupLocation := location + ".backup"

	osFile, err := os.OpenFile(location, os.O_WRONLY|os.O_APPEND, 0o665)
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

// Backup creates a copy of hosts file with .backup suffix.
func (f *File) Backup() error {
	log.Debug().Str("location", f.backupLocation).Msg("backup hosts file..")
	return fsutil.CopyFile(f.fileLocation, f.backupLocation)
}

// Restore restores the original hosts file from its backup.
func (f *File) Restore() error {
	return fsutil.CopyFile(f.backupLocation, f.fileLocation)
}

// Write writes content to file
// It appends content to the end of file intead of rewriting it.
func (f *File) Write(content string) error {
	if _, err := f.file.WriteString(content); err != nil {
		return err
	}

	return nil
}

// Rewrite rewrites the entire file to the content provided.
// os.Create method truncated file completely, then WriteString writes new
// content.
func (f *File) Rewrite(content string) error {
	// If the file already exists, if will be truncated.
	newFile, err := os.Create(f.fileLocation)
	if err != nil {
		return err
	}

	if _, err := newFile.WriteString(content); err != nil {
		return err
	}

	return nil
}

// Read returns content of the file by its location.
func (f *File) Read() string {
	content, _ := os.ReadFile(f.fileLocation)
	return string(content)
}

// Status checks if hosts file already has domains that are being blocked.
// If StartTag exists, it means domains blocking enabled.
func (f *File) Status() Status {
	content := f.Read()

	if strings.Contains(content, StartTag) {
		return Enabled
	}

	return Disabled
}

// RemoveDomainsBlocking removes domains located between StartTag and EndTag
// that were parsed from blocklists.
func (f *File) RemoveDomainsBlocking() error {
	content := f.Read()

	startIndex := strings.Index(content, StartTag)
	if startIndex == -1 {
		return ErrStartTagNotFound
	}

	endIndex := strings.Index(content, EndTag)
	if endIndex == -1 {
		return ErrEndTagNotFound
	}
	endIndex += len(EndTag)

	resultContent := content[:startIndex] + content[endIndex:]

	if err := f.Rewrite(resultContent); err != nil {
		return err
	}

	return nil
}

func location() string {
	if runtime.GOOS == "windows" {
		// TODO: add for windows as well
	}

	return "/etc/hosts"
}
