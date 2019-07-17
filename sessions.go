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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type (
	// TSession is an opaque session data store.
	TSession struct {
		sID    string
		sValue interface{} // used only when requesting a data value
	}
)

// `changeID()` generates a new SID for the current session's data.
//
// Since the ID changes are handle internally by the `Wrap()` function
// this method is not exported but kept private.
func (so *TSession) changeID() *TSession {
	result := so.request(smChangeSession, "", nil)
	so.sID = result.sID

	return so
} // ChangeID()

// Delete removes the session data identified by `aKey`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) Delete(aKey string) *TSession {
	so.request(smDeleteKey, aKey, nil)

	return so
} // Delete()

// Destroy a session.
//
// All internal references and external session files are removed.
func (so *TSession) Destroy() {
	so.request(smDestroySession, "", nil)
	so.sID = ""
} // Destroy()

// EmptySession returns whether the current session has no associated data.
func (so *TSession) EmptySession() bool {
	result := so.request(smEmptySession, "", nil)
	if b, ok := result.sValue.(bool); ok {
		return b
	}

	return false
} // EmptySession()

// Get returns the session data identified by `aKey`.
//
// If `aKey` doesn't exist the method returns `nil`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) Get(aKey string) interface{} {
	result := so.request(smGetKey, aKey, nil)

	return result.sValue
} // Get()

// GetBool returns the `boolean` session data identified by `aKey`.
//
// The second (`bool`) return value signals whether a session
// value of type `bool` is associated with `aKey`.
//
// If `aKey` doesn't exist the method returns `false`
// and `false`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) GetBool(aKey string) (bool, bool) {
	result := so.request(smGetKey, aKey, nil)
	if b, ok := result.sValue.(bool); ok {
		return b, true
	}

	return false, false
} // GetBool()

// GetFloat returns the `float64` session data identified by `aKey`.
//
// The second (`bool`) return value signals whether a session
// value of type `float64` is associated with `aKey`.
//
// If `aKey` doesn't exist the method returns `0` (zero)
// and `false`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) GetFloat(aKey string) (float64, bool) {
	result := so.request(smGetKey, aKey, nil)
	if f, ok := result.sValue.(float64); ok {
		return f, true
	}
	if f, ok := result.sValue.(float32); ok {
		return float64(f), true
	}

	return 0, false
} // GetFloat()

// GetInt returns the `int` session data identified by `aKey`.
//
// The second (`bool`) return value signals whether a session
// value of type `int` is associated with `aKey`.
//
// If `aKey` doesn't exist the method returns `0` (zero)
// and `false`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) GetInt(aKey string) (int64, bool) {
	result := so.request(smGetKey, aKey, nil)
	if i, ok := result.sValue.(int64); ok {
		return i, true
	}
	if i, ok := result.sValue.(int); ok {
		return int64(i), true
	}

	return 0, false
} // GetInt()

// GetString returns the `string` session data identified by `aKey`.
//
// The second (`bool`) return value signals whether a session
// value of type `string` is associated with `aKey`.
//
// If `aKey` doesn't exist the method returns an empty string
// and `false`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) GetString(aKey string) (string, bool) {
	result := so.request(smGetKey, aKey, nil)
	if str, ok := result.sValue.(string); ok {
		return str, true
	}

	return "", false
} // GetString()

// GetTime returns the `time.Time` session data identified by `aKey`.
//
// The second (`bool`) return value signals whether a session
// value of type `time.Time` is associated with `aKey`.
//
// If `aKey` doesn't exist the method returns a zero time and `false`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) GetTime(aKey string) (rTime time.Time, rOK bool) {
	result := so.request(smGetKey, aKey, nil)
	if t, ok := result.sValue.(time.Time); ok {
		rTime, rOK = t, true
	}

	return
} // GetTime()

// ID returns the session's ID.
func (so *TSession) ID() string {
	return so.sID
} // ID()

// Len returns the current length of the list of session vars.
func (so *TSession) Len() int {
	result := so.request(smSessionLen, "", nil)
	if len, ok := result.sValue.(int); ok {
		return len
	}

	return 0
} // Len()

// `request()` queries the session manager for certain data.
//
//	`aType` The lookup type.
//	`aKey` Optional session variable name/key.
//	`aValue` Optional session variable value.
func (so *TSession) request(aType tShLookupType, aKey string, aValue interface{}) (rSession *TSession) {
	answer := make(chan *TSession)
	defer close(answer)

	chSession <- tShRequest{
		rKey:   aKey,
		rSID:   so.sID,
		rType:  aType,
		rValue: aValue,
		reply:  answer,
	}
	rSession = <-answer

	return
} // request()

