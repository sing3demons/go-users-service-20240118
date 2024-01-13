package model

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommonModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Href      string             `json:"href,omitempty" bson:"href,omitempty"`
	Type      string             `json:"@type,omitempty" bson:"@type,omitempty"`
	UpdatedAt string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt string             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Version   string             `json:"version,omitempty" bson:"version,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
}

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Href      string             `json:"href,omitempty" bson:"href,omitempty"`
	Type      string             `json:"@type,omitempty" bson:"@type,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Version   string             `json:"version,omitempty" bson:"version,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`

	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`

	Role         string    `json:"role,omitempty" bson:"role,omitempty"`
	ProfileImage string    `json:"profileImage,omitempty" bson:"profileImage,omitempty"`
	Gender       string    `json:"gender,omitempty" bson:"gender,omitempty"`
	Birthday     string    `json:"birthday,omitempty" bson:"birthday,omitempty"`
	Profiles     []Profile `json:"profiles,omitempty" bson:"profiles,omitempty"`
}

type Profile struct {
	ID           string `json:"id,omitempty" bson:"id,omitempty"`
	LanguageCode string `json:"languageCode,omitempty" bson:"languageCode,omitempty"`
	Prefix       string `json:"prefix,omitempty" bson:"prefix,omitempty"`
	FirstName    string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	NickName     string `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
}

type IUser struct{}

type Login struct {
	Email    string
	Password string
}

type Register struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	FistName string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	NickName string `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Password string `json:"password" bson:"password"`
}

func MaskEmail(email string) string {
	// Split the email address into local part and domain
	parts := strings.Split(email, "@")
	localPart, domain := parts[0], parts[1]

	// Mask the local part (keep the first character and replace the rest with '*')
	maskedLocalPart := string(localPart[0]) + strings.Repeat("*", len(localPart)-1)

	// Create the masked email address
	maskedEmail := fmt.Sprintf("%s@%s", maskedLocalPart, domain)

	return maskedEmail
}
