package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/users/model"
	"github.com/sing3demons/users/repository"
	"github.com/sing3demons/users/security"
	"github.com/sing3demons/users/utils"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserService interface {
	Register(session string, user model.Register) (any, error)
	Login(session string, req model.Login) (token any, err error)
	GetProfile(session string, userId string) (*model.User, error)
}

type userService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (u *userService) GetProfile(session string, userId string) (*model.User, error) {
	user, err := u.repo.FindOne(session, bson.M{"_id": u.repo.ConvertStringToObjectID(userId)}, &options.FindOneOptions{Projection: bson.M{"password": 0}})
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":   session,
			"error":  err.Error(),
			"func":   "FindOne",
			"file":   "service/user.go",
			"tag":    "GetProfile",
			"result": nil,
		}).Error("error")

		return nil, fmt.Errorf("user not found")
	}

	user.Href = fmt.Sprintf("/%s/%s", user.Type, user.ID.Hex())

	logger.WithFields(logger.Fields{
		"uuid":   session,
		"error":  nil,
		"func":   "FindOne",
		"file":   "service/user.go",
		"tag":    "GetProfile",
		"result": utils.MaskSensitiveData(user),
	}).Debug("find user success")

	return user, nil
}

func (u *userService) Register(session string, user model.Register) (any, error) {
	exist := u.repo.CheckUserExist(session, user.Email)
	if exist {
		logger.WithFields(logger.Fields{
			"uuid":    session,
			"error":   nil,
			"func":    "Register",
			"file":    "service/user.go",
			"tag":     "service",
			"headers": nil,
			"result":  exist,
		}).Debug("user already exist")
		return model.User{}, fmt.Errorf("user already exist")
	}

	hash, err := security.EncryptPassword(user.Password)
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":    session,
			"error":   err.Error(),
			"func":    "Register",
			"file":    "service/user.go",
			"tag":     "service",
			"headers": nil,
			"result":  nil,
		}).Error("error")
		return model.User{}, err
	}

	result, err := u.repo.CreateUser(session, model.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  hash,
		Type:      "users",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":    session,
			"error":   err.Error(),
			"func":    "Register",
			"file":    "service/user.go",
			"tag":     "service",
			"headers": nil,
			"result":  nil,
		}).Error("error")
		return model.User{}, err
	}

	logger.WithFields(logger.Fields{
		"uuid":    session,
		"error":   nil,
		"func":    "Register",
		"file":    "service/user.go",
		"tag":     "service",
		"headers": nil,
		"result":  result,
	}).Debug("insert success")

	return result, nil
}

func (u *userService) Login(session string, req model.Login) (any, error) {
	user, err := u.repo.FindOne(session, bson.M{"email": req.Email})
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":   session,
			"error":  err.Error(),
			"func":   "FindOne",
			"file":   "service/user.go",
			"tag":    "Login",
			"result": nil,
		}).Error("error")

		return nil, fmt.Errorf("user not found")
	}

	logger.WithFields(logger.Fields{
		"uuid":   session,
		"error":  nil,
		"func":   "FindOne",
		"file":   "service/user.go",
		"tag":    "Login",
		"result": utils.MaskSensitiveData(user),
	}).Debug("find user success")

	if err := security.VerifyPassword(user.Password, req.Password); err != nil {
		logger.WithFields(logger.Fields{
			"uuid":   session,
			"error":  err.Error(),
			"func":   "VerifyPassword",
			"file":   "service/user.go",
			"tag":    "Login",
			"result": nil,
		}).Error("error")

		return nil, err
	}

	token, err := security.GenerateToken(*user)
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":   session,
			"error":  err.Error(),
			"func":   "GenerateToken",
			"file":   "service/user.go",
			"tag":    "Login",
			"result": nil,
		}).Error("error")

		return nil, err
	}

	logger.WithFields(logger.Fields{
		"uuid":   session,
		"error":  nil,
		"func":   "GenerateToken",
		"file":   "service/user.go",
		"tag":    "Login",
		"result": utils.MaskSensitiveData(token),
	}).Debug("generate token success")

	return token, nil
}
