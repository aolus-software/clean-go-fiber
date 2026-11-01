package utils

import (
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Errors  []apperrors.FieldError `json:"errors"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalCount int64 `json:"totalCount"`
}

type PaginatedData[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func OK(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{Status: fiber.StatusOK, Success: true, Message: message, Data: data})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Status: fiber.StatusCreated, Success: true, Message: message, Data: data})
}

func Error(c *fiber.Ctx, code int, message string, fields []apperrors.FieldError) error {
	if code == fiber.StatusBadRequest || code == fiber.StatusUnprocessableEntity {
		if fields == nil {
			fields = []apperrors.FieldError{}
		}
		return c.Status(code).JSON(ErrorResponse{Status: code, Success: false, Message: message, Errors: fields})
	}
	if code == fiber.StatusTooManyRequests {
		return c.Status(code).JSON(fiber.Map{"status": code, "success": false, "message": message})
	}
	return c.Status(code).JSON(Response{Status: code, Success: false, Message: message, Data: nil})
}
