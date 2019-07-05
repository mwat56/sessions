/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"testing"
	"time"
)

func initTestSession() string {
	sdir, _ := checkSessionDir("./sessions")
	go sessionMonitor(sdir, chSession)
	sid := "aTestSID"
	session := doRequest(shLoadSession, sid, "", nil)
	session.Set("Zeichenkette", "eine Zeichenkette").
		Set("Zahl", 123456789).
		Set("Datum", time.Now()).
		Set("Real", 12345.6789)
	// doRequest(shStoreSession, sid, "", nil)

	return sid
} // initTestSession()

func Test_doRequest(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	sid2 := "aTestSID2"

	type args struct {
		aRequest tShLookupType
		aSID     string
		aKey     string
		aValue   interface{}
	}
	tests := []struct {
		name string
		args args
		want int //*TSession
	}{
		// TODO: Add test cases.
		{" 1", args{shLoadSession, sid, "", nil}, 4},
		{" 2", args{shLoadSession, sid2, "", nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doRequest(tt.args.aRequest, tt.args.aSID, tt.args.aKey, tt.args.aValue); got.Len() != tt.want {
				t.Errorf("doRequest() = %v,\nwant %v", got.Len(), tt.want)
			}
		})
	}
} // Test_doRequest()

func TestTSession_Len(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	s1 := TSession{sID: sid}
	w1 := 4
	s2 := TSession{sID: "aTestSID2"}
	w2 := 0
	tests := []struct {
		name   string
		fields TSession // fields
		want   int
	}{
		// TODO: Add test cases.
		{" 1", s1, w1},
		{" 2", s2, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &TSession{
				sID:    tt.fields.sID,
				sValue: tt.fields.sValue,
			}
			if got := so.Len(); got != tt.want {
				t.Errorf("TSession.Len() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTSession_Len()

func Test_newID(t *testing.T) {
	tests := []struct {
		name string
		want int //string
	}{
		// TODO: Add test cases.
		{" 1", 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSID(); len(got) != tt.want {
				t.Errorf("newID() = %v, want %v", len(got), tt.want)
			}
		})
	}
} // Test_newID()
/*
func TestNewSession(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	go sessionMonitor(sdir, chSession)
	defer func() {
		chSession <- tShRequest{rReq: shTerminate}
	}()

	tests := []struct {
		name string
		want int //*TSession
	}{
		// TODO: Add test cases.
		{" 1", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSession(); got.Len() != tt.want {
				t.Errorf("NewSession() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestNewSession()
*/
