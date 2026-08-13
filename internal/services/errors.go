package services

import "errors"

var (
	// auth
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already registered")
	ErrContextNotFound    = errors.New("user context not found")

	// token
	ErrInvalidToken   = errors.New("invalid token")
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenNotExists = errors.New("token is not exists")
	ErrTokenRevoked   = errors.New("token has been revoked")

	// server
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidInput   = errors.New("invalid input")

	// Workspace
	ErrWorkspaceNotFound = errors.New("workspace not found")
)
