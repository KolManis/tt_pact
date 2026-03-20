package model

type SessionStatus string

const (
	SessionStatusPending SessionStatus = "pending"
	SessionStatusReady   SessionStatus = "ready"
	SessionStatusClosed  SessionStatus = "closed"
)

type Session struct {
	ID      string
	QRCode  string
	Status  SessionStatus
	Created int64
}

type Message struct {
	ID        int64
	From      string
	Text      string
	Timestamp int64
	SessionID string
}
