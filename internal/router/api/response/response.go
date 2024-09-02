package reponse

import (
	"errors"

	"ticket-booking/internal/domain/common"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ErrorMessage struct {
	Name       string                 `json:"name"`
	Code       int                    `json:"code"`
	Message    string                 `json:"message,omitempty"`
	RemoteCode int                    `json:"remoteCode,omitempty"`
	Detail     map[string]interface{} `json:"detail,omitempty"`
}

func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

func RespondWithoutBody(c *gin.Context, code int) {
	c.Status(code)
}

func RespondWithError(c *gin.Context, err error) {
	errMessage := ParseError(err)

	ctx := c.Request.Context()
	zerolog.Ctx(ctx).Error().Err(err).Str("component", "handler").Msg(errMessage.Message)
	_ = c.Error(err)
	c.AbortWithStatusJSON(errMessage.Code, errMessage)
}

func ParseError(err error) ErrorMessage {
	var domainError common.DomainError
	// We don't check if errors.As is valid or not
	// because an empty common.DomainError would return default error data.
	_ = errors.As(err, &domainError)

	return ErrorMessage{
		Name:       domainError.Name(),
		Code:       domainError.HTTPStatus(),
		Message:    domainError.ClientMsg(),
		RemoteCode: domainError.RemoteHTTPStatus(),
		Detail:     domainError.Detail(),
	}
}
