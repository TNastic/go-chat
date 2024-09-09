package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")
	ErrFileUploadError     = newError(600, "File upload failed.")

	// more biz errors user
	ErrEmailAlreadyUse      = newError(1001, "The email is already in use.")
	ErrUserNameAlreadyUse   = newError(1002, "The username is already in use.")
	ErrUserInfoUpdateFailed = newError(1003, "User info update failed.")
	ErrUserNotFound         = newError(1004, "User not found.")
	ErrUserPasswordError    = newError(1005, "The password is incorrect.")
	ErrUserEmailNotFound    = newError(1006, "The email is not found.")
	ErrSendEmailFailed      = newError(1007, "Send email failed.")
	ErrEmailCodeError       = newError(1008, "The email code is incorrect.")
)
