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
	so := &TSession{
		sID: "aTestSID",
	}
	so.request(shLoadSession, "", nil)
	so.Set("Zeichenkette", "eine Zeichenkette").
		Set("Zahl", 123456789).
		Set("Datum", time.Now()).
		Set("Real", 12345.6789).
		Set("Wahr", true)
	so.request(shStoreSession, "", nil)

	return so.sID
} // initTestSession()

func TestTSession_request(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	sid2 := "aTestSID2"
	s1 := TSession{sID: sid}
	w1 := &TSession{sID: sid}
	s2 := TSession{sID: sid2}
	w2 := &TSession{sID: sid2}
	type args struct {
		aType  tShLookupType
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name         string
		fields       TSession
		args         args
		wantRSession *TSession
		wantLen      int
	}{
		// TODO: Add test cases.
		{" 1", s1, args{shLoadSession, "", nil}, w1, 5},
		{" 2", s2, args{shLoadSession, "", nil}, w2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			gotRSession := so.request(tt.args.aType, tt.args.aKey, tt.args.aValue)
			if !reflect.DeepEqual(gotRSession, tt.wantRSession) {
				t.Errorf("TSession.request() = %v,\nwant %v", gotRSession, tt.wantRSession)
			}
			if tt.wantLen != gotRSession.Len() {
				t.Errorf("doRequest() = %v, want %v", gotRSession.Len(), tt.wantLen)
			}
		})
	}
} // TestTSession_request()

func TestTSession_Len(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	s1 := TSession{sID: sid}
	w1 := 5
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

func TestTSession_GetString(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	s1 := TSession{sID: sid}
	k1 := "gibbet nich"
	ws1 := ""
	wb1 := false
	k2 := "Zeichenkette"
	ws2 := "eine Zeichenkette"
	wb2 := true
	k3 := "Zahl"
	ws3 := ""
	wb3 := false
	type args struct {
		aKey string
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   string
		want1  bool
	}{
		// TODO: Add test cases.
		{" 1", s1, args{k1}, ws1, wb1},
		{" 2", s1, args{k2}, ws2, wb2},
		{" 3", s1, args{k3}, ws3, wb3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			got, got1 := so.GetString(tt.args.aKey)
			if got != tt.want {
				t.Errorf("TSession.GetString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSession.GetString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
} // TestTSession_GetString()

func TestTSession_GetInt(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	s1 := TSession{sID: sid}
	k1 := "gibbet nich"
	k2 := "Zahl"
	k3 := "Zeichenkette"
	type args struct {
		aKey string
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   int64
		want1  bool
	}{
		// TODO: Add test cases.
		{" 1", s1, args{k1}, 0, false},
		{" 2", s1, args{k2}, 123456789, true},
		{" 3", s1, args{k3}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			got, got1 := so.GetInt(tt.args.aKey)
			if got != tt.want {
				t.Errorf("TSession.GetInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSession.GetInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
} // TestTSession_GetInt()

func TestTSession_GetFloat(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: shTerminate}
	}()
	s1 := TSession{sID: sid}
	k1 := "gibbet nich"
	ws1 := float64(0)
	wb1 := false
	k2 := "Real"
	ws2 := 12345.6789
	wb2 := true
	k3 := "Zeichenkette"
	ws3 := float64(0)
	wb3 := false
	k4 := "Zahl"
	ws4 := float64(0)
	wb4 := false
	type fields struct {
		sID    string
		sValue interface{}
	}
	type args struct {
		aKey string
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   float64
		want1  bool
	}{
		// TODO: Add test cases.
		{" 1", s1, args{k1}, ws1, wb1},
		{" 2", s1, args{k2}, ws2, wb2},
		{" 3", s1, args{k3}, ws3, wb3},
		{" 4", s1, args{k4}, ws4, wb4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &TSession{
				sID:    tt.fields.sID,
				sValue: tt.fields.sValue,
			}
			got, got1 := so.GetFloat(tt.args.aKey)
			if got != tt.want {
				t.Errorf("TSession.GetFloat() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSession.GetFloat() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
} // TestTSession_GetFloat()
