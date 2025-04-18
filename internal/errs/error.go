package errs

import "errors"

var (
	ErrInvalidParam           = errors.New("invalid param")
	ErrSendNotificationFailed = errors.New("failed to send notification")
	ErrInvalidChannel         = errors.New("invalid channel")
)
