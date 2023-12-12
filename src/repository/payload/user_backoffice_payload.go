package payload

import (
	"database/sql"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/constants"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/common/utility"
	sqlc "github.com/wit-id/blueprint-backend-go/src/repository/pgbo_sqlc"
	"github.com/wit-id/blueprint-backend-go/toolkit/config"
)

type RegisterUserBackofficePayload struct {
	Name                   string `json:"name" valid:"required"`
	ProfilePictureImageURL string `json:"profile_picture_image_url"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email" valid:"required,email"`
	RoleID                 int32  `json:"role_id" valid:"required"`
}

type UpdateUserBackofficePayload struct {
	Name                   string `json:"name" valid:"required"`
	ProfilePictureImageURL string `json:"profile_picture_image_url"`
	Phone                  string `json:"phone"`
}

type UpdateUserBackofficePasswordPayload struct {
	OldPassword        string `json:"old_password" valid:"required,length(5|50)"`
	NewPassword        string `json:"new_password" valid:"required,length(5|50)"`
	ConfirmNewPassword string `json:"confirm_new_password" valid:"required,length(5|50)"`
}

type ListUserBackofficePayload struct {
	Filter ListUserBackofficeFilterPayload `json:"filter"`
	Limit  int32                           `json:"limit" valid:"required"`
	Offset int32                           `json:"page" valid:"required"`
	Order  string                          `json:"order" valid:"required"`
	Sort   string                          `json:"sort" valid:"required"` // ASC, DESC
}

type ListUserBackofficeFilterPayload struct {
	SetName     bool   `json:"set_name"`
	Name        string `json:"name"`
	SetPhone    bool   `json:"set_phone"`
	Phone       string `json:"phone"`
	SetEmail    bool   `json:"set_email"`
	Email       string `json:"email"`
	SetRoleID   bool   `json:"set_role_id"`
	RoleID      int32  `json:"role_id"`
	SetIsActive bool   `json:"set_is_active"`
	IsActive    bool   `json:"is_active"`
}

type readRegisterUserBackofficePayload struct {
	GUID                   string                        `json:"id"`
	Name                   string                        `json:"name"`
	ProfilePictureImageURL *string                       `json:"profile_picture_image_url"`
	Phone                  string                        `json:"phone"`
	Email                  string                        `json:"email"`
	Role                   readUserBackofficeRolePayload `json:"role"`
	Password               string                        `json:"password,omitempty"`
	IsActive               bool                          `json:"is_active"`
	CreatedAt              time.Time                     `json:"created_at"`
	CreatedBy              string                        `json:"created_by"`
}

type readUpdateUserBackofficePayload struct {
	GUID                   string                        `json:"id"`
	Name                   string                        `json:"name"`
	ProfilePictureImageURL *string                       `json:"profile_picture_image_url"`
	Phone                  string                        `json:"phone"`
	Email                  string                        `json:"email"`
	Role                   readUserBackofficeRolePayload `json:"role"`
	IsActive               bool                          `json:"is_active"`
	CreatedAt              time.Time                     `json:"created_at"`
	CreatedBy              string                        `json:"created_by"`
}

type readUserBackofficePayload struct {
	GUID                   string                        `json:"id"`
	Name                   string                        `json:"name"`
	ProfilePictureImageURL *string                       `json:"profile_picture_image_url"`
	Phone                  string                        `json:"phone"`
	Email                  string                        `json:"email"`
	Role                   readUserBackofficeRolePayload `json:"role"`
	IsActive               bool                          `json:"is_active"`
	CreatedAt              time.Time                     `json:"created_at"`
	CreatedBy              string                        `json:"created_by"`
	UpdatedAt              *time.Time                    `json:"updated_at"`
	UpdatedBy              *string                       `json:"updated_by"`
	LastLogin              *time.Time                    `json:"last_login"`
}

