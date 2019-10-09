/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides functions to monitor all session requests
 * (i.e. reading/writing session data).
 */

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"
)

type (
	// `tSessionData` stores the session data.
	tSessionData map[string]interface{}

	// `tShList` is the list of known sessions
	tShList map[string]*tSessionData

	// `tShLookupType` is the kind of request to `goMonitor()`
	tShLookupType int

	// `tShRequest` is the request structure channelled to `goMonitor()`
	tShRequest struct {
		rKey   string
		rSID   string
		rType  tShLookupType
		rValue interface{}
		reply  chan *TSession
	}

	// Structure to store the session data:
	//
	//	tStoreStruct{
	//		"data":    tSessionData,
	//		"expires": Unix.secs,
	//		"sid":     aSID,
	//	}
	//
	tStoreStruct map[string]interface{}
)

const (
	// Possible request types send to SessionMonitor `goMonitor()`
	smTerminate = tShLookupType(1 << iota) // for testing only: terminate `goMonitor()`
	smChangeSession
	smDeleteKey
	smDestroySession
	smGetKey
	smLoadSession
	smSessionLen
	smSetKey
	smStoreSession
)

// `goDel()` deletes the file and session data for `aSID`.
//
// This function is called from `goGC()`
//
// We need this additional level of indirection to delete both, the
// session data in memory and the session file.
func goDel(aSID string) {
	answer := make(chan *TSession)
	defer close(answer)

	soSessionChannel <- tShRequest{
		rSID:  aSID,
		rType: smDestroySession,
		reply: answer,
	}
	<-answer // ignore the result
} // goDel()

// `goGC()` cleans up old sessions.
//
// Sessions that have not been updated for at least
// `SessionTTL()` seconds will be removed.
//
//	`aSessionDir` The directory where the session files are stored.
func goGC(aSessionDir string) {
	secs := time.Now().Unix() - int64(soSessionTTL)
	expired := time.Unix(secs, 0)
	files, err := filepath.Glob(aSessionDir + "/*.sid")
	if nil != err {
		return
	}
	for _, file := range files {
		fi, err := os.Stat(file)
		if nil != err {
			continue
		}
		if fi.ModTime().Before(expired) {
			fName := filepath.Base(file)
			go goDel(fName[:len(fName)-4])
		}
	}
} // goGC()

// `goMonitor()` handles the access to the internal list of session data.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aRequest` The channel to receive request through.
func goMonitor(aSessionDir string, aRequest <-chan tShRequest) {
	shList := make(tShList, 32) // list of active sessions
	go goGC(aSessionDir)        // cleanup old session files
	gcInterval := time.Duration(soSessionTTL<<1)*time.Second + 1
	timer := time.NewTimer(gcInterval)
	defer timer.Stop()

	for { // wait for requests
		select {
		case request, more := <-aRequest:
			if !more { // channel closed
				return
			}

			switch request.rType {
			case smChangeSession:
				newsid := newSID()
				if data, ok := shList[request.rSID]; ok {
					shList[newsid] = data
					delete(shList, request.rSID)
				} else {
					list := make(tSessionData)
					shList[newsid] = &list
				}
				go goRemove(aSessionDir, request.rSID)
				request.reply <- &TSession{sID: newsid}

			case smDeleteKey:
				if data, ok := shList[request.rSID]; ok {
					delete(*data, request.rKey)
				}
				request.reply <- &TSession{sID: request.rSID}

			case smDestroySession:
				delete(shList, request.rSID)
				go goRemove(aSessionDir, request.rSID)
				request.reply <- &TSession{}

			case smGetKey:
				result := &TSession{
					sID: request.rSID,
				}
				data, ok := shList[request.rSID]
				if !ok {
					data = loadSession(aSessionDir, request.rSID)
					shList[request.rSID] = data
				}
				if val, ok := (*data)[request.rKey]; ok {
					result.sValue = val
				}
				request.reply <- result

			case smLoadSession:
				if _, ok := shList[request.rSID]; !ok {
					shList[request.rSID] = loadSession(aSessionDir, request.rSID)
				}
				request.reply <- &TSession{sID: request.rSID}

			case smSessionLen:
				result := &TSession{
					sID:    request.rSID,
					sValue: 0,
				}
				if data, ok := shList[request.rSID]; ok {
					result.sValue = len(*data)
				}
				request.reply <- result

			case smSetKey:
				if data, ok := shList[request.rSID]; ok {
					(*data)[request.rKey] = request.rValue
				} else {
					data = loadSession(aSessionDir, request.rSID)
					(*data)[request.rKey] = request.rValue
					shList[request.rSID] = data
				}
				request.reply <- &TSession{sID: request.rSID}

			case smStoreSession:
				if data, ok := shList[request.rSID]; ok {
					if 0 == len(*data) {
						// free unused memory
						delete(shList, request.rSID)
					} else {
						go goStore(aSessionDir, request.rSID, data)
					}
				}
				request.reply <- &TSession{sID: request.rSID}

			case smTerminate:
				return
			} // switch

		case <-timer.C:
			go goGC(aSessionDir)
			timer.Reset(gcInterval)
		} // select
	} // for
} // goMonitor()

// `goRemove()` removes the session file.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aSID` The session ID being destroyed.
func goRemove(aSessionDir, aSID string) {
	fName := filepath.Join(aSessionDir, aSID) + ".sid"
	// we try to remove the file w/o any checks
	_ = os.Remove(fName)
} // goRemove()

// `goStore()` saves `aData` of `aSID` on disk.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aSID` The session ID of the data to be stored.
//	`aData` The session data to store.
func goStore(aSessionDir, aSID string, aData *tSessionData) {
	now := time.Now()
	ss := tStoreStruct{
		"data":    *aData,
		"expires": now.Unix() + int64(soSessionTTL) + 1,
		"sid":     aSID,
	}
	gob.Register(*aData)
	gob.Register(now)
	gob.Register(ss)

	fName := filepath.Join(aSessionDir, aSID) + ".sid"
	file, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if nil != err {
		return
	}
	defer file.Close()

	_ = gob.NewEncoder(file).Encode(ss)
} // goStore()

// `loadSession()` reads the data for `aSID` from disk.
// If no (previous) session data is available, an empty session
// is returned.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aSID` The session ID whose data are to be read from disk.
func loadSession(aSessionDir, aSID string) *tSessionData {
	sData := make(tSessionData)
	fName := filepath.Join(aSessionDir, aSID) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDONLY, 0)
	if nil != err {
		return &sData
	}
	defer file.Close()

	var ss tStoreStruct
	now := time.Now()
	gob.Register(sData)
	gob.Register(now)
	gob.Register(ss)
	decoder := gob.NewDecoder(file)
	_ = decoder.Decode(&ss) // ignoring error: the following tests will catch it
	if e, ok := ss["expires"]; ok {
		if expireSecs, ok := e.(int64); ok &&
			time.Unix(expireSecs, 0).After(now) {
			if id, ok := ss["sid"]; ok {
				if sid, ok := id.(string); ok &&
					(sid == aSID) {
					if d, ok := ss["data"]; ok {
						if data, ok := d.(tSessionData); ok {
							return &data
						}
					}
				}
			}
		}
	}

	return &sData
} // loadSession()

/* _EoF_ */
