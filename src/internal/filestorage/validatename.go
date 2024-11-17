package filestorage

import (
	"errors"
	"fmt"
)

// MaxFileNameLength is the maximum length of a file name.
const MaxFileNameLength = 255

var (
	ErrFileNameTooLong = fmt.Errorf(
		"file name is too long (max length %d characters)",
		MaxFileNameLength,
	)
	ErrFileNameEmpty = errors.New("file name is empty")
)

func validateFileName(fileName string) error {
	if fileName == "" {
		return ErrFileNameEmpty
	}

	if len(fileName) > MaxFileNameLength {
		return ErrFileNameTooLong
	}

	return nil
}
