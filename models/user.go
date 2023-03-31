package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	IDUser          int        `json:"id_user" gorm:"type:serial;default:uuid_generate_v4();primaryKey"`
	Surname         string     `json:"surname" gorm:"type:varchar(100);not null"`
	Name            string     `json:"name" gorm:"type:varchar(100);not null"`
	MiddleName      *string    `json:"middle_name" gorm:"type:varchar(100)"`
	Login           string     `json:"login" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password        string     `json:"password" gorm:"type:varchar(255);not null"`
	AvatarURL       *string    `json:"avatar_url" gorm:"type:varchar(100);not null;default:'default.png'"`
	IsActive        bool       `json:"is_active" gorm:"not null;default:false"`
	LastActive      *time.Time `json:"last_active" gorm:"type:timestamp without time zone;not null;default:now()"`
	Department      *string    `json:"department" gorm:"type:varchar(100)"`
	DecryptPassword string     `json:"decrypt_password" gorm:"type:varchar(255);not null"`
}

type UserResponse struct {
	IDUser     int       `json:"id_user,omitempty"`
	Surname    string    `json:"surname,omitempty"`
	Name       string    `json:"name,omitempty"`
	MiddleName string    `json:"middle_name,omitempty"`
	Login      string    `json:"login,omitempty"`
	AvatarUrl  string    `json:"avatar_url,omitempty"`
	IsActive   bool      `json:"is_active,omitempty"`
	LastActive time.Time `json:"last_active,omitempty"`
	Department *string   `json:"department,omitempty"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		IDUser:     user.IDUser,
		Surname:    user.Surname,
		Name:       user.Name,
		MiddleName: *user.MiddleName,
		Login:      user.Login,
		Department: user.Department,
		AvatarUrl:  *user.AvatarURL,
		IsActive:   user.IsActive,
		LastActive: *user.LastActive,
	}
}

type SignUpRequest struct {
	Surname                string `json:"surname" validate:"required"`
	Name                   string `json:"name" validate:"required"`
	MiddleName             string `json:"middle_name"`
	Login                  string `json:"login" validate:"required,min=8"`
	Password               string `json:"password" validate:"required,min=8"`
	ConfirmPassword        string `json:"confirm_password" validate:"required"`
	DecryptPassword        string `json:"decrypt_password" validate:"required,min=8"`
	ConfirmDecryptPassword string `json:"confirm_decrypt_password" validate:"required"`
}

type AuthRequest struct {
	Login    string
	Password string
}

type UpdateUserRequest struct {
	Surname    string  `json:"surname"`
	Name       string  `json:"name"`
	MiddleName *string `json:"middle_name"`
	AvatarURL  *string `json:"avatar_url"`
	Department *string `json:"department"`
}

type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
