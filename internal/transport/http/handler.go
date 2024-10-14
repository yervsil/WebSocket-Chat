package http

import (
	"context"
	"log/slog"

	_ "github.com/yervsil/auth_service/docs"
	"github.com/yervsil/auth_service/domain"
	"github.com/yervsil/auth_service/internal/token"

	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yervsil/auth_service/internal/utils"
	"gopkg.in/go-playground/validator.v9"
)


type Producer interface {
	ProduceJSONMessage(ctx context.Context, data any) error
	Close() error
}	

type Service interface{
	EncryptPassword(req *domain.SignupRequest) (int, error)
	Login(req *domain.SigninRequest) (*token.TokenPair, error)
	RefreshToken(req *domain.RefreshTokenRequest) (*token.TokenPair, error)
}

type Handler struct {
	service Service
	producer Producer
	log *slog.Logger 
}

func New(service Service, producer Producer, log *slog.Logger ) *Handler {
	return &Handler{
		service: service,
		log: log,
	}
}

// @Summary Create user
// @Tags User-creation
// @Description User creatiom
// @ModuleID userCreate
// @Accept  json
// @Produce  json
// @Param input body domain.SignupRequest true "User data"
// @Success 200 {object} utils.Response
// @Failure 400,404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Failure default {object} utils.Response
// @Router /sign-up [post]
func(h *Handler) Sign_up() func(http.ResponseWriter, *http.Request){
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error(fmt.Sprintf("could not read request body: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return
		}

		var req *domain.SignupRequest
	
		err = json.Unmarshal(body, &req) 
		if err != nil {
			h.log.Error(fmt.Sprintf("failed to decoding: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return 
		}

		err = validate.Struct(req)
    	if err != nil {
        	errors := utils.InvalidFields(err.(validator.ValidationErrors))
        	utils.SendJson(w, errors, http.StatusBadRequest)
        	return
    	}

		id, err := h.service.EncryptPassword(req)
		if err != nil {
			h.log.Error(fmt.Sprintf("failed to decoding: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		utils.SendJson(w, id, http.StatusOK)
	}
}

// SignIn godoc
// @Summary Authentication
// @Description Authentication of users by email and password, with token returning.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param data body domain.SigninRequest true "request for login from user"
// @Success 200 {object} utils.Response "Access and refresh tokens"
// @Failure 400 {object} utils.Response "Invaild request data"
// @Failure 500 {object} utils.Response "Server error"
// @Router /auth/sign-in [post]
func(h *Handler) Sign_in() func(http.ResponseWriter, *http.Request){
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error(fmt.Sprintf("could not read request body: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return
		}

		var req *domain.SigninRequest
	
		err = json.Unmarshal(body, &req) 
		if err != nil {
			h.log.Error(fmt.Sprintf("failed to decoding: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return 
		}

		err = validate.Struct(req)
    	if err != nil {
        	errors := utils.InvalidFields(err.(validator.ValidationErrors))
        	utils.SendJson(w, errors, http.StatusBadRequest)
        	return
    	}

		tp, err := h.service.Login(req)
		if err != nil {
			utils.SendJson(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		utils.SendJson(w, tp, http.StatusOK)
	}
}

// RefreshToken godoc
// @Summary Token update
// @Description updates request tokens by refresh token
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param data body domain.RefreshTokenRequest true "Request to update token"
// @Success 200 {object} utils.Response "New token"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Server error"
// @Router /auth/refresh_token [post]
func(h *Handler) Refresh_token() func(http.ResponseWriter, *http.Request){
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error(fmt.Sprintf("could not read request body: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return
		}

		var req *domain.RefreshTokenRequest
	
		err = json.Unmarshal(body, &req) 
		if err != nil {
			h.log.Error(fmt.Sprintf("failed to decoding: %s", err.Error()))
			utils.SendJson(w, err.Error(), http.StatusBadGateway)
			return 
		}

		err = validate.Struct(req)
    	if err != nil {
        	errors := utils.InvalidFields(err.(validator.ValidationErrors))
        	utils.SendJson(w, errors, http.StatusBadRequest)
        	return
    	}

		tp, err := h.service.RefreshToken(req)
		if err != nil {
			utils.SendJson(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		utils.SendJson(w, tp, http.StatusOK)
	}
}