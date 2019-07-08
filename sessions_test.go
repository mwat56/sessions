/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"reflect"
	"testing"
	"time"
)

func initTestSession() string {
	sdir, _ := checkSessionDir("./sessions")
	go goMonitor(sdir, chSession)
	sid := "aTestSID"
	session := doRequest(shLoadSession, sid, "", nil)
	session.Set("Zeichenkette", "eine Zeichenkette").
		Set("Zahl", 123456789).
		Set("Datum", time.Now()).
		Set("Real", 12345.6789)

	return sid
} // initTestSession()

func Test_doRequest(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	sid2 := "aTestSID2"
	w1 := &TSession{sID: sid}
	w2 := &TSession{sID: sid2}
	type args struct {
		aType  tShLookupType
		aSID   string
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name         string
		args         args
		wantRSession *TSession
		wantLen      int
	}{
		// TODO: Add test cases.
		{" 1", args{shLoadSession, sid, "", nil}, w1, 4},
		{" 2", args{shLoadSession, sid2, "", nil}, w2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRSession := doRequest(tt.args.aType, tt.args.aSID, tt.args.aKey, tt.args.aValue)
			if !reflect.DeepEqual(gotRSession, tt.wantRSession) {
				t.Errorf("doRequest() = %v, want %v", gotRSession, tt.wantRSession)
			}
			if tt.wantLen != gotRSession.Len() {
				t.Errorf("doRequest() = %v, want %v", gotRSession.Len(), tt.wantLen)
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
		fields TSession
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
		{" 2", 32},
		{" 3", 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSID(); len(got) != tt.want {
				t.Errorf("newID() = %v, want %v", len(got), tt.want)
			}
		})
	}
} // Test_newID()
