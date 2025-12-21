package infra

import (
	"context"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/ray-d-song/yan/internal/model"
	"github.com/ray-d-song/yan/internal/repo"
)

// DBStore implements gorilla/sessions Store interface with database backend
type DBStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	repo    repo.SessionRepo
}

// NewDBStore creates a new database-backed session store with provided keys
func NewDBStore(repo repo.SessionRepo, keyPairs ...[]byte) *DBStore {
	return &DBStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
		},
		repo: repo,
	}
}

// NewSessionStore creates a new session store with auto-generated secure keys
func NewSessionStore(sessionRepo repo.SessionRepo) *DBStore {
	authKey := securecookie.GenerateRandomKey(32)
	encryptKey := securecookie.GenerateRandomKey(32)
	return NewDBStore(sessionRepo, authKey, encryptKey)
}

// Get returns a session for the given name after adding it to the registry
func (s *DBStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New creates a new session
func (s *DBStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true

	// Try to get session from cookie
	if cookie, err := r.Cookie(name); err == nil {
		if err := securecookie.DecodeMulti(name, cookie.Value, &session.ID, s.Codecs...); err == nil {
			// Load session from database
			if err := s.load(r.Context(), session); err == nil {
				session.IsNew = false
			}
		}
	}

	return session, nil
}

// Save stores the session in the database
func (s *DBStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Delete session if MaxAge < 0
	if session.Options.MaxAge < 0 {
		if err := s.repo.Delete(r.Context(), session.ID); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	// Generate session ID if new
	if session.ID == "" {
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32),
			), "=")
	}

	// Serialize session values
	data, err := json.Marshal(session.Values)
	if err != nil {
		return err
	}

	// Get user ID from session
	userID, ok := session.Values["user_id"].(int64)
	if !ok {
		userID = 0
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(session.Options.MaxAge) * time.Second)

	// Check if session exists
	existing, err := s.repo.GetByID(r.Context(), session.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existing == nil {
		// Create new session
		dbSession := &model.Session{
			SessionID: session.ID,
			UserID:    userID,
			Data:      string(data),
			ExpiresAt: expiresAt,
		}
		if err := s.repo.Create(r.Context(), dbSession); err != nil {
			return err
		}
	} else {
		// Update existing session
		existing.Data = string(data)
		existing.ExpiresAt = expiresAt
		if userID > 0 {
			existing.UserID = userID
		}
		if err := s.repo.Update(r.Context(), existing); err != nil {
			return err
		}
	}

	// Encode session ID and set cookie
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// load loads session data from database
func (s *DBStore) load(ctx context.Context, session *sessions.Session) error {
	dbSession, err := s.repo.GetByID(ctx, session.ID)
	if err != nil {
		return err
	}

	// Check if session is expired
	if dbSession.IsExpired() {
		return sql.ErrNoRows
	}

	// Deserialize session data
	if dbSession.Data != "" {
		if err := json.Unmarshal([]byte(dbSession.Data), &session.Values); err != nil {
			return err
		}
	}

	return nil
}
