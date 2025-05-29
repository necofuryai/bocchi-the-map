package repositories

import (
	"context"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// GetByAuthProvider retrieves a user by auth provider and provider ID
	GetByAuthProvider(ctx context.Context, provider entities.AuthProvider, providerID string) (*entities.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error
}