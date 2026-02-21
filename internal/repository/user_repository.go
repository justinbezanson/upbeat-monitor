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

// UserRepository defines the interface for user and account data operations.
type UserRepository interface {
	CreateAccount(account *Account) error
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
