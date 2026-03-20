package errors

import (
	"fmt"
	"net/http"
)

type Code int

const (
	CodeSuccess Code = 0

	// 认证/授权错误 10001-10099
	CodeAuthFailed       Code = 10001
	CodeTokenExpired     Code = 10002
	CodePermissionDenied Code = 10003
	CodeUserNotFound     Code = 10004
	CodeUserExists       Code = 10005
	CodePasswordWrong    Code = 10006
	CodeInvalidToken     Code = 10007
	CodeInvalidAPIKey    Code = 10008

	// 文件相关错误 20001-20099
	CodeFileNotFound      Code = 20001
	CodeFileUploadFailed  Code = 20002
	CodePackageCreateFail Code = 20003
	CodePackageNotFound   Code = 20004
	CodeFileSizeExceed    Code = 20005
	CodeInvalidFileType   Code = 20006

	// 密码验证相关错误 30001-30099
	CodePasswordIncorrect Code = 30001
	CodePasswordExpired   Code = 30002
	CodePasswordNotValid  Code = 30003
	CodeNoActivePassword  Code = 30004
	CodeTooManyAttempts   Code = 30005

	// 加密相关错误 40001-40099
	CodeEncryptConfigInvalid Code = 40001
	CodeCryptoFailed         Code = 40002
	CodeSignatureInvalid     Code = 40003

	// 系统内部错误 90001-99999
	CodeInternalError Code = 90001
	CodeDatabaseError Code = 90002
	CodeCacheError    Code = 90003
	CodeUnknownError  Code = 99999
)

type PTMError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *PTMError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
}

func (e *PTMError) HTTPStatus() int {
	switch e.Code {
	case CodeSuccess:
		return http.StatusOK
	case CodeAuthFailed, CodeTokenExpired, CodeInvalidToken, CodeInvalidAPIKey,
		CodePasswordWrong:
		return http.StatusUnauthorized
	case CodePermissionDenied:
		return http.StatusForbidden
	case CodeUserNotFound, CodeFileNotFound, CodePackageNotFound,
		CodePasswordNotValid:
		return http.StatusNotFound
	case CodeUserExists:
		return http.StatusConflict
	case CodePasswordIncorrect, CodePasswordExpired, CodeNoActivePassword,
		CodeTooManyAttempts, CodeEncryptConfigInvalid, CodeSignatureInvalid:
		return http.StatusUnprocessableEntity
	case CodeFileUploadFailed, CodePackageCreateFail, CodeFileSizeExceed,
		CodeInvalidFileType, CodeCryptoFailed:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

var (
	ErrAuthFailed = &PTMError{Code: CodeAuthFailed, Message: "authentication failed"}

	ErrTokenExpired = &PTMError{Code: CodeTokenExpired, Message: "token expired"}

	ErrPermissionDenied = &PTMError{Code: CodePermissionDenied, Message: "permission denied"}

	ErrUserNotFound = &PTMError{Code: CodeUserNotFound, Message: "user not found"}

	ErrUserExists = &PTMError{Code: CodeUserExists, Message: "user already exists"}

	ErrPasswordWrong = &PTMError{Code: CodePasswordWrong, Message: "password is incorrect"}

	ErrInvalidToken = &PTMError{Code: CodeInvalidToken, Message: "invalid token"}

	ErrInvalidAPIKey = &PTMError{Code: CodeInvalidAPIKey, Message: "invalid api key"}

	ErrFileNotFound = &PTMError{Code: CodeFileNotFound, Message: "file not found"}

	ErrFileUploadFailed = &PTMError{Code: CodeFileUploadFailed, Message: "file upload failed"}

	ErrPackageCreateFail = &PTMError{Code: CodePackageCreateFail, Message: "package creation failed"}

	ErrPackageNotFound = &PTMError{Code: CodePackageNotFound, Message: "package not found"}

	ErrFileSizeExceed = &PTMError{Code: CodeFileSizeExceed, Message: "file size exceeds limit"}

	ErrInvalidFileType = &PTMError{Code: CodeInvalidFileType, Message: "invalid file type"}

	ErrPasswordIncorrect = &PTMError{Code: CodePasswordIncorrect, Message: "password is incorrect"}

	ErrPasswordExpired = &PTMError{Code: CodePasswordExpired, Message: "password has expired"}

	ErrPasswordNotValid = &PTMError{Code: CodePasswordNotValid, Message: "password is not yet valid"}

	ErrNoActivePassword = &PTMError{Code: CodeNoActivePassword, Message: "no active password"}

	ErrTooManyAttempts = &PTMError{Code: CodeTooManyAttempts, Message: "too many verification attempts"}

	ErrEncryptConfigInvalid = &PTMError{Code: CodeEncryptConfigInvalid, Message: "encryption configuration is invalid"}

	ErrCryptoFailed = &PTMError{Code: CodeCryptoFailed, Message: "encryption/decryption failed"}

	ErrSignatureInvalid = &PTMError{Code: CodeSignatureInvalid, Message: "signature verification failed"}

	ErrInternalError = &PTMError{Code: CodeInternalError, Message: "internal server error"}

	ErrDatabaseError = &PTMError{Code: CodeDatabaseError, Message: "database error"}
)

func New(code Code, message string) *PTMError {
	return &PTMError{Code: code, Message: message}
}

func NewWithDetail(code Code, message string, detail string) *PTMError {
	return &PTMError{Code: code, Message: message, Detail: detail}
}

func Is(err, target error) bool {
	if e, ok := err.(*PTMError); ok {
		if t, ok := target.(*PTMError); ok {
			return e.Code == t.Code
		}
	}
	return false
}
