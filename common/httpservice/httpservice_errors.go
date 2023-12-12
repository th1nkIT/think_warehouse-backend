package httpservice

import "errors"

// error message.
var (
	ErrBadRequest    = errors.New("bad request payload")
	ErrInvalidAppKey = errors.New("invalid app key")
	ErrUnknownSource = errors.New("unknown error")

	ErrMissingHeaderData = errors.New("missing header data")

	ErrInvalidToken            = errors.New("invalid token")
	ErrUnauthorizedTokenData   = errors.New("unauthorized token data")
	ErrInvalidOTP              = errors.New("invalid otp")
	ErrInvalidOTPToken         = errors.New("invalid otp token")
	ErrPasswordNotMatch        = errors.New("password not match")
	ErrConfirmPasswordNotMatch = errors.New("confirm password not match")
	ErrNoResultData            = errors.New("no result data")

	ErrUserAlreadyRegistered = errors.New("user is already registered")
	ErrUserNotFound          = errors.New("user not found")
	ErrUnauthorizedUser      = errors.New("unauthorized user")
	ErrInActiveUser          = errors.New("user not active")

	ErrRoleNotFound = errors.New("role not found")

	ErrInvalidESBPromotionCode     = errors.New("invalid ESB promotion code")
	ErrInsufficientQuantityVoucher = errors.New("insufficient quantities of voucher")
	ErrVoucherIsNotActive          = errors.New("voucher is not active")
	ErrVoucherIsExpired            = errors.New("voucher is expired")
	ErrInvalidVisitPurposeID       = errors.New("invalid visit purpose id")
	ErrInvalidMenuID               = errors.New("invalid menu id")
	ErrInvalidESBOrderID           = errors.New("invalid esb order id")

	ErrInvalidPaymentID = errors.New("invalid payment id")
)

// error message.
var (
	MsgHeaderTokenNotFound            = "Header `token` not found"
	MsgHeaderRefreshTokenNotFound     = "Header `refresh-token` not found"
	MsgHeaderTokenUnauthorized        = "Unauthorized token"
	MsgHeaderRefreshTokenUnauthorized = "Unauthorized refresh token"
	MsgIsNotLogin                     = "Please login first"
	MsgUnauthorizedUser               = "Unauthorized user"
	MsgUserNotActive                  = "User not active"
	MsgInvalidIDParam                 = "invalid id parameter"
)
