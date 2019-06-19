/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"encoding/gob"
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

	// TFileSessionHandler implements a file-based session handler.
	TFileSessionHandler struct {
		fsDir         string        // directory to store session files
		fsList        tSessionList  // list of known sessions
		fsMaxLifetime int64         // lifetime seconds
		fsMtx         *sync.RWMutex // guard file operations
	}
)

// Close the session.
func (fs *TFileSessionHandler) Close() error {
	return fs.GC(0)
} // Close()

// Destroy a session.
//
// `aSID` The session ID being destroyed.
func (fs *TFileSessionHandler) Destroy(aSID string) error {
	// If m is nil or there is no such element, delete is a no-op.
	delete(fs.fsList, aSID)

	file := filepath.Join(fs.fsDir, aSID) + ".sid"
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
func (fs *TFileSessionHandler) GC(aMaxlifetime int64) error {
	if 0 == aMaxlifetime {
		aMaxlifetime = fs.fsMaxLifetime * int64(time.Second)
	} else {
		aMaxlifetime *= int64(time.Second)
	}

	expired := time.Now().Add(time.Duration(0 - aMaxlifetime))
	files, err := filepath.Glob(fs.fsDir + "/*.sid")
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
			if 5 < len(fName) {
				sid := fName[:len(fName)-4]
				delete(fs.fsList, sid)
			}
			os.Remove(file)
		}
	}

	return nil
} // GC()

// Init (initialise) the session.
//
// `aSavePath` The path where to store/retrieve the session data.
func (fs *TFileSessionHandler) Init(aSavePath string) error {
	return fs.setSessionDir(aSavePath)
} // Open()

// `load()` reads the session data for `aSID` from disk.
func (fs *TFileSessionHandler) load(aSID string) (*TMapSession, error) {
	// locking is done by the caller
	result := newSession(aSID)

	if data, ok := fs.fsList[aSID]; ok {
		result.dmData = *data

		return result, nil
	}

	fName := filepath.Join(fs.fsDir, aSID) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDONLY, 0)
	if nil != err {
		if os.IsNotExist(err) {
			err = nil
		}
		fs.fsList[aSID] = &result.dmData

		return result, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	gob.Register(result.dmData)
	gob.Register(time.Now())
	var store tStoreStruct
	err = decoder.Decode(&store)
	if e, ok := store["expires"]; ok {
		if expires, ok := e.(time.Time); ok && expires.After(time.Now()) {
			if id, ok := store["sid"]; ok {
				if key, ok := id.(string); ok && (key == aSID) {
					if d, ok := store["data"]; ok {
						if data, ok := d.(tSessionData); ok {
							result.dmData = data
							err = nil
						}
					}
				}
			}
		}
	}
	fs.fsList[aSID] = &result.dmData

	return result, err
} // load()

// Load reads the session data from disk.
//
// `aSID` The session ID to read data for.
func (fs *TFileSessionHandler) Load(aSID string) (*TMapSession, error) {
	fs.fsMtx.Lock()
	defer fs.fsMtx.Unlock()

	return fs.load(aSID)
} // Load()

// `setSessionDir()` assigns the directory to store the session files.
func (fs *TFileSessionHandler) setSessionDir(aSavePath string) error {
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
		return fmt.Errorf("Not as directory: %q", dir)
	}
	fs.fsDir = dir

	return nil
} // setSessionDir()

// `store()` saves the session data on disk.
func (fs *TFileSessionHandler) store(aSession *TMapSession) error {
	// locking is done by the caller
	sid := aSession.dmID
	fs.fsList[sid] = &aSession.dmData
	fName := filepath.Join(fs.fsDir, sid) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if nil != err {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	expires := time.Now().Add(time.Duration(fs.fsMaxLifetime)*time.Second + 1)
	store := &tStoreStruct{
		"data":    aSession.dmData,
		"expires": expires,
		"sid":     sid,
	}
	gob.Register(aSession.dmData)
	gob.Register(expires)

	return encoder.Encode(store)
} // store()

// Store writes session data to disk.
//
// `aSID` The current session ID.
// `aValue` The session data to store.
func (fs *TFileSessionHandler) Store(aSession *TMapSession) error {
	fs.fsMtx.Lock()
	defer fs.fsMtx.Unlock()

	return fs.store(aSession)
} // Store()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// newFilehandler returns a new `TFileSessionHandler` instance.
//
// `aSavePath` is the directory to use for storing sessions files.
func newFilehandler(aSavePath string) (*TFileSessionHandler, error) {
	result := TFileSessionHandler{
		fsList:        make(tSessionList, 32),
		fsMaxLifetime: defaultLifetime,
	}
	err := result.setSessionDir(aSavePath)

	return &result, err
} // NewFilehandler()

/* _EoF_ */
