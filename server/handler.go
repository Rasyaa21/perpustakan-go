package server

import (
	"base-gin/config"
	"base-gin/domain/dto"
	"base-gin/exception"
	"base-gin/repository"
	"base-gin/util"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTrans "github.com/go-playground/validator/v10/translations/id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mssola/user_agent"
)

var (
	ErrRequestThrottled = errors.New("ratelimit")
)

type BindingErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Handler struct {
	cfg         config.Config
	idValidator ut.Translator
	accountRepo *repository.AccountRepository
}

func NewHandler(
	cfg *config.Config,
	accountRepo *repository.AccountRepository,
) *Handler {
	var idValidator ut.Translator

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		idNew := id.New()
		uni := ut.New(idNew, idNew)
		idValidator, _ = uni.GetTranslator("id")
		err := idTrans.RegisterDefaultTranslations(v, idValidator)
		if err != nil {
			log.Error().Err(err).Msg("RegisterDefaultTranslations")
		}
	}
	return &Handler{
		cfg:         *cfg,
		idValidator: idValidator,
		accountRepo: accountRepo,
	}
}

func (h *Handler) BindingError(err error) (int, dto.ErrorResponse) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		messageBag := make([]BindingErrorMessage, len(ve))
		for i, fe := range ve {
			translatedErrMsg := fe.Translate(h.idValidator)
			messageBag[i] = BindingErrorMessage{
				Field:   fe.Field(),
				Message: translatedErrMsg,
			}
		}
		return http.StatusUnprocessableEntity, dto.ErrorResponse{
			Success: false,
			Message: "Validasi error",
			Errors:  messageBag,
		}
	}
	log.Error().Err(err).Msg("Handler.BindingError")
	return http.StatusBadRequest, dto.ErrorResponse{
		Success: false,
		Message: http.StatusText(http.StatusBadRequest),
		Errors:  "terdapat kesalahan input",
	}
}

func (h *Handler) ErrorResponse(message string) dto.ErrorResponse {
	return dto.ErrorResponse{
		Success: false,
		Message: message,
	}
}

func (h *Handler) ErrorInternalServer(c *gin.Context, err error) {
	log.Error().Err(err).Msg("Handler.ErrorIntenalServer")
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		Success: false,
		Message: "terdapat kesalahan server",
	})
}

func (h *Handler) verifyAuthAccessToken(r *http.Request) (jwt.MapClaims, error) {
	strArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(strArr) != 2 {
		return nil, exception.ErrBearerTokenInvalid
	}
	return util.VerifyAuthAccessToken(h.cfg, strArr[1])
}

func (h *Handler) verifyAuthRefreshToken(r *http.Request) (jwt.MapClaims, error) {
	strArr := strings.Split(r.Header.Get("Authorization"), " ")

	if len(strArr) != 2 {
		return nil, exception.ErrBearerTokenInvalid
	}

	return util.VerifyAuthRefreshToken(h.cfg, strArr[1])
}

func (h *Handler) AuthAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := h.verifyAuthAccessToken(c.Request)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Handler.AuthAccess")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		account, err := h.accountRepo.GetByUsername(token["sub"].(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if account.ID == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Message: exception.ErrUserNotFound.Error(),
			})
			return
		}

		c.Set(ParamTokenUserID, account.ID)
		c.Set(ParamTokenUsername, account.Username)
		c.Next()
	}
}

func (h *Handler) AuthRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := h.verifyAuthRefreshToken(c.Request)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Handler.AuthRefresh")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.Set(ParamTokenUsername, token["sub"])
		c.Next()
	}
}

func (h *Handler) MaxPostSizeKb(maxSizeInKB int64) gin.HandlerFunc {
	maxSizeInByte := maxSizeInKB * (1 << 10)
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSizeInByte)
		buff, errRead := c.GetRawData()
		if errRead != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{
				Success: false,
				Message: fmt.Sprintf("berkas terlalu besar. Maksimal %d KB", maxSizeInKB),
			})
			return
		}
		buf := bytes.NewBuffer(buff)
		c.Request.Body = io.NopCloser(buf)
	}
}

func (h *Handler) MaxPostSizeMb(maxSizeInMB int64) gin.HandlerFunc {
	maxSizeInByte := maxSizeInMB * (1 << 20)
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSizeInByte)
		buff, errRead := c.GetRawData()
		if errRead != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{
				Success: false,
				Message: fmt.Sprintf("berkas terlalu besar. Maksimal %d MB", maxSizeInMB),
			})
			return
		}
		buf := bytes.NewBuffer(buff)
		c.Request.Body = io.NopCloser(buf)
	}
}

func (h *Handler) ClientInfo(c *gin.Context) dto.ClientInfo {
	userAgent := c.GetHeader("User-Agent")
	ua := user_agent.New(userAgent)
	return dto.ClientInfo{
		IPAddress: c.ClientIP(),
		UserAgent: userAgent,
		UserOS:    ua.OS(),
		UserGeo:   "",
	}
}
