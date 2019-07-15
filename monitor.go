/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

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

	// `tShLookupType` is the kind of request to `sessionMonitor()`
	tShLookupType int

	// `tShRequest` is the request structure channeled to `sessionMonitor()`
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
	// The possible request types send to `sessionMonitor()`
	shNone = tShLookupType(1 << iota)
	shChangeSession
	shDeleteKey
	shDestroySession
	shGetKey
	shLoadSession
	shSessionLen
	shSetKey
	shStoreSession
	shTerminate // for testing only: terminate `sessionMonitor()`
)

// `goGC()` cleans up old sessions.
//
// Sessions that have not been updated for at least
// `SessionTTL()` seconds will be removed.
//
//	`aSessionDir` The directory where the session files are stored.
func goGC(aSessionDir string) {
	secs := time.Now().Unix() - int64(sessionTTL)
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
			sid := fName[:len(fName)-4]
			go goRemove(aSessionDir, sid)
		}
	}
} // goGC()

// `goMonitor()` handles the access to the internal list of session data.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aRequest` is the channel to receive request through.
func goMonitor(aSessionDir string, aRequest <-chan tShRequest) {
	shList := make(tShList, 32) // list of active sessions
	timer := time.NewTimer(time.Duration(sessionTTL<<4)*time.Second + 1)
	defer timer.Stop()

	for { // wait for requests
		select {
		case request, more := <-aRequest:
			if !more { // channel closed
				return
			}
			switch request.rType {

			case shChangeSession:
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

			case shDeleteKey:
				if data, ok := shList[request.rSID]; ok {
					delete(*data, request.rKey)
				}
				request.reply <- &TSession{sID: request.rSID}

			case shDestroySession:
				delete(shList, request.rSID)
				go goRemove(aSessionDir, request.rSID)
				request.reply <- &TSession{}

			case shGetKey:
				result := &TSession{
					sID: request.rSID,
				}
				if data, ok := shList[request.rSID]; ok {
					if val, ok := (*data)[request.rKey]; ok {
						result.sValue = val
					}
				}
				request.reply <- result

			case shLoadSession:
				data, ok := shList[request.rSID]
				if !ok {
					data = loadSession(aSessionDir, request.rSID)
				}
				shList[request.rSID] = data
				request.reply <- &TSession{sID: request.rSID}

			case shSessionLen:
				result := &TSession{
					sID: request.rSID,
				}
				if data, ok := shList[request.rSID]; ok {
					result.sValue = len(*data)
				}
				request.reply <- result

			case shSetKey:
				if data, ok := shList[request.rSID]; ok {
					(*data)[request.rKey] = request.rValue
				}
				request.reply <- &TSession{sID: request.rSID}

			case shStoreSession:
				if data, ok := shList[request.rSID]; ok {
					go goStore(aSessionDir, request.rSID, *data)
				}
				request.reply <- &TSession{sID: request.rSID}

			case shTerminate:
				return
			} // switch

		case <-timer.C:
			go goGC(aSessionDir)
			timer.Reset(time.Duration(sessionTTL<<4)*time.Second + 1)
		} // select
	} // for
} // goMonitor()

// `goRemove()` removes the session file.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aSID` The session ID being destroyed.
func goRemove(aSessionDir, aSID string) {
	fName := filepath.Join(aSessionDir, aSID) + ".sid"
	if _, err := os.Stat(fName); nil != err {
		return
	}

	os.Remove(fName)
} // goRemove()

// `goStore()` saves `aData` of `aSID` on disk.
//
//	`aSessionDir` The directory where the session files are stored.
//	`aSID` the session ID of the datra to be stored.
//	`aData` The session data to store.
func goStore(aSessionDir string, aSID string, aData tSessionData) {
	now := time.Now()
	ss := tStoreStruct{
		"data":    aData,
		"expires": now.Unix() + int64(sessionTTL) + 1,
		"sid":     aSID,
	}
	gob.Register(aData)
	gob.Register(now)
	gob.Register(ss)

	fName := filepath.Join(aSessionDir, aSID) + ".sid"
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if nil != err {
		return
	}
	defer file.Close()

	gob.NewEncoder(file).Encode(ss)
} // goStore()

// `loadSession()` reads the data for `aSID` from disk.
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
	err = decoder.Decode(&ss)
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
