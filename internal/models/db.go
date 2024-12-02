package models

import (
	"net"
	"time"
)

// User represents the main user entity
type User struct {
	ID           int64      `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	Username     string     `db:"username" json:"username"`
	PasswordHash string     `db:"password_hash" json:"-"` // "-" to never expose in JSON
	IsActive     bool       `db:"is_active" json:"is_active"`
	IsVerified   bool       `db:"is_verified" json:"is_verified"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLogin    *time.Time `db:"last_login" json:"last_login,omitempty"`

	// Virtual fields for related data
	OAuthAccounts       []OAuthAccount       `db:"-" json:"oauth_accounts,omitempty"`
	Sessions            []UserSession        `db:"-" json:"sessions,omitempty"`
	VerificationTokens  []VerificationToken  `db:"-" json:"verification_tokens,omitempty"`
	PasswordResetTokens []PasswordResetToken `db:"-" json:"password_reset_tokens,omitempty"`
	URLs                []URL                `db:"-" json:"urls,omitempty"`
}

// OAuthAccount represents OAuth provider account linkage
type OAuthAccount struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"user_id"`
	Provider       string     `db:"provider" json:"provider"`
	ProviderUserID string     `db:"provider_user_id" json:"provider_user_id"`
	AccessToken    string     `db:"access_token" json:"-"`  // Sensitive data
	RefreshToken   string     `db:"refresh_token" json:"-"` // Sensitive data
	ExpiresAt      *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`

	// Virtual fields
	User *User `db:"-" json:"user,omitempty"`
}

// UserSession represents an active user session
type UserSession struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	TokenType string    `db:"token_type" json:"token_type"` // "refresh" or "access"
	Token     string    `db:"token" json:"-"`               // Sensitive data
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IsValid   bool      `db:"is_valid" json:"is_valid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	IPAddress net.IP    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`

	// Virtual fields
	User *User `db:"-" json:"user,omitempty"`
}

// VerificationToken represents an email verification token
type VerificationToken struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"-"` // Sensitive data
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Virtual fields
	User *User `db:"-" json:"user,omitempty"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	Token     string     `db:"token" json:"-"` // Sensitive data
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UsedAt    *time.Time `db:"used_at" json:"used_at,omitempty"`

	// Virtual fields
	User *User `db:"-" json:"user,omitempty"`
}

// URL represents a shortened URL
type URL struct {
	ID          int64      `db:"id" json:"id"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	LastVisited *time.Time `db:"last_visited" json:"last_visited,omitempty"`
	LastUpdated time.Time  `db:"last_updated" json:"last_updated"`
	LongURL     string     `db:"long_url" json:"long_url"`
	ShortURL    string     `db:"short_url" json:"short_url"`
	UserID      int64      `db:"user_id" json:"user_id"`

	// Virtual fields
	User *User `db:"-" json:"user,omitempty"`
}
