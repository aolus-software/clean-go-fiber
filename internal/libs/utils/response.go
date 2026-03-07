package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Status:  "success",
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Status:  "success",
		Data:    data,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Success: false,
		Status:  "error",
		Message: message,
	})
}

func PaginatedOK(c *fiber.Ctx, data any, total int64, perPage, currentPage int) error {
	lastPage := 0
	if perPage > 0 {
		lastPage = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  "success",
		"data":    data,
		"meta": fiber.Map{
			"total":        total,
			"per_page":     perPage,
			"current_page": currentPage,
			"last_page":    lastPage,
		},
	})
}
