package sessionstore

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
)

type SessionStore struct {
	Store            *sessions.FilesystemStore
	Login            *connection.Login
	NetworkContracts connection.NetworkContract
}

func NewSessionStore(storeKey string) (*SessionStore, error) {
	store := sessions.NewFilesystemStore("", []byte(storeKey))
	gob.Register(connection.Login{})

	return &SessionStore{
		Store:            store,
		Login:            &connection.Login{},
		NetworkContracts: make(connection.NetworkContract),
	}, nil
}

func getLoginFromSession(s *sessions.Session) connection.Login {
	val := s.Values["login"]
	user := connection.Login{}
	user, _ = val.(connection.Login)
	return user
}

func (s *SessionStore) CheckLoginFromSession(r *http.Request, loginsMap connection.UsersLoginMap) (bool, error) {
	loggedIn := false
	networkContracts := make(connection.NetworkContract)
	session, e := s.Store.Get(r, "auth-session")
	if e != nil {
		return false, e
	}

	loginSession := getLoginFromSession(session)
	if (connection.Login{}) != loginSession && connection.IsRegistered(loginsMap, loginSession.Name, loginSession.Password) {
		loggedIn = true
		s.Login = &loginSession

		_, networkContracts, e = connection.GetGatewayObjects(s.Login.Name)
		if e != nil {
			log.Err(e).Msg("Error getting gateway objects")
		}
	} else {
		s.Login = new(connection.Login)
	}

	s.NetworkContracts = networkContracts
	return loggedIn, nil
}
