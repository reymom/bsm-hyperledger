package sessionstore

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
)

var Store *sessions.FilesystemStore

func Init() error {
	// authKeyOne := securecookie.GenerateRandomKey(64)
	// encryptionKeyOne := securecookie.GenerateRandomKey(32)

	// Store = sessions.NewFilesystemStore(
	// 	"",
	// 	authKeyOne,
	// 	encryptionKeyOne,
	// )

	Store = sessions.NewFilesystemStore("", []byte("it-was-me-who-killed-elvis"))

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(connection.Login{})

	return nil
}

func GetLoginFromSession(s *sessions.Session) *connection.Login {
	val := s.Values["login"]
	var user = connection.Login{}
	user, _ = val.(connection.Login)
	return &user
}
