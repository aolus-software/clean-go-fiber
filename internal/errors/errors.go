package errors

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Code    int
	Message string
	Errors  []FieldError
	cause   error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func Field(field, message string) FieldError {
	return FieldError{Field: field, Message: message}
}

func NewBadRequestError(message string, fields ...FieldError) *AppError {
	return &AppError{Code: 400, Message: message, Errors: fields}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: 401, Message: message}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{Code: 403, Message: message}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: 404, Message: message}
}

func NewUnprocessableEntityError(message string, fields ...FieldError) *AppError {
	return &AppError{Code: 422, Message: message, Errors: fields}
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
	return &AppError{Code: appErr.Code, Message: appErr.Message, Errors: appErr.Errors, cause: cause}
}
