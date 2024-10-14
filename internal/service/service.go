package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lib/pq"
	"github.com/yervsil/auth_service/domain"
	"github.com/yervsil/auth_service/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(name, email, passwordHash string) (int, error)
	UserByEmail(email string) (*domain.User, error)
	UserById(id int) (*domain.User, error)
}

type Service struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository, log  *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log: log,
	}
}

func(s *Service) EncryptPassword(req *domain.SignupRequest) (int, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to encrypt password: %s", err.Error()))
		return -1, err
	}

	id, err := s.repo.CreateUser(req.Name, req.Email, string(bytes))
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to create user in db: %s", err.Error()))
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return -1, errors.New("the user is already exist")
		}
		return -1, err
	}

	return id, nil
}

func(s *Service) Login(req *domain.SigninRequest) (*token.TokenPair, error){
	user, err := s.repo.UserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows{
			s.log.Error(fmt.Sprintf("email not found: %s", err.Error()))
			return nil, errors.New("there is no user with this email")
		}
		s.log.Error(fmt.Sprintf("failed to get user by email: %s", err.Error()))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(req.Password))
	if err != nil {
		s.log.Error(fmt.Sprintf("passwords don't match: %s", err.Error()))
		return nil, err
	}

	tp, err := token.Token(user)
	if err != nil {
		s.log.Error(fmt.Sprintf("could not generate tokens: %s", err.Error()))
		return nil, err
	}

	return tp, nil
}

func(s *Service) RefreshToken(req *domain.RefreshTokenRequest) (*token.TokenPair, error){
	data, err := token.ParseToken(req.RefreshToken)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	user, err := s.repo.UserById(data.Id)
	if err != nil {
		if err == sql.ErrNoRows{
			s.log.Error(fmt.Sprintf("id not found: %s", err.Error()))
			return nil, errors.New("user does not exitst")
		}
		s.log.Error(fmt.Sprintf("failed to get user by id: %s", err.Error()))
		return nil, err
	}

	tp, err := token.Token(user)
	if err != nil {
		s.log.Error(fmt.Sprintf("could not generate tokens: %s", err.Error()))
		return nil, err
	}

	return tp, nil
}