package error

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadRequest        = errors.New("your request is in a bad format")
	ErrKeyParams         = "error_params"
	ErrKeyInternalServer = "error_internal_server"
	ErrKeyIDNotFound     = "error_id_not_found"
)

type AppErrorOption func(*AppError)

// AppError is the default error struct containing detailed information about the error
type AppError struct {
	// HTTP Status code to be set in response
	Status int `json:"-"`
	// Message displays "Failed" value
	Message string `json:"message"`
	// ErrorKey is the error message that represent error_params or error_internal_server
	ErrorKey string `json:"error_key"`
	// ErrorMessage is the error message that may be displayed to end users
	ErrorMessage string      `json:"error_message"`
	ErrorData    interface{} `json:"error_data"`
}

// New generates an application error
func New(status int, opts ...AppErrorOption) *AppError {
	err := new(AppError)
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(err)
	}
	err.Status = status
	return err
}

// Error returns the error message.
func (e AppError) Error() string {
	return e.Message
}

func WithDefinition(errorMessage string, errKey string) AppErrorOption {
	return func(h *AppError) {
		h.Message = "Failed"
		h.ErrorKey = errKey
		h.ErrorMessage = errorMessage
		h.ErrorData = map[string]interface{}{}
	}
}

// Response writes an error response to client
func Response(c *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *AppError:
		return c.Status(e.Status).JSON(e)
	case validation.Errors:
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
}

// Default error bad request
func ResponseBadRequest(errMessage string) error {
	return New(fiber.StatusOK,
		WithDefinition(
			errMessage,
			ErrKeyParams))
}

func ResponseInternalServerError(errMessage string) error {
	return New(fiber.StatusOK,
		WithDefinition(
			errMessage,
			ErrKeyInternalServer))
}

func ResponseNotFound(errMessage string) error {
	return New(fiber.StatusOK,
		WithDefinition(
			errMessage,
			ErrKeyIDNotFound))
}
