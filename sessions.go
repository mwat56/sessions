/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type (
	// `tSessionData` stores the session data.
	tSessionData map[string]interface{}

	// TSession is a `map`-based session store
	TSession struct {
		sData *tSessionData
		sID   string
	}
)

// ChangeID generates a new SID for the current session's data.
func (so *TSession) ChangeID() (*TSession, error) {
	return sessionHandler.ChangeID(so.sID)
} // ChangeID()

// Delete removes the session data identified by `aKey`.
func (so *TSession) Delete(aKey string) (*TSession, error) {
	delete(*so.sData, aKey)

	return so, nil
} // Delete()

// Destroy a session.
//
// All internal references and external session files are removed.
func (so *TSession) Destroy() error {
	go sessionHandler.Destroy(so.sID)
	so.sData, so.sID = nil, ""

	return nil
} // Destroy()

// Get returns the session data identified by `aKey`.
//
// If `aKey` doesn't exist the method returns `nil`.
func (so *TSession) Get(aKey string) interface{} {
	if result, ok := (*so.sData)[aKey]; ok {
		return result
	}

	return nil
} // Get()

// SessionID returns the session's ID.
func (so *TSession) SessionID() string {
	return so.sID
} // SessionID()

// Set adds/updates the session data of `aKey` with `aValue`.
//
// This implementation always returns `nil`.
func (so *TSession) Set(aKey string, aValue interface{}) error {
	(*so.sData)[aKey] = aValue

	return nil
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newSession()` returns a new `TSession` instance.
func newSession(aSID string) *TSession {
	list := make(tSessionData)
	result := TSession{
		sData: &list,
		sID:   aSID,
	}

	return &result
} // newSession()

var (
	// `defaultLifetime` is the lifetime for an unused session.
	// It defaults to 1800 seconds (30 minutes).
	defaultLifetime = int64(60 * 30)

	// sessionHandler is the global session handler.
	sessionHandler *tSessionHandler

	// `sidName` is the GET/POST identifier fo the session ID.
	sidName = "SID"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// DefaultLifetime returns the default max. TTL of a session.
func DefaultLifetime() int64 {
	return defaultLifetime
} // DefaultLifetime()

// GetSession returns the `TSession` for `aRequest`.
//
// If `aRequest` doesn't provide a session ID in its form values
// a new (empty) session is returned.
//
// `aRequest` is the HTTP request received by the server.
func GetSession(aRequest *http.Request) *TSession {
	sid := aRequest.FormValue(sidName)
	if 0 == len(sid) {
		if c, err := aRequest.Cookie(sidName); nil == err {
			sid = c.Value
		} else {
			sid = newSID()
		}
	}
	result, _ := sessionHandler.Load(sid)

	return result
} // GetSession()

// NewSession returns a new `TSession` instance.
func NewSession() *TSession {
	result, _ := sessionHandler.Load(newSID())

	return result
} // NewSession()

// `newSID()` returns an ID based on time and random bytes.
func newSID() string {
	b := make([]byte, 16)
	rand.Read(b)
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), b)
	b = []byte(id[:24])

	return base64.URLEncoding.EncodeToString(b)
} // newSID()

// SetDefaultLifetime sets the default max. lifetime of a session.
//
// `aMaxLifetime` is the number of seconds a session's life lasts.
func SetDefaultLifetime(aMaxLifetime int64) {
	if 0 >= aMaxLifetime {
		defaultLifetime = 1800 // 1800 seconds = 30 minutes
	} else {
		defaultLifetime = aMaxLifetime
	}
} // SetDefaultLifetime()

// SetSIDname sets the session name.
//
// `aSID` identifies the session data.
func SetSIDname(aSID string) {
	if 0 < len(aSID) {
		sidName = aSID
	}
} // SetSIDname

// SIDname returns the configured session name.
//
// This name is expected to be used as a FORM field's name or
// the name of a CGI argument.
// Its default value is `SID`.
func SIDname() string {
	return sidName
} // SIDname()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Wrap initialises the session handling.
//
// `aHandler` responds to the actual HTTP request.
//
// `aSessionDir` is the name of the directory to store session files.
func Wrap(aHandler http.Handler, aSessionDir string) http.Handler {
	sessionHandler, _ = newFilehandler(aSessionDir)

	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			var usersession *TSession

			sid := aRequest.FormValue(sidName)
			if 0 == len(sid) {
				sid = newSID()
			}
			// load session file from disk
			usersession, _ = sessionHandler.Load(sid)

			// store a reference for other handlers
			c := http.Cookie{
				Name:  sidName,
				Value: sid,
			}
			aRequest.AddCookie(&c)

			// the original handler can access the session now
			aHandler.ServeHTTP(aWriter, aRequest)

			// save the (possibly modified) session data
			go sessionHandler.Store(usersession)
		})
} // Wrap()

/* _EoF_ */
