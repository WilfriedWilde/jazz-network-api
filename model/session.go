package model

type SessionCreate struct {
	UUID      string
	UserID    int64
	CreatedAt string
	ExpiresAt string
}
