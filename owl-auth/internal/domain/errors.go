package domain

import "github.com/pkg/errors"

var ErrNotFound = errors.New("not found")

type DuplicateEntryError struct {
	Entity string
	Err    error
}

func (d DuplicateEntryError) Error() string { return d.Entity + ":" + d.Err.Error() }

func (d DuplicateEntryError) Unwrap() error { return d.Err }
