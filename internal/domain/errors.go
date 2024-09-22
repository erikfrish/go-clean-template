package domain

var ErrNotFound = NotFoundError{Message: "not found"}
var ErrAlreadyProcessed = AlreadyProcessedError{Message: "already processed"}
var ErrValidationError = ValidationError{Message: "validation error"}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type AlreadyProcessedError struct {
	Message string
}

func (e AlreadyProcessedError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
