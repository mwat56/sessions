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
	// tSessionData stores the session data.
	tSessionData map[string]interface{}

	// TSession is the type to get/set session data.
	TSession interface {

		// Data returns the complete set of session data.
		Data() *tSessionData

		// Get returns the session data identified by `aKey`.
		Get(aKey string) interface{}

		// Set adds/updates the session data of `aKey` with `aValue`.
		Set(aKey string, aValue interface{}) error

		// Delete removes the session data identified by `aKey`.
		Delete(aKey string) error

		// SessionID returns the session's ID.
		SessionID() string
	}

	// TSessionHandler defines the interface of a session handler.
	TSessionHandler interface {

		// Init (initialise) the session.
		//
		// `aSavePath` The path where to store/retrieve the session data.
		Init(aSavePath string) error

		// Close the session.
		Close() error

		// Load reads the session data from disk.
		//
		// `aSID` The session ID to read data for.
		Load(aSID string) (*TSession, error)

		// Store writes session data to disk..
		//
		// `aSID` The current session ID.
		// `aValue` The session data to store.
		Store(aSession *TSession) error

		// Destroy a session.
		//
		// `aSID` The session ID being destroyed.
		Destroy(aSID string) error

		// GC cleans up old sessions.
		//
		// `aMaxlifetime` Sessions that have not updated for the
		// last `aMaxlifetime` seconds will be removed.
		GC(aMaxlifetime int64) error
	}
)

var (
	// Sessions is the global session handler.
	Sessions *TFileSessionHandler
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// DefaultLifetime returns the default max. lifetime of a session.
func DefaultLifetime() int64 {
	return defaultLifetime
} // DefaultLifetime()

// `newID()` returns an ID based on time and random bytes.
func newID() string {
	// var r string
	b := make([]byte, 16)
	if _, err := rand.Read(b); nil == err {
		// r = string(b)
	}
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), b)
	b = []byte(id[:24])
	/*
		id := strconv.FormatInt(time.Now().UnixNano(), 10) + r
		for 32 > len(id) {
			id += id
		}
		b = []byte(id[:32])
	*/
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

			// (1) get `aSID` from `aRequest`
			// store `aSID` as `aRequest.Header["Cookie"]` (or "SID"?)
			if sid := aRequest.FormValue("SID"); 0 < len(sid) {
				// load session file
			}

			// (2) init session handling for `aSID`
			// provide session access (how?)

			// (3) call wrapped `ServeHttp()` method
			aHandler.ServeHTTP(aWriter, aRequest)

			// (4) close session for `aSID`

		})
} // Wrap()

/* _EoF_ */
