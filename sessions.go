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

var (
	// `defaultLifetime` is the lifetime for an unused session.
	// It defaults to 30 minutes.
	defaultLifetime = int64(60 * 30)
)

type (
	// tTSession is the type to get/set session data.
	tTSession interface {

		// Get returns the session data identified by `aKey`.
		Get(aKey string) interface{}

		// Set adds/updates the session data of `aKey` with `aValue`.
		Set(aKey string, aValue interface{}) error

		// Delete removes the session data identified by `aKey`.
		Delete(aKey string) error

		// SessionID returns the session's ID.
		SessionID() string
	}

	// tTSessionHandler defines the interface of a session handler.
	tTSessionHandler interface {

		// Init (initialise) the session.
		//
		// `aSavePath` The path where to store/retrieve the session data.
		Init(aSavePath string) error

		// Close the session.
		Close() error

		// Load reads the session data from disk.
		//
		// `aSID` The session ID to read data for.
		Load(aSID string) (*tTSession, error)

		// Store writes session data to disk..
		//
		// `aSID` The current session ID.
		// `aValue` The session data to store.
		Store(aSession *tTSession) error

		// Destroy a session.
		//
		// `aSID` The session ID being destroyed.
		Destroy(aSID string) error

		// GC cleans up old sessions.
		GC() error
	}
)

var (
	// Sessions is the global session handler.
	Sessions *TSessionHandler
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// DefaultLifetime returns the default max. lifetime of a session.
func DefaultLifetime() int64 {
	return defaultLifetime
} // DefaultLifetime()

// `newID()` returns an ID based on time and random bytes.
func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), b)
	b = []byte(id[:24])

	return base64.URLEncoding.EncodeToString(b)
} // newID()

// SetDefaultLifetime sets the default max. lifetime of a session.
//
// `aMaxLifetime` is the number of seconds a session's life lasts.
func SetDefaultLifetime(aMaxLifetime int64) {
	if 0 >= aMaxLifetime {
		defaultLifetime = int64(time.Minute)
	} else {
		defaultLifetime = int64(time.Second) * aMaxLifetime
	}
} // SetDefaultLifetime()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Wrap initialises the session handling.
//
// `aHandler` responds to the actual HTTP request.
//
// `aSessionDir` is the name of the directory to store session files.
func Wrap(aHandler http.Handler, aSessionDir string) http.Handler {
	Sessions, _ = newFilehandler(aSessionDir)

	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			var usersession *TSession

			// (1) get `aSID` from `aRequest`
			// store `aSID` as `aRequest.Header["Cookie"]` (or "SID"?)
			if sid := aRequest.FormValue("SID"); 0 < len(sid) {
				// load session file
				usersession, _ = Sessions.Load(sid)
			}

			// (2) init session handling for `aSID`
			// provide session access (how?)

			// (3) call wrapped `ServeHttp()` method
			aHandler.ServeHTTP(aWriter, aRequest)

			// (4) close session for `aSID`
			Sessions.Store(usersession)

		})
} // Wrap()

/* _EoF_ */
