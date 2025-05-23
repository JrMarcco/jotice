package errs

import "errors"

var (
	ErrInvalidParam               = errors.New("[jotice] invalid param")
	ErrSendNotificationFailed     = errors.New("[jotice] failed to send notification")
	ErrInvalidChannel             = errors.New("[jotice] invalid channel")
	ErrInvalidSendStrategy        = errors.New("[jotice] invalid send strategy")
	ErrNoAvailableFailoverService = errors.New("[jotice] no service needs to be take over")
)
