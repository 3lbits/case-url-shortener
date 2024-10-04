// Note to candidates:
// This is the file you will be working in! Please implement [LinkFile.Long].

package main

import (
	"context"
	"errors"
	"net/url"
)

// LinkFile is what you should implement :)
// LinkFile reads short/long link pairs from [LinkFile.Filename].
// The syntax for the filename is explained in README.md.
type LinkFile struct {
	Filename string
}

// This is an "interface guard". This line doesn't do anything, but it doesn't
// compile if the type doesn't implement the interface.
// If you're new to Go, disregard it.
var _ Shortlinks = (*LinkFile)(nil)

// Long returns the long URL corresponding to the short name.
func (l *LinkFile) Long(ctx context.Context, short string) (*url.URL, error) {
	// TODO: Implement this!
	return nil, errors.New("implement me!")
}
