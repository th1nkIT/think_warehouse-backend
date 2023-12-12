package jwt

import (
	"time"

	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type RequestJWTToken struct {
	AppName    string
	DeviceID   string
	DeviceType string
}

type ResponseJwtToken struct {
	AppName             string
	DeviceID            string
	DeviceType          string
	Token               string
	TokenExpired        time.Time
	RefreshToken        string
	RefreshTokenExpired time.Time
}

type ResponseJwtTokenOTP struct {
	Token        string
	TokenExpired time.Time
}

// JWT token ...
func CreateJWTToken(cfg config.KVStore, request RequestJWTToken) (response ResponseJwtToken, err error) {
	tokenJwt := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	expiredToken := time.Now().Add(cfg.GetDuration("jwt.expired"))
	// The backend can also decode the token and get admin etc.
	claims := tokenJwt.Claims.(jwt.MapClaims)
	claims["app_name"] = request.AppName
	claims["device_id"] = request.DeviceID
	claims["device_type"] = request.DeviceType
	claims["exp"] = expiredToken.Unix()

	// The signing string should be secret (a generated UUID works too)
	token, err := tokenJwt.SignedString([]byte(cfg.GetString("jwt.key")))
	if err != nil {
		err = errors.Wrap(err, "failed generate jwt token")
		return
	}

	refreshTokenJwt := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	expiredRefreshToken := time.Now().Add(cfg.GetDuration("jwt.refresh_expired"))
	// The backend can also decode the token and get admin etc.
	rtClaims := refreshTokenJwt.Claims.(jwt.MapClaims)
	rtClaims["app_name"] = request.AppName
	rtClaims["device_id"] = request.DeviceID
	rtClaims["device_type"] = request.DeviceType
	rtClaims["exp"] = expiredRefreshToken.Unix()

	// The signing string should be secret (a generated UUID works too)
	refreshToken, err := refreshTokenJwt.SignedString([]byte(cfg.GetString("jwt.key")))
	if err != nil {
		err = errors.Wrap(err, "failed generate jwt refresh token")
		return
	}

	response = ResponseJwtToken{
		AppName:             request.AppName,
		DeviceID:            request.DeviceID,
		DeviceType:          request.DeviceType,
		Token:               token,
		TokenExpired:        expiredToken,
		RefreshToken:        refreshToken,
		RefreshTokenExpired: expiredRefreshToken,
	}

	return
}

func ClaimsJwtToken(cfg config.KVStore, token string) (response RequestJWTToken, err error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, errors.Wrapf(err, "Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.GetString("jwt.key")), nil
	})
	if err != nil {
		return
	}

	if jwtToken == nil {
		err = errors.WithStack(httpservice.ErrInvalidOTPToken)
		return
	}

	if !jwtToken.Valid {
		err = errors.WithStack(httpservice.ErrInvalidOTPToken)
		return
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		response = RequestJWTToken{
			AppName:    claims["app_name"].(string),
			DeviceID:   claims["device_id"].(string),
			DeviceType: claims["device_type"].(string),
		}
	}

	return
}

// JWT for OTP...
func CreateJWTTokenOTP(phoneNumber, otp string, cfg config.KVStore) (response ResponseJwtTokenOTP, err error) {
	tokenJwt := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	expiredToken := time.Now().Add(time.Duration(cfg.GetInt64("jwt.expired-otp")))
	// The backend can also decode the token and get admin etc.
	claims := tokenJwt.Claims.(jwt.MapClaims)
	claims["phone_number"] = phoneNumber
	claims["otp"] = otp

	// The signing string should be secret (a generated UUID works too)
	token, err := tokenJwt.SignedString([]byte(cfg.GetString("jwt.key-otp")))
	if err != nil {
		err = errors.Wrap(err, "failed generate jwt token otp")
		return
	}

	response = ResponseJwtTokenOTP{
		Token:        token,
		TokenExpired: expiredToken,
	}

	return
}

func ClaimsJWTTokenOtp(cfg config.KVStore, token string) (phoneNumber string, err error) {
	tokenOtp, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, errors.Wrapf(err, "Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.GetString("jwt.key-otp")), nil
	})
	if err != nil {
		return
	}

	if tokenOtp == nil {
		err = errors.WithStack(httpservice.ErrInvalidOTPToken)
		return
	}

	if !tokenOtp.Valid {
		err = errors.WithStack(httpservice.ErrInvalidOTPToken)
		return
	}

	if claims, ok := tokenOtp.Claims.(jwt.MapClaims); ok && tokenOtp.Valid {
		phoneNumber = claims["phone_number"].(string)
	}

	return
}
