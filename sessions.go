/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"context"
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

// Len returns the current length of the list of session vars.
func (so *TSession) Len() int {
	return len(*so.sData)
} // Len()

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
	// `sessionTTL` is the max. TTL for an unused session.
	// It defaults to 600 seconds (10 minutes).
	sessionTTL = 600

	// `sessionHandler` is the global session handler.
	sessionHandler *tSessionHandler
)

type (
	// `tSIDname` is a string that is not a string
	// (builtin types should not be used as `Context` keys).
	tSIDname string
)

var (
	// `sidName` is the GET/POST identifier fo the session ID.
	sidName = tSIDname("SID")
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// GetSession returns the `TSession` for `aRequest`.
//
// If `aRequest` doesn't provide a session ID in its form values
// a new (empty) `TSession` instance is returned.
//
// `aRequest` is the HTTP request received by the server.
func GetSession(aRequest *http.Request) *TSession {
	sid := aRequest.FormValue(string(sidName))
	if 0 == len(sid) {
		ctx := aRequest.Context()
		if id, ok := ctx.Value(sidName).(tSIDname); ok {
			sid = string(id)
		} else {
			sid = newSID()
		}
	}
	result, _ := sessionHandler.Load(sid)

	return result
} // GetSession()

// `newSID()` returns an ID based on time and random bytes.
func newSID() string {
	b := make([]byte, 16)
	rand.Read(b)
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), b)
	b = []byte(id[:24])

	return base64.URLEncoding.EncodeToString(b)
} // newSID()

// SessionTTL returns the Time-To-Life of a session (in seconds).
func SessionTTL() int {
	return sessionTTL
} // SessionTTL()

// SetSessionTTL sets the default max. lifetime of a session.
//
// `aTTL` is the number of seconds a session's life lasts.
func SetSessionTTL(aTTL int) {
	if 0 >= aTTL {
		sessionTTL = 600 // 600 seconds == 10 minutes
	} else {
		sessionTTL = aTTL
	}
} // SetSessionTTL()

// SetSIDname sets the session name.
//
// `aSID` identifies the session data.
func SetSIDname(aSID string) {
	if 0 < len(aSID) {
		sidName = tSIDname(aSID)
	}
} // SetSIDname

// SIDname returns the configured session name.
//
// This name is expected to be used as a FORM field's name or
// the name of a CGI argument.
// Its default value is `SID`.
func SIDname() string {
	return string(sidName)
} // SIDname()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Wrap initialises the session handling.
//
// `aHandler` responds to the actual HTTP request.
//
// `aSessionDir` is the name of the directory to store session files.
func Wrap(aHandler http.Handler, aSessionDir string) http.Handler {
	sessionHandler, _ = newSessionHandler(aSessionDir)

	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			var usersession *TSession

			sid := aRequest.FormValue(string(sidName))
			if 0 == len(sid) {
				sid = newSID()
			}
			// store a reference for `GetSession()`
			ctx := context.WithValue(aRequest.Context(), sidName, sid)
			aRequest = aRequest.WithContext(ctx)

			// load session file from disk
			usersession, _ = sessionHandler.Load(sid)

			// the original handler can access the session now
			aHandler.ServeHTTP(aWriter, aRequest)

			// save the updated session data
			go sessionHandler.Store(usersession)
		})
} // Wrap()

/* _EoF_ */
