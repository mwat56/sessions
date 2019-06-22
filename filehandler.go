/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type (
	// List of known sessions
	tSessionList map[string]*tSessionData

	// Structure to store the session data:
	//
	//	tStoreStruct{
	//		"data":    tSessionData,
	//		"expires": expiration date,
	//		"sid":     aSID,
	//	}
	//
	tStoreStruct map[string]interface{}

	// TSessionHandler implements a file-based session handler.
	TSessionHandler struct {
		shDir      string        // directory to store session files
		shList     tSessionList  // list of known sessions
		shLifetime int64         // lifetime seconds
		shMtx      *sync.RWMutex // guard file operations
	}
)

// `changeID()` moves the data associated with `aOldSID` to `aNewSID`.
func (sh *TSessionHandler) changeID(aOldSID, aNewSID string) (*TSession, error) {
	// locking is done be the caller

	if data, ok := sh.shList[aOldSID]; ok {
		sh.shList[aNewSID] = data
		delete(sh.shList, aOldSID)
		result := TSession{
			sData: *data,
			sID:   aNewSID,
		}

		return &result, nil
	}

	return nil, errors.New("Session ID '" + aOldSID + "' not found")
} // changeSID()

// ChangeID generates a new session ID for the data associated with `aSID`.
func (sh *TSessionHandler) ChangeID(aSID string) (*TSession, error) {
	sh.shMtx.Lock()
	defer sh.shMtx.Unlock()

	newSID := newID()
	if result, err := sh.changeID(aSID, newSID); nil == err {
		return result, nil
	}

	// Reaching this point of execution means the session `aSID`
	// hasn't been loaded (or found) yet.
	_, err := sh.load(aSID)
	if nil == err {
		return sh.changeID(aSID, newSID)
	}

	return nil, err
} // ChangeID()

// Close the session.
func (sh *TSessionHandler) Close() error {
	return sh.GC()
} // Close()

// Destroy a session.
//
// `aSID` The session ID being destroyed.
func (sh *TSessionHandler) Destroy(aSID string) error {
	sh.shMtx.Lock()
	defer sh.shMtx.Unlock()

	// If m is nil or there is no such element, delete is a no-op.
	delete(sh.shList, aSID)

	file := filepath.Join(sh.shDir, aSID) + ".sid"
	if _, err := os.Stat(file); nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
		return err
	}

	return os.Remove(file)
} // Destroy()

// GC cleans up old sessions.
//
// `aMaxlifetime` Sessions that have not updated for the
// last `aMaxlifetime` seconds will be removed.
func (sh *TSessionHandler) GC() error {
	ttl := sh.shLifetime * int64(time.Second)
	expired := time.Now().Add(time.Duration(0 - ttl))
	files, err := filepath.Glob(sh.shDir + "/*.sid")
	if nil != err {
		return err
	}
	for _, file := range files {
		fi, err := os.Stat(file)
		if nil != err {
			continue
		}
		if fi.ModTime().Before(expired) {
			fName := filepath.Base(file)
			sid := fName[:len(fName)-4]
			sh.Destroy(sid)
		}
	}

	return nil
} // GC()

// Init initialises the session.
//
// `aSavePath` The path where to store/retrieve the session data.
func (sh *TSessionHandler) Init(aSavePath string) error {
	return sh.initSessionDir(aSavePath)
} // Init()

func (sh *TSessionHandler) initSessionDir(aSavePath string) error {
	dir, err := filepath.Abs(aSavePath)
	if nil != err {
		return err
	}
	if fi, err := os.Stat(dir); nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			fmode := os.ModeDir | 0775
			if err := os.MkdirAll(filepath.FromSlash(dir), fmode); nil != err {
				return err
			}
		} else {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("Not a directory: %q", dir)
	}
	sh.shDir = dir

	return nil
} // Init()

// `load()` reads the session data for `aSID` from disk.
func (sh *TSessionHandler) load(aSID string) (*TSession, error) {
	// locking is done by the caller
	result := newSession(aSID)

	if data, ok := sh.shList[aSID]; ok {
		result.sData = *data

		return result, nil
	}

	fName := filepath.Join(sh.shDir, aSID) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDONLY, 0)
	if nil != err {
		if os.IsNotExist(err) {
			err = nil
		}
		sh.shList[aSID] = &result.sData

		return result, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	gob.Register(result.sData)
	gob.Register(time.Now())
	var ss tStoreStruct
	err = decoder.Decode(&ss)
	if e, ok := ss["expires"]; ok {
		if expires, ok := e.(time.Time); ok && expires.After(time.Now()) {
			if id, ok := ss["sid"]; ok {
				if key, ok := id.(string); ok && (key == aSID) {
					if d, ok := ss["data"]; ok {
						if data, ok := d.(tSessionData); ok {
							result.sData = data
							err = nil
						}
					}
				}
			}
		}
	}
	sh.shList[aSID] = &result.sData

	return result, err
} // load()

// Load reads the session data from disk.
//
// `aSID` The session ID to read data for.
func (sh *TSessionHandler) Load(aSID string) (*TSession, error) {
	sh.shMtx.Lock()
	defer sh.shMtx.Unlock()

	return sh.load(aSID)
} // Load()

// `store()` saves the session data on disk.
func (sh *TSessionHandler) store(aSession *TSession) error {
	// locking is done by the caller
	sid := aSession.sID
	sh.shList[sid] = &aSession.sData
	fName := filepath.Join(sh.shDir, sid) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if nil != err {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	expires := time.Now().Add(time.Duration(sh.shLifetime)*time.Second + 1)
	store := &tStoreStruct{
		"data":    aSession.sData,
		"expires": expires,
		"sid":     sid,
	}
	gob.Register(aSession.sData)
	gob.Register(expires)

	return encoder.Encode(store)
} // store()

// Store writes session data to disk.
//
// `aSID` The current session ID.
// `aValue` The session data to store.
func (sh *TSessionHandler) Store(aSession *TSession) error {
	sh.shMtx.Lock()
	defer sh.shMtx.Unlock()

	return sh.store(aSession)
} // Store()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// newFilehandler returns a new `TFileSessionHandler` instance.
//
// `aSavePath` is the directory to use for storing sessions files.
func newFilehandler(aSavePath string) (*TSessionHandler, error) {
	result := TSessionHandler{
		shList:     make(tSessionList, 32),
		shLifetime: defaultLifetime,
	}
	err := result.Init(aSavePath)

	return &result, err
} // NewFilehandler()

/* _EoF_ */
