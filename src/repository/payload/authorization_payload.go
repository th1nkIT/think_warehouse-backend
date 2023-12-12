package payload

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
)

type AuthorizationBackofficePayload struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(5|50)"`
}

type AuthorizationHandheldPayload struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(5|50)"`
}

func (payload *AuthorizationBackofficePayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}

func (payload *AuthorizationHandheldPayload) Validate() (err error) {
	// Validate Payload
	if _, err = govalidator.ValidateStruct(payload); err != nil {
		err = errors.Wrapf(httpservice.ErrBadRequest, "bad request: %s", err.Error())
		return
	}

	return
}
