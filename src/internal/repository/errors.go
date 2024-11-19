package repository

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrNotOwner = errors.New("not owner")
)
