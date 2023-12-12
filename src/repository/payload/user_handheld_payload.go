package payload

import (
	"database/sql"
	"time"

	"github.com/wit-id/blueprint-backend-go/common/constants"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/utility"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type RegisterUserHandheldPayload struct {
	Name                   string `json:"name" valid:"required"`
	ProfilePictureImageURL string `json:"profile_picture_image_url" valid:"url"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email" valid:"required,email"`
	Gender                 string `json:"gender" valid:"required"`
	Address                string `json:"address"`
	Password               string `json:"password" valid:"required,length(5|50)"`
	ConfirmPassword        string `json:"confirm_password" valid:"required,length(5|50)"`
	FCMToken               string `json:"fcm_token"`
}

type UpdateUserHandheldPayload struct {
	Name                   string `json:"name" valid:"required"`
	ProfilePictureImageURL string `json:"profile_picture_image_url" valid:"url"`
	Phone                  string `json:"phone"`
	Gender                 string `json:"gender" valid:"required"`
	Address                string `json:"address"`
}

type UpdateUserHandheldFCMPayload struct {
	FCM string `json:"fcm" valid:"required"`
}

type UpdateUserHandheldPasswordPayload struct {
	OldPassword        string `json:"old_password" valid:"required,length(5|50)"`
	NewPassword        string `json:"new_password" valid:"required,length(5|50)"`
	ConfirmNewPassword string `json:"confirm_new_password" valid:"required,length(5|50)"`
}

type ListUserHandheldPayload struct {
	Filter ListUserHandheldFilterPayload `json:"filter"`
	Limit  int32                         `json:"limit" valid:"required"`
	Offset int32                         `json:"page" valid:"required"`
	Order  string                        `json:"order" valid:"required"`
	Sort   string                        `json:"sort" valid:"required"` // ASC, DESC
}

type ListUserHandheldFilterPayload struct {
	SetName     bool   `json:"set_name"`
	Name        string `json:"name"`
	SetPhone    bool   `json:"set_phone"`
	Phone       string `json:"phone"`
	SetEmail    bool   `json:"set_email"`
	Email       string `json:"email"`
	SetGender   bool   `json:"set_gender"`
	Gender      string `json:"gender"`
	SetAddress  bool   `json:"set_address"`
	Address     string `json:"address"`
	SetIsActive bool   `json:"set_is_active"`
	IsActive    bool   `json:"is_active"`
}

