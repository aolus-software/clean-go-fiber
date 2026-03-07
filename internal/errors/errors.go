package errors

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	cause   error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: 400, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

func NewPaymentRequiredError(message string) *AppError {
	return &AppError{Code: 402, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func NewConflictError(message string) *AppError {
	return &AppError{Code: 409, Message: message}
}

func NewUnprocessableEntityError(message string) *AppError {
	return &AppError{Code: 422, Message: message}
}

func NewTooManyRequestsError(message string) *AppError {
	return &AppError{Code: 429, Message: message}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{Code: 500, Message: message}
}

func NewServiceUnavailableError(message string) *AppError {
	return &AppError{Code: 503, Message: message}
}

func Wrap(cause error, appErr *AppError) *AppError {
	return &AppError{Code: appErr.Code, Message: appErr.Message, cause: cause}
}
