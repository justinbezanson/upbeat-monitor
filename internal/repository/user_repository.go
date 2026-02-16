package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the system.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't serialize password hash
	Tier         string    `json:"tier"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
}

// SQLiteUserRepository implements UserRepository for a SQLite database.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// CreateUser inserts a new user into the database.
func (r *SQLiteUserRepository) CreateUser(user *User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, tier, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.Tier, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *SQLiteUserRepository) GetUserByEmail(email string) (*User, error) {
	const query = `
		SELECT id, email, password_hash, tier, created_at
		FROM users
		WHERE email = ?
	`
	user := &User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Tier, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *SQLiteUserRepository) GetUserByID(id string) (*User, error) {
	const query = `
		SELECT id, email, password_hash, tier, created_at
		FROM users
		WHERE id = ?
	`
	user := &User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Tier, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}
