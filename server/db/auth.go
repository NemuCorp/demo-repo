package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/NemuCorp/demo-repo/server/logger"
)

type AuthDB struct {
	createUser     *sql.Stmt
	getUserByEmail *sql.Stmt
	getUserByID    *sql.Stmt
	createSession  *sql.Stmt
	getSession     *sql.Stmt
	deleteSession  *sql.Stmt
	deleteUserSessions *sql.Stmt
}

func NewAuthDB(conn *sql.DB) (*AuthDB, error) {
	var a AuthDB

	stmt, err := conn.Prepare(`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`)
	if err != nil {
		return nil, err
	}
	a.createUser = stmt

	stmt, err = conn.Prepare(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`)
	if err != nil {
		return nil, err
	}
	a.getUserByEmail = stmt

	stmt, err = conn.Prepare(`SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	a.getUserByID = stmt

	stmt, err = conn.Prepare(`INSERT INTO sessions (id, user_id, session_hash, expires_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, created_at, expires_at`)
	if err != nil {
		return nil, err
	}
	a.createSession = stmt

	stmt, err = conn.Prepare(`SELECT id, user_id, session_hash, created_at, expires_at FROM sessions WHERE session_hash = $1 AND expires_at > $2`)
	if err != nil {
		return nil, err
	}
	a.getSession = stmt

	stmt, err = conn.Prepare(`DELETE FROM sessions WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	a.deleteSession = stmt

	stmt, err = conn.Prepare(`DELETE FROM sessions WHERE user_id = $1`)
	if err != nil {
		return nil, err
	}
	a.deleteUserSessions = stmt

	logger.Info.Println("AuthDB prepared statements initialized")
	return &a, nil
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Session struct {
	ID          uuid.UUID `json:"id"`
	UserID      int       `json:"user_id"`
	SessionHash string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (a *AuthDB) CreateUser(email, passwordHash string) (*User, error) {
	u := &User{}
	err := a.createUser.QueryRow(email, passwordHash).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a *AuthDB) GetUserByEmail(email string) (*User, error) {
	u := &User{}
	err := a.getUserByEmail.QueryRow(email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a *AuthDB) GetUserByID(id int) (*User, error) {
	u := &User{}
	err := a.getUserByID.QueryRow(id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a *AuthDB) CreateSession(userID int, sessionHash string, expiresAt time.Time) (*Session, error) {
	s := &Session{}
	err := a.createSession.QueryRow(uuid.New(), userID, sessionHash, expiresAt).Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (a *AuthDB) GetSession(sessionHash string) (*Session, error) {
	s := &Session{}
	err := a.getSession.QueryRow(sessionHash, time.Now()).Scan(&s.ID, &s.UserID, &s.SessionHash, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (a *AuthDB) DeleteSession(id uuid.UUID) error {
	_, err := a.deleteSession.Exec(id)
	return err
}

func (a *AuthDB) DeleteUserSessions(userID int) error {
	_, err := a.deleteUserSessions.Exec(userID)
	return err
}