type readUserHandheld struct {
	GUID                   string     `json:"id"`
	Name                   string     `json:"name"`
	ProfilePictureImageURL *string    `json:"profile_picture_image_url"`
	Phone                  *string    `json:"phone"`
	Email                  string     `json:"email"`
	Gender                 string     `json:"gender"`
	Address                *string    `json:"address"`
	IsActive               bool       `json:"is_active"`
	FcmToken               *string    `json:"fcm_token"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
	LastLogin              *time.Time `json:"last_login"`
}

func (payload *RegisterUserHandheldPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	// Validate confirm password
	if payload.Password != payload.ConfirmPassword {
		err = errors.WithStack(httpservice.ErrConfirmPasswordNotMatch)
		return
	}

	return
}

func (payload *UpdateUserHandheldPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateUserHandheldFCMPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateUserHandheldPasswordPayload) Validate(userData sqlc.UserHandheld) (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	// Validate Old Password
	password := utility.GeneratePassword(userData.Salt, payload.OldPassword)
	if userData.Password != password {
		err = errors.WithStack(httpservice.ErrPasswordNotMatch)
		return
	}

	// Validate confirm password
	if payload.NewPassword != payload.ConfirmNewPassword {
		err = errors.WithStack(httpservice.ErrConfirmPasswordNotMatch)
		return
	}

	return
}

func (payload *ListUserHandheldPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *RegisterUserHandheldPayload) ToEntity(cfg config.KVStore) (data sqlc.InsertUserHandheldParams) {
	salt := utility.GeneratePasswordSalt(cfg)
	password := utility.GeneratePassword(salt, payload.Password)

	data = sqlc.InsertUserHandheldParams{
		Guid:     utility.GenerateGoogleUUID(),
		Name:     payload.Name,
		Email:    payload.Email,
		Gender:   payload.Gender,
		Salt:     salt,
		Password: password,
	}

	if payload.ProfilePictureImageURL != "" {
		data.ProfilePictureImageUrl = sql.NullString{
			String: payload.ProfilePictureImageURL,
			Valid:  true,
		}
	}

	if payload.Phone != "" {
		data.Phone = sql.NullString{
			String: utility.FormatPhoneNumber(payload.Phone),
			Valid:  true,
		}
	}

	if payload.Address != "" {
		data.Address = sql.NullString{
			String: payload.Address,
			Valid:  true,
		}
	}

	if payload.FCMToken != "" {
		data.FcmToken = sql.NullString{
			String: payload.FCMToken,
			Valid:  true,
		}
	}

	return
}

func (payload *UpdateUserHandheldPayload) ToEntity(userData sqlc.UserHandheld) (data sqlc.UpdateUserHandheldParams) {
	data = sqlc.UpdateUserHandheldParams{
		Name:   payload.Name,
		Gender: payload.Gender,
		Guid:   userData.Guid,
	}

	if payload.ProfilePictureImageURL != "" {
		data.ProfilePictureImageUrl = sql.NullString{
			String: payload.ProfilePictureImageURL,
			Valid:  true,
		}
	}

	if payload.Phone != "" {
		data.Phone = sql.NullString{
			String: utility.FormatPhoneNumber(payload.Phone),
			Valid:  true,
		}
	}

	if payload.Address != "" {
		data.Address = sql.NullString{
			String: payload.Address,
			Valid:  true,
		}
	}

	return
}

func (payload *UpdateUserHandheldFCMPayload) ToEntity(userData sqlc.UserHandheld) (data sqlc.UpdateUserHandheldFcmTokenParams) {
	data = sqlc.UpdateUserHandheldFcmTokenParams{
		FcmToken: sql.NullString{
			String: payload.FCM,
			Valid:  true,
		},
		Guid: userData.Guid,
	}

	return
}

func (payload *UpdateUserHandheldPasswordPayload) ToEntity(cfg config.KVStore, userData sqlc.UserHandheld) (data sqlc.UpdateUserHandheldPasswordParams) {
	// Generate Salt & Password
	salt := utility.GeneratePasswordSalt(cfg)
	password := utility.GeneratePassword(salt, payload.NewPassword)

	data = sqlc.UpdateUserHandheldPasswordParams{
		Password: password,
		Salt:     salt,
		Guid:     userData.Guid,
	}

	return
}

func (payload *ListUserHandheldPayload) ToEntity() (data sqlc.ListUserHandheldParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListUserHandheldParams{
		SetName:     payload.Filter.SetName,
		Name:        "%" + payload.Filter.Name + "%",
		SetPhone:    payload.Filter.SetPhone,
		Phone:       "%" + payload.Filter.Phone + "%",
		SetEmail:    payload.Filter.SetEmail,
		Email:       "%" + payload.Filter.Email + "%",
		SetGender:   payload.Filter.SetGender,
		Gender:      "%" + payload.Filter.Gender + "%",
		SetAddress:  payload.Filter.SetAddress,
		Address:     "%" + payload.Filter.Address + "%",
		SetIsActive: payload.Filter.SetIsActive,
		IsActive: sql.NullBool{
			Bool:  payload.Filter.IsActive,
			Valid: true,
		},
		LimitData: payload.Limit,
	}

	if payload.Limit == 0 {
		data.LimitData = 10
	}

	if payload.Offset == 0 {
		data.OffsetPage = (1 * data.LimitData) - data.LimitData
	} else {
		data.OffsetPage = (payload.Offset * data.LimitData) - data.LimitData
	}

	if payload.Order != "" {
		orderParam = payload.Order + " ASC"

		if payload.Sort != "" {
			orderParam = payload.Order + " " + payload.Sort
		}
	}

	data.OrderParam = orderParam

	return
}

func ToPayloadUserHandheld(userHandheld sqlc.UserHandheld) (payload readUserHandheld) {
	payload = readUserHandheld{
		GUID:      userHandheld.Guid,
		Name:      userHandheld.Name,
		Email:     userHandheld.Email,
		Gender:    userHandheld.Gender,
		IsActive:  userHandheld.IsActive.Bool,
		CreatedAt: userHandheld.CreatedAt,
		UpdatedAt: nil,
		LastLogin: nil,
	}

	if userHandheld.ProfilePictureImageUrl.Valid {
		profilePictureImageURL := userHandheld.ProfilePictureImageUrl.String
		payload.ProfilePictureImageURL = &profilePictureImageURL
	}

	if userHandheld.Phone.Valid {
		phone := userHandheld.Phone.String
		payload.Phone = &phone
	}

	if userHandheld.Address.Valid {
		address := userHandheld.Address.String
		payload.Address = &address
	}

	if userHandheld.FcmToken.Valid {
		fcmToken := userHandheld.FcmToken.String
		payload.FcmToken = &fcmToken
	}

	if userHandheld.UpdatedAt.Valid {
		updatedAt := userHandheld.UpdatedAt.Time
		payload.UpdatedAt = &updatedAt
	}

	if userHandheld.LastLogin.Valid {
		lastLogin := userHandheld.LastLogin.Time
		payload.LastLogin = &lastLogin
	}

	return
}

func ToPayloadListUserHandheld(listUserHandheld []sqlc.UserHandheld) (payload []*readUserHandheld) {
	payload = make([]*readUserHandheld, len(listUserHandheld))

	for i := range listUserHandheld {
		payload[i] = new(readUserHandheld)
		data := ToPayloadUserHandheld(listUserHandheld[i])
		payload[i] = &data
	}

	return
}
