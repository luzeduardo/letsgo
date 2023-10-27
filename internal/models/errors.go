package models

import "errors"

// this error models helps encapsulate the model completely.
// the application is not concerned with the underlying datastore or datastore-specific errors
var (
	ErrNoRecord = errors.New("models: no mathcing record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicatedEmail = errors.New("models: dupplicate email")
)
