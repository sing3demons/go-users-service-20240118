package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/users/model"
	"github.com/sing3demons/users/service"
	"github.com/sing3demons/users/utils"
	logger "github.com/sirupsen/logrus"
)

type IUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

type userHandler struct {
	service service.IUserService
}

func NewUserHandler(service service.IUserService) IUserHandler {
	return &userHandler{service: service}
}

const xSessionId = "X-Session-Id"

func (u *userHandler) GetProfile(c *gin.Context) {
	sessionId := c.Writer.Header().Get(xSessionId)

	userId, ok := c.Get("userId")
	if !ok {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": "user id not found",
			"type":  "handler",
			"func":  "GetProfile",
			"file":  "userHandler",
			"tag":   "Get Profile",
		}).Error("UNAUTHORIZED")

		c.JSON(401, gin.H{
			"message": "unauthorized",
		})
		return
	}

	user, err := u.service.GetProfile(sessionId, userId.(string))
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": err.Error(),
			"type":  "handler",
			"func":  "GetProfile",
			"file":  "userHandler",
			"tag":   "error",
		}).Error("GET_PROFILE")

		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	logger.WithFields(logger.Fields{
		"uuid":   sessionId,
		"func":   "GetProfile",
		"file":   "userHandler",
		"tag":    "info",
		"result": utils.MaskSensitiveData(user),
	}).Info("GET_PROFILE")

	c.JSON(200, gin.H{
		"message": "success",
		"user":    user,
	})
}

func (u *userHandler) Register(c *gin.Context) {
	sessionId := c.Writer.Header().Get(xSessionId)

	var body model.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": err.Error(),
			"type":  "binding json",
			"body":  utils.MaskSensitiveData(body),
			"func":  "Register",
			"file":  "userHandler",
			"tag":   "error",
		}).Error("VALIDATION_ERROR")

		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := u.service.Register(sessionId, body)

	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": err.Error(),
			"type":  "handler",
			"body":  utils.MaskSensitiveData(body),
			"func":  "Register",
			"file":  "userHandler",
			"tag":   "error",
		}).Error("REGISTER")

		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	logger.WithFields(logger.Fields{
		"uuid":   sessionId,
		"body":   utils.MaskSensitiveData(body),
		"func":   "Register",
		"file":   "userHandler",
		"tag":    "info",
		"result": utils.MaskSensitiveData(result),
	}).Info("REGISTER")

	c.JSON(200, gin.H{
		"message": "success",
		"id":      result,
	})
}

func (u *userHandler) Login(c *gin.Context) {
	sessionId := c.Writer.Header().Get(xSessionId)
	var body model.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": err.Error(),
			"type":  "binding json",
			"body":  utils.MaskSensitiveData(body),
			"func":  "Login",
			"file":  "userHandler",
			"tag":   "error",
		}).Error("VALIDATION_ERROR")

		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := u.service.Login(sessionId, body)
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":  sessionId,
			"error": err.Error(),
			"type":  "handler",
			"body":  utils.MaskSensitiveData(body),
			"func":  "Login",
			"file":  "userHandler",
			"tag":   "handler",
		}).Error("LOGIN")

		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	logger.WithFields(logger.Fields{
		"uuid":   sessionId,
		"body":   utils.MaskSensitiveData(body),
		"func":   "Login",
		"file":   "userHandler",
		"tag":    "info",
		"result": token,
	}).Info("LOGIN")

	c.JSON(200, gin.H{
		"message": "success",
		"token":   token,
	})
}
