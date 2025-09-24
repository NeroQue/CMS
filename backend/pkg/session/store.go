package session

import (
	"context"
	"log"
	"sync"

	"github.com/NeroQue/course-management-backend/internal/database"
	"github.com/google/uuid"
)

// SessionStore manages user sessions - kinda like a simple auth system
type SessionStore struct {
	DB             *database.Queries
	mu             sync.RWMutex      // for thread safety
	currentSession *database.Session // cache current user
}

// global session store - not ideal but works for now
var store *SessionStore

// Initialize sets up the session store with database
func Initialize(db *database.Queries) {
	store = &SessionStore{
		DB:             db,
		currentSession: nil,
	}

	// try to load any existing session on startup
	go loadActiveSession()
}

// loadActiveSession tries to restore the last active session
func loadActiveSession() {
	if store == nil || store.DB == nil {
		log.Println("Warning: Cannot load active session, session store not initialized")
		return
	}

	session, err := store.DB.GetActiveSession(context.Background())
	if err != nil {
		// no big deal if there's no active session
		return
	}

	store.mu.Lock()
	store.currentSession = &session
	store.mu.Unlock()
}

// SetCurrentUser sets the currently logged in user
func SetCurrentUser(userID uuid.UUID) {
	if store == nil || store.DB == nil {
		log.Println("Warning: Cannot set current user, session store not initialized")
		return
	}

	// Delete any existing sessions first
	// This ensures we only have one active session
	ClearAllSessions()

	// Create a new session in the database
	sessionID := uuid.New()
	session, err := store.DB.CreateSession(context.Background(), database.CreateSessionParams{
		ID:     sessionID,
		UserID: userID,
	})

	if err != nil {
		log.Printf("Error creating session: %v", err)
		return
	}

	// Cache the new session
	store.mu.Lock()
	store.currentSession = &session
	store.mu.Unlock()
}

// GetCurrentUser retrieves the currently logged in user ID
func GetCurrentUser() uuid.UUID {
	if store == nil {
		return uuid.Nil
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	if store.currentSession == nil {
		return uuid.Nil
	}

	return store.currentSession.UserID
}

// IsLoggedIn checks if any user is currently logged in
func IsLoggedIn() bool {
	return GetCurrentUser() != uuid.Nil
}

// ClearCurrentUser clears the current user session
func ClearCurrentUser() {
	if store == nil || store.DB == nil {
		return
	}

	store.mu.RLock()
	session := store.currentSession
	store.mu.RUnlock()

	if session != nil {
		// Delete the session from the database
		err := store.DB.DeleteSession(context.Background(), session.ID)
		if err != nil {
			log.Printf("Error deleting session: %v", err)
		}
	}

	// Clear the cached session
	store.mu.Lock()
	store.currentSession = nil
	store.mu.Unlock()
}

// ClearAllSessions removes all sessions from the database
// Typically used for testing or when you need to force logout all users
func ClearAllSessions() error {
	if store == nil || store.DB == nil {
		return nil
	}

	err := store.DB.DeleteAllSessions(context.Background())
	if err != nil {
		return err
	}

	// Clear the cached session
	store.mu.Lock()
	store.currentSession = nil
	store.mu.Unlock()

	return nil
}
