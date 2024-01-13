package repository

import (
	"context"
	"time"

	"github.com/sing3demons/users/model"
	"github.com/sing3demons/users/utils"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	FindById(session string, id string) (*model.User, error)
	FindAll(session string) ([]model.User, error)
	CreateUser(session string, user model.User) (any, error)
	UpdateUser(session string, user model.User) (any, error)
	DeleteUser(session string, id string) (any, error)
	ConvertStringToObjectID(objectID string) primitive.ObjectID
	ConvertObjectIDToString(objectID primitive.ObjectID) string
	CheckUserExist(session string, email string) bool
	FindOneByEmail(session string, email string) (*model.User, error)
	FindOne(session string, filter primitive.M, opts ...*options.FindOneOptions) (*model.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) IUserRepository {
	return &userRepository{collection}
}

func (u *userRepository) FindById(session string, id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *model.User

	if err := u.collection.FindOne(ctx, bson.M{
		"_id":        u.ConvertStringToObjectID(id),
		"deleteDate": nil,
	}).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) FindAll(session string) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []model.User

	cursor, err := u.collection.Find(ctx, bson.M{
		"deleteDate": nil,
	})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *userRepository) FindOneByEmail(session string, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *model.User

	if err := u.collection.FindOne(ctx, bson.M{
		"email":      email,
		"deleteDate": nil,
	}).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepository) CheckUserExist(session string, email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User

	if err := u.collection.FindOne(ctx, bson.M{
		"email":      email,
		"deleteDate": nil,
	}).Decode(&user); err != nil {
		return false
	}

	return true
}

func (u *userRepository) CreateUser(session string, user model.User) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := u.collection.InsertOne(ctx, &user)
	if err != nil {
		logger.WithFields(logger.Fields{
			"uuid":    session,
			"error":   err.Error(),
			"func":    "CreateUser",
			"file":    "user.go",
			"tag":     "repository",
			"headers": nil,
			"result":  nil,
		}).Error("error")

		return nil, err
	}

	logger.WithFields(logger.Fields{
		"uuid":    session,
		"error":   nil,
		"func":    "CreateUser",
		"file":    "user.go",
		"tag":     "repository",
		"headers": nil,
		"result":  result,
	}).Debug("insert success")

	return result.InsertedID, nil
}

func (u *userRepository) UpdateUser(session string, user model.User) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := u.collection.UpdateOne(ctx, bson.M{
		"_id":        user.ID,
		"deleteDate": nil,
	}, bson.M{
		"$set": user,
	})
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

func (u *userRepository) DeleteUser(session string, id string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := u.collection.UpdateOne(ctx, bson.M{
		"_id":        u.ConvertStringToObjectID(id),
		"deleteDate": nil,
	}, bson.M{
		"$set": bson.M{
			"deleteDate": time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

func (u *userRepository) ConvertStringToObjectID(objectID string) primitive.ObjectID {
	primitiveObjectID, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return primitive.ObjectID{}
	}
	return primitiveObjectID
}

func (u *userRepository) ConvertObjectIDToString(objectID primitive.ObjectID) string {
	return objectID.Hex()
}

func (u *userRepository) FindOne(session string, filter primitive.M, opts ...*options.FindOneOptions) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := model.User{}
	if err := u.collection.FindOne(ctx, filter, opts...).Decode(&user); err != nil {
		logger.WithFields(logger.Fields{
			"uuid":   session,
			"error":  err.Error(),
			"func":   "CreateUser",
			"file":   "repository/user.go",
			"tag":    "repository",
			"result": nil,
			"filter": filter,
		}).Error("error")

		return nil, err
	}

	logger.WithFields(logger.Fields{
		"uuid":   session,
		"error":  nil,
		"func":   "FindOne",
		"file":   "user.go",
		"tag":    "repository",
		"result": utils.MaskSensitiveData(user),
		"filter": filter,
	}).Debug("insert success")

	return &user, nil
}
