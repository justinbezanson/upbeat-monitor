package repository

import (
	"database/sql"
	"fmt"
	"time"
)

// Account represents a user's account, which holds billing and tier info.
type Account struct {
	ID        string    `json:"id"`
	Name      sql.NullString `json:"name"`
	Tier      string    `json:"tier"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents a user in the system.
type User struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't serialize password hash
	Role         string    `json:"role"`
	Tier         string    `json:"tier"` // Joined from accounts table
	CreatedAt    time.Time `json:"created_at"`
}

// PasswordResetToken represents a token used to reset a user's password.
type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the interface for user and account data operations.
type UserRepository interface {
	CreateAccount(account *Account) error
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUserPassword(userID string, passwordHash string) error

	CreatePasswordResetToken(token *PasswordResetToken) error
	GetPasswordResetToken(tokenHash string) (*PasswordResetToken, error)
	DeletePasswordResetTokensByUserID(userID string) error
	DeletePasswordResetToken(id string) error
}

// SQLiteUserRepository implements UserRepository for a SQLite database.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// CreateAccount inserts a new account into the database.
func (r *SQLiteUserRepository) CreateAccount(account *Account) error {
	const query = `
		INSERT INTO accounts (id, name, tier, created_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, account.ID, account.Name, account.Tier, account.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

// CreateUser inserts a new user into the database.
func (r *SQLiteUserRepository) CreateUser(user *User) error {
	const query = `
		INSERT INTO users (id, account_id, email, password_hash, role, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, user.ID, user.AccountID, user.Email, user.PasswordHash, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address, joining the tier from the accounts table.
func (r *SQLiteUserRepository) GetUserByEmail(email string) (*User, error) {
	const query = `
		SELECT u.id, u.account_id, u.email, u.password_hash, u.role, u.created_at, a.tier
		FROM users u
		JOIN accounts a ON u.account_id = a.id
		WHERE u.email = ?
	`
	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.AccountID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.Tier,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID, joining the tier from the accounts table.
func (r *SQLiteUserRepository) GetUserByID(id string) (*User, error) {
	const query = `
		SELECT u.id, u.account_id, u.email, u.password_hash, u.role, u.created_at, a.tier
		FROM users u
		JOIN accounts a ON u.account_id = a.id
		WHERE u.id = ?
	`
	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.AccountID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.Tier,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// UpdateUserPassword updates a user's password hash.
func (r *SQLiteUserRepository) UpdateUserPassword(userID string, passwordHash string) error {
	const query = `UPDATE users SET password_hash = ? WHERE id = ?`
	_, err := r.db.Exec(query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}

// CreatePasswordResetToken inserts a new password reset token into the database.
func (r *SQLiteUserRepository) CreatePasswordResetToken(token *PasswordResetToken) error {
	const query = `
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}
	return nil
}

// GetPasswordResetToken retrieves a password reset token by its hash.
func (r *SQLiteUserRepository) GetPasswordResetToken(tokenHash string) (*PasswordResetToken, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = ?
	`
	token := &PasswordResetToken{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Token not found
		}
		return nil, fmt.Errorf("failed to get password reset token: %w", err)
	}
	return token, nil
}

// DeletePasswordResetTokensByUserID deletes all password reset tokens for a user.
func (r *SQLiteUserRepository) DeletePasswordResetTokensByUserID(userID string) error {
	const query = `DELETE FROM password_reset_tokens WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete password reset tokens for user: %w", err)
	}
	return nil
}

// DeletePasswordResetToken deletes a password reset token by its ID.
func (r *SQLiteUserRepository) DeletePasswordResetToken(id string) error {
	const query = `DELETE FROM password_reset_tokens WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}
	return nil
}
