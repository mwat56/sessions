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
	"syscall"
	"time"
)

type (
	// TSession is an opaque session data store.
	TSession struct {
		sID    string
		sValue interface{} // used when requesting a data value
	}
)

// `changeID()` generates a new SID for the current session's data.
//
// Since the ID changes are handle internally by the `Wrap()` function
// this method is not exported but kept private.
func (so *TSession) changeID() *TSession {
	result := doRequest(shChangeSession, so.sID, "", nil)
	so.sID = result.sID

	return so
} // ChangeID()

// Delete removes the session data identified by `aKey`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) Delete(aKey string) *TSession {
	doRequest(shDeleteKey, so.sID, aKey, nil)

	return so
} // Delete()

// Destroy a session.
//
// All internal references and external session files are removed.
func (so *TSession) Destroy() {
	doRequest(shDestroySession, so.sID, "", nil)
	so.sID = ""
} // Destroy()

// Get returns the session data identified by `aKey`.
//
// If `aKey` doesn't exist the method returns `nil`.
//
//	`aKey` The identifier to lookup.
func (so *TSession) Get(aKey string) interface{} {
	result := doRequest(shGetKey, so.sID, aKey, nil)

	return result.sValue
} // Get()

// ID returns the session's ID.
func (so *TSession) ID() string {
	return so.sID
} // ID()

// Len returns the current length of the list of session vars.
func (so *TSession) Len() int {
	result := doRequest(shSessionLen, so.sID, "", nil)
	if len, ok := result.sValue.(int); ok {
		return len
	}

	return 0
} // Len()

// Set adds/updates the session data of `aKey` with `aValue`.
//
//	`aKey` The identifier to lookup.
//	`aValue` The value to assign.
func (so *TSession) Set(aKey string, aValue interface{}) *TSession {
	doRequest(shSetKey, so.sID, aKey, aValue)

	return so
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// The channel to send requests through to `sessionMonitor()`
	// (`tShRequest` defined in `monitor.go`).
	chSession = make(chan tShRequest, 64)
)

// `doRequest()` queries the session manager for certain data.
//
//	`aType` The lookup type.
//	`aSID` The session ID to use for lookup.
//	`aKey` Optional session variable name/key.
//	`aValue` Optional session variable value.
func doRequest(aType tShLookupType, aSID, aKey string, aValue interface{}) (rSession *TSession) {
	answer := make(chan *TSession)
	defer close(answer)

	request := tShRequest{
		rKey:   aKey,
		rSID:   aSID,
		rType:  aType,
		rValue: aValue,
		reply:  answer,
	}
	chSession <- request
	rSession = <-answer

	return
} // doRequest()

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
			return doRequest(shNewSession, newSID(), "", nil)
		}
	}

	return doRequest(shLoadSession, sid, "", nil)
} // GetSession()

/*
// NewSession returns a new (empty) session.
func NewSession() *TSession {
	return doRequest(shNewSession, newSID(), "", nil)
} // NewSession()
*/

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
	if 0 >= aTTL {
		sessionTTL = 600 // 600 seconds == 10 minutes
	} else {
		sessionTTL = aTTL
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
	if dir, err := checkSessionDir(aSessionDir); nil != err {
		log.Fatalf("%s: %v", os.Args[0], err)
	} else {
		go sessionMonitor(dir, chSession)
	}

	return http.HandlerFunc(
		func(aWriter http.ResponseWriter, aRequest *http.Request) {
			defer func() {
				// make sure a `panic` won't kill the program
				if err := recover(); err != nil {
					log.Printf("[%v] caught panic: %v", aRequest.RemoteAddr, err)
				}
			}()

			sid := aRequest.FormValue(string(sidName))
			if 0 == len(sid) {
				sid = newSID()
			}

			// prepare a reference for `GetSession()`
			ctx := context.WithValue(aRequest.Context(), sidName, sid)
			aRequest = aRequest.WithContext(ctx)

			// load session file from disk
			session := doRequest(shLoadSession, sid, "", nil)

			// replace the old SID by a new ID
			sid = session.changeID().ID()

			// keep a session reference with the writer
			hr := &tHRefWriter{aWriter, sid}

			// the original handler can access the session now
			aHandler.ServeHTTP(hr, aRequest)

			// save the possibly updated session data
			doRequest(shStoreSession, sid, "", nil)
		})
} // Wrap()

/* _EoF_ */