// Set adds/updates the session data of `aKey` with `aValue`.
//
//	`aKey` The identifier to lookup.
//	`aValue` The value to assign.
func (so *TSession) Set(aKey string, aValue interface{}) *TSession {
	so.request(smSetKey, aKey, aValue)

	return so
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// The channel to send requests through to `goMonitor()`
	// (`tShRequest` defined in `monitor.go`).
	chSession = make(chan tShRequest, 2)
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
	so := &TSession{sID: sid}

	return so.request(smLoadSession, "", nil)
} // GetSession()

// `newSID()` returns an ID based on time and random bytes.
func newSID() string {
	b := make([]byte, 16)
	rand.Read(b)
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), b)
	b = []byte(id[:24])

	return base64.URLEncoding.EncodeToString(b)
} // newSID()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// `sessionTTL` is the max. TTL for an unused session.
	// It defaults to 600 seconds (10 minutes).
	sessionTTL = 600
)

// SessionTTL returns the Time-To-Life of a session (in seconds).
func SessionTTL() int {
	return sessionTTL
} // SessionTTL()

// SetSessionTTL sets the lifetime of a session.
//
// `aTTL` is the number of seconds a session's life lasts.
func SetSessionTTL(aTTL int) {
	if 0 < aTTL {
		sessionTTL = aTTL
	} else {
		sessionTTL = 600 // 600 seconds == 10 minutes
	}
} // SetSessionTTL()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// `tSIDname` is a string that is not a string
	// (builtin types should not be used as `Context` keys).
	tSIDname string
)

var (
	// `sidName` is the GET/POST identifier fo the session ID.
	sidName = tSIDname("SID")
)

// SetSIDname sets the name of the session ID.
//
// `aSID` identifies the session data.
func SetSIDname(aSID string) {
	if 0 < len(aSID) {
		sidName = tSIDname(aSID)
	}
} // SetSIDname

// SIDname returns the configured session name.
//
// This name is expected to be used as a FORM field's name or the
// name of a CGI argument.
// Its default value is `SID`.
func SIDname() string {
	return string(sidName)
} // SIDname()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `checkSessionDir()` checks whether `aSessionDir` exists and
// creates it if neccessary.
//
// This function is a helper of and called by `Wrap()`.
func checkSessionDir(aSessionDir string) (rDir string, rErr error) {
	if rDir, rErr = filepath.Abs(aSessionDir); nil != rErr {
		return
	}
	if fi, err := os.Stat(rDir); nil != err {
		if e, ok := err.(*os.PathError); ok && syscall.ENOENT == e.Err {
			fmode := os.ModeDir | 0775
			rErr = os.MkdirAll(filepath.FromSlash(rDir), fmode)
		} else {
			rErr = err
		}
	} else if !fi.IsDir() {
		rErr = fmt.Errorf("Not a directory: %q", rDir)
	}

	return
} // checkSessionDir()

// Wrap initialises the session handling.
//
//	`aHandler` responds to the actual HTTP request.
//	`aSessionDir` is the name of the directory to store session files.
func Wrap(aHandler http.Handler, aSessionDir string) http.Handler {
	var doOnce sync.Once
	doOnce.Do(func() {
		if dir, err := checkSessionDir(aSessionDir); nil != err {
			log.Fatalf("%s: %v", os.Args[0], err)
		} else {
			go goMonitor(dir, chSession)
		}
	})

	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			switch aRequest.Method {
			case "GET", "POST":
				session := &TSession{
					sID: aRequest.FormValue(string(sidName)),
				}
				if 0 < len(session.sID) {
					// load session file from disk
					session.request(smLoadSession, "", nil)

					// replace the old SID by a new ID
					session.changeID()
				} else {
					session.sID = newSID()
				}

				// prepare a reference for `GetSession()`
				ctx := context.WithValue(aRequest.Context(), sidName, session.sID)
				aRequest = aRequest.WithContext(ctx)

				// keep a session reference with the writer
				hr := &tHRefWriter{
					aWriter,
					session.sID,
				}

				// the original handler can access the session now
				aHandler.ServeHTTP(hr, aRequest)

				// save the possibly updated session data
				session.request(smStoreSession, "", nil)

			default:
				// run the original handler
				aHandler.ServeHTTP(aWriter, aRequest)
			}
		})
} // Wrap()

/* _EoF_ */
