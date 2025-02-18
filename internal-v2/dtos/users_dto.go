// api/dto/user.go
package dto

import (
	"time"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
)

type User struct {
	ID        string    `json:"id" openapi:"desc=User unique identifier"`
	FullName  string    `json:"name" openapi:"desc=Users name"`
	Email     string    `json:"email" openapi:"desc=Users email"`
	CreatedAt time.Time `json:"created_at" openapi:"desc=Creation timestamp"`
	UpdatedAt time.Time `json:"updated_at" openapi:"desc=Creation timestamp"`
}

func FromSQLCUser(u sqlitestore.User) User {
	return User{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type CreateUserRequest struct {
	FullName string `db:"full_name" json:"full_name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id" desc:"User unique identifier"`
	FullName  string    `json:"name" desc:"Users name"`
	Email     string    `json:"email" desc:"Users email"`
	CreatedAt time.Time `json:"created_at" desc:"Creation timestamp"`
}

func NewUserResponse(u sqlitestore.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func NewUserResponseList(users []sqlitestore.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = NewUserResponse(u)
	}
	return result
}
