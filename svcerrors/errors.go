package svcerrors

import "errors"

// ErrNotFound is returned when a requested resource is not found.
var NotFound = errors.New("not found") // I disagree that this is an error, but the GoLang mafia says it is idiomatic to make this mistake