type readUserBackofficeRolePayload struct {
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (payload *RegisterUserBackofficePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateUserBackofficePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *UpdateUserBackofficePasswordPayload) Validate(userData sqlc.GetUserBackofficeRow) (err error) {
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

func (payload *ListUserBackofficePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *RegisterUserBackofficePayload) ToEntity(cfg config.KVStore, userData sqlc.GetUserBackofficeRow) (data sqlc.InsertUserBackofficeParams) {
	// Generate Salt & Password
	salt := utility.GeneratePasswordSalt(cfg)
	password := utility.GeneratePassword(salt, cfg.GetString("password.default"))

	data = sqlc.InsertUserBackofficeParams{
		Guid: utility.GenerateGoogleUUID(),
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		Phone:    utility.FormatPhoneNumber(payload.Phone),
		Email:    payload.Email,
		RoleID:   payload.RoleID,
		Password: password,
		Salt:     salt,
		IsActive: sql.NullBool{
			Bool:  constants.DefaultIsActiveValue,
			Valid: true,
		},
		CreatedBy: userData.CreatedBy,
	}

	if payload.ProfilePictureImageURL != "" {
		data.ProfilePictureImageUrl = sql.NullString{
			String: payload.ProfilePictureImageURL,
			Valid:  true,
		}
	}

	return
}

func (payload *UpdateUserBackofficePayload) ToEntity(userData sqlc.GetUserBackofficeRow) (data sqlc.UpdateUserBackofficeParams) {
	data = sqlc.UpdateUserBackofficeParams{
		Name: sql.NullString{
			String: payload.Name,
			Valid:  true,
		},
		Phone: utility.FormatPhoneNumber(payload.Phone),
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		Guid: userData.Guid,
	}

	if payload.ProfilePictureImageURL != "" {
		data.ProfilePictureImageUrl = sql.NullString{
			String: payload.ProfilePictureImageURL,
			Valid:  true,
		}
	}

	return
}

func (payload *UpdateUserBackofficePasswordPayload) ToEntity(cfg config.KVStore, userData sqlc.GetUserBackofficeRow) (data sqlc.UpdateUserBackofficePasswordParams) {
	// Generate Salt & Password
	salt := utility.GeneratePasswordSalt(cfg)
	password := utility.GeneratePassword(salt, payload.NewPassword)

	data = sqlc.UpdateUserBackofficePasswordParams{
		Password: password,
		Salt:     salt,
		UpdatedBy: sql.NullString{
			String: userData.Guid,
			Valid:  true,
		},
		Guid: userData.Guid,
	}

	return
}

func (payload *ListUserBackofficePayload) ToEntity() (data sqlc.ListUserBackofficeParams) {
	orderParam := constants.DefaultOrderValue

	data = sqlc.ListUserBackofficeParams{
		SetName:     payload.Filter.SetName,
		Name:        "%" + payload.Filter.Name + "%",
		SetPhone:    payload.Filter.SetPhone,
		Phone:       "%" + payload.Filter.Phone + "%",
		SetEmail:    payload.Filter.SetEmail,
		Email:       "%" + payload.Filter.Email + "%",
		SetRoleID:   payload.Filter.SetRoleID,
		RoleID:      payload.Filter.RoleID,
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

func ToPayloadRegisterUserBackoffice(cfg config.KVStore, userBackoffice sqlc.UserBackoffice, role sqlc.UserBackofficeRole) (payload readRegisterUserBackofficePayload) {
	payload = readRegisterUserBackofficePayload{
		GUID:  userBackoffice.Guid,
		Name:  userBackoffice.Name.String,
		Phone: userBackoffice.Phone,
		Email: userBackoffice.Email,
		Role: readUserBackofficeRolePayload{
			RoleID:   role.ID,
			RoleName: role.Name,
		},
		Password:  cfg.GetString("password.default"),
		IsActive:  userBackoffice.IsActive.Bool,
		CreatedAt: userBackoffice.CreatedAt,
		CreatedBy: userBackoffice.CreatedBy,
	}

	if userBackoffice.ProfilePictureImageUrl.Valid {
		payload.ProfilePictureImageURL = &userBackoffice.ProfilePictureImageUrl.String
	}

	return
}

func ToPayloadUpdateUserBackoffice(userBackoffice sqlc.UserBackoffice, role sqlc.UserBackofficeRole) (payload readUpdateUserBackofficePayload) {
	payload = readUpdateUserBackofficePayload{
		GUID:  userBackoffice.Guid,
		Name:  userBackoffice.Name.String,
		Phone: userBackoffice.Phone,
		Email: userBackoffice.Email,
		Role: readUserBackofficeRolePayload{
			RoleID:   role.ID,
			RoleName: role.Name,
		},
		IsActive:  userBackoffice.IsActive.Bool,
		CreatedAt: userBackoffice.CreatedAt,
		CreatedBy: userBackoffice.CreatedBy,
	}

	if userBackoffice.ProfilePictureImageUrl.Valid {
		payload.ProfilePictureImageURL = &userBackoffice.ProfilePictureImageUrl.String
	}

	return
}

func ToPayloadUserBackoffice(userBackoffice sqlc.GetUserBackofficeRow) (payload readUserBackofficePayload) {
	payload = readUserBackofficePayload{
		GUID:  userBackoffice.Guid,
		Name:  userBackoffice.Name.String,
		Phone: userBackoffice.Phone,
		Email: userBackoffice.Email,
		Role: readUserBackofficeRolePayload{
			RoleID:   userBackoffice.ID,
			RoleName: userBackoffice.RoleName,
		},
		IsActive:  userBackoffice.IsActive.Bool,
		CreatedAt: userBackoffice.CreatedAt,
		CreatedBy: userBackoffice.CreatedBy,
	}

	if userBackoffice.ProfilePictureImageUrl.Valid {
		payload.ProfilePictureImageURL = &userBackoffice.ProfilePictureImageUrl.String
	}

	if userBackoffice.UpdatedAt.Valid {
		payload.UpdatedAt = &userBackoffice.UpdatedAt.Time
	}

	if userBackoffice.UpdatedBy.Valid {
		payload.UpdatedBy = &userBackoffice.UpdatedBy.String
	}

	if userBackoffice.LastLogin.Valid {
		payload.LastLogin = &userBackoffice.LastLogin.Time
	}

	return
}

func ToPayloadUserBackofficeByMail(userBackoffice sqlc.GetUserBackofficeByEmailRow) (payload readUserBackofficePayload) {
	payload = readUserBackofficePayload{
		GUID:  userBackoffice.Guid,
		Name:  userBackoffice.Name.String,
		Phone: userBackoffice.Phone,
		Email: userBackoffice.Email,
		Role: readUserBackofficeRolePayload{
			RoleID:   userBackoffice.ID,
			RoleName: userBackoffice.RoleName,
		},
		IsActive:  userBackoffice.IsActive.Bool,
		CreatedAt: userBackoffice.CreatedAt,
		CreatedBy: userBackoffice.CreatedBy,
	}

	if userBackoffice.ProfilePictureImageUrl.Valid {
		payload.ProfilePictureImageURL = &userBackoffice.ProfilePictureImageUrl.String
	}

	if userBackoffice.UpdatedAt.Valid {
		payload.UpdatedAt = &userBackoffice.UpdatedAt.Time
	}

	if userBackoffice.UpdatedBy.Valid {
		payload.UpdatedBy = &userBackoffice.UpdatedBy.String
	}

	if userBackoffice.LastLogin.Valid {
		payload.LastLogin = &userBackoffice.LastLogin.Time
	}

	return
}

func ToPayloadListUserBackoffice(listUserBackoffice []sqlc.ListUserBackofficeRow) (payload []*readUserBackofficePayload) {
	payload = make([]*readUserBackofficePayload, len(listUserBackoffice))

	for i := range listUserBackoffice {
		payload[i] = new(readUserBackofficePayload)
		data := ToPayloadUserBackoffice(sqlc.GetUserBackofficeRow(listUserBackoffice[i]))
		payload[i] = &data
	}

	return
}
