package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int32      `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	HashedPassword    string     `json:"-"`
	EmailVerifiedAt   NullTime   `json:"emailVerifiedAt"`
	RoleID            int32      `json:"roleId"`
	RememberTokenUUID NullString `json:"rememberTokenUuid"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	DeletedAt         NullTime   `json:"deletedAt"`
}

func (u *User) GetID() int32 {
	return u.ID
}

func (u *User) SetID(id int32) {
	u.ID = id
}

func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

type Role struct {
	ID          int32      `json:"id"`
	Name        string     `json:"name"`
	Description NullString `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (r *Role) GetID() int32 {
	return r.ID
}

func (r *Role) SetID(id int32) {
	r.ID = id
}

func (r *Role) GetCreatedAt() time.Time {
	return r.CreatedAt
}

func (r *Role) SetCreatedAt(t time.Time) {
	r.CreatedAt = t
}

func (r *Role) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

func (r *Role) SetUpdatedAt(t time.Time) {
	r.UpdatedAt = t
}

type PasswordResetToken struct {
	Email     string   `json:"email"`
	Token     string   `json:"token"`
	CreatedAt NullTime `json:"createdAt"`
}

type Session struct {
	ID           string        `json:"id"`
	UserID       sql.NullInt32 `json:"userId"`
	IpAddress    NullString    `json:"ipAddress"`
	UserAgent    NullString    `json:"userAgent"`
	Payload      string        `json:"payload"`
	LastActivity int32         `json:"lastActivity"`
}
