package utils

import "fmt"

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Message)
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s", e.Message)
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.Message)
}

func NewNotFound(msg string) error {
	return &NotFoundError{Message: msg}
}

func NewBadRequest(msg string) error {
	return &BadRequestError{Message: msg}
}

func NewUnauthorized(msg string) error {
	return &UnauthorizedError{Message: msg}
}

func NewForbidden(msg string) error {
	return &ForbiddenError{Message: msg}
}

func NewConflict(msg string) error {
	return &ConflictError{Message: msg}
}

func NewInternal(msg string) error {
	return &InternalError{Message: msg}
}
