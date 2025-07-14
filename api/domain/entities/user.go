package entities

import (
	"time"
)

// User represents a user in the domain (MVP version)
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User instance (MVP version)
func NewUser(email, name string) *User {
	now := time.Now()
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}