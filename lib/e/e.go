// Package e ...
package e

import "fmt"

// Wrap ...
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr ...
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
