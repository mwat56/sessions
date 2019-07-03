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
	// List of known sessions
	tShList map[string]*tSessionData

	// the kind of request to `sessionMonitor()`
	tShLookupType int

	// the request structure transported to `sessionMonitor()`
	tShRequest struct {
		req   tShLookupType
		sid   string
		reply chan *TSession
	}
)

const (
	// The possible request types send to `sessionMonitor()`
	shNone = tShLookupType(1 << iota)
	shChangeSession
	shCloseSession
	shDestroySession
	shLoadSession
	shNewSession
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
//	`aSID` The session ID whose data are to be stored.
func goStore(aSessionDir, aSID string, aData *tSessionData) {
	now := time.Now()
	expireSec := now.Unix() + int64(sessionTTL) + 1
	ss := tStoreStruct{
		"data":    aData,
		"expires": expireSec,
		"sid":     aSID,
	}
	gob.Register(aData)
	gob.Register(now)
	gob.Register(ss)
	for _, val := range *aData {
		gob.Register(val)
	}

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
	gob.Register(&sData)
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
						if data, ok := d.(*tSessionData); ok {
							return data
						}
					}
				}
			}
		}
	}

	return &sData
} // loadSession()

// `sessionMonitor()` handles the access to the internal list of session data.
//
//	`aSessionDir` The directory where the session files are stored.
func sessionMonitor(aSessionDir string, aRequest <-chan tShRequest) {
	shList := make(tShList, 32) // list of known/active sessions
	timer := time.NewTimer(time.Duration(sessionTTL)*time.Second + 1)
	defer timer.Stop()

	for { // wait for requests
		select {
		case request := <-aRequest:
			switch request.req {
			case shChangeSession: // X
				newsid := newSID()
				result := TSession{sID: newsid}
				if data, ok := shList[request.sid]; ok {
					shList[newsid] = data
					delete(shList, request.sid)
					result.sData = data
				} else {
					list := make(tSessionData)
					shList[newsid] = &list
					result.sData = &list
				}
				go goRemove(aSessionDir, request.sid)
				request.reply <- &result

			case shCloseSession:
				go goGC(aSessionDir)
				request.reply <- &TSession{sID: request.sid}

			case shDestroySession: // X
				delete(shList, request.sid)
				go goRemove(aSessionDir, request.sid)
				request.reply <- &TSession{}

			case shLoadSession: // XX
				result := &TSession{
					sID: request.sid,
				}
				if data, ok := shList[request.sid]; ok {
					result.sData = data
				} else {
					result.sData = loadSession(aSessionDir, request.sid)
					shList[request.sid] = result.sData
				}
				request.reply <- result

			case shNewSession: // X
				data := make(tSessionData, 16)
				result := &TSession{
					sID:   request.sid,
					sData: &data,
				}
				shList[request.sid] = &data

				request.reply <- result

			case shStoreSession: // X
				if data, ok := shList[request.sid]; ok {
					go goStore(aSessionDir, request.sid, data)
				}
				request.reply <- &TSession{}

			case shTerminate:
				return
			} // switch

		case <-timer.C:
			go goGC(aSessionDir)
			timer.Reset(time.Duration(sessionTTL)*time.Second + 1)
		} // select
	} // for
} // sessionMonitor()

/* _EoF_ */
