package response

import (
	"github.com/gofiber/fiber/v2"
)

type AppSuccessOption func(*AppSuccess)

type AppSuccess struct {
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Meta    *interface{} `json:"meta,omitempty"`
}

func SuccessMessage(message string) AppSuccessOption {
	return func(h *AppSuccess) {
		h.Message = message
	}
}

func SuccessData(data interface{}) AppSuccessOption {
	return func(h *AppSuccess) {
		h.Message = "Success"
		h.Data = &data
	}
}

func Success(c *fiber.Ctx, status int, option ...AppSuccessOption) error {
	appSuccess := new(AppSuccess)
	// Loop through each option
	for _, opt := range option {
		// Call the option giving the instantiated
		opt(appSuccess)
	}
	if status == 200 {
		return c.Status(fiber.StatusOK).JSON(*appSuccess)
	} else if status > 200 && status < 300 {
		appSuccess = new(AppSuccess)
		return c.Status(status).JSON(appSuccess)
	}
	return nil
}
