/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func initTestSession() string {
	sdir, _ := checkSessionDir("./sessions")
	go goMonitor(sdir, chSession)
	sData := make(tSessionData)
	sData["Datum"] = time.Now()
	sData["Real"] = 12345.6789
	sData["Wahr"] = true
	sData["Zahl"] = 123456789
	sData["Zeichenkette"] = "eine Zeichenkette"
	sid := newSID() // "aTestSID"
	goStore(sdir, sid, sData)
	so := &TSession{
		sID: sid,
	}
	_ = so.request(smLoadSession, "", nil)

	return sid
} // initTestSession()

func initRequest() (string, *http.Request) {
	sid := initTestSession()
	result := httptest.NewRequest("GET", "/", nil)

	// prepare a reference for `GetSession()`
	ctx := context.WithValue(result.Context(), sidName, sid)
	result = result.WithContext(ctx)

	return sid, result
} // initRequest()

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

func TestGetSession(t *testing.T) {
	sid, req := initRequest()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	w1 := &TSession{sID: sid}
	type args struct {
		aRequest *http.Request
	}
	tests := []struct {
		name string
		args args
		want *TSession
	}{
		// TODO: Add test cases.
		{" 1", args{req}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSession(tt.args.aRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestGetSession()

func TestTSession_Get(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
	now := time.Now()
	s1.Set("Datum", now)
	type args struct {
		aKey string
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
		{" 1", s1, args{"Datum"}, now},
		{" 2", s1, args{"Real"}, 12345.6789},
		{" 3", s1, args{"Wahr"}, true},
		{" 4", s1, args{"Zahl"}, 123456789},
		{" 5", s1, args{"Zeichenkette"}, "eine Zeichenkette"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			if got := so.Get(tt.args.aKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TSession.Get() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTSession_Get()

func TestTSession_GetBool(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
	type args struct {
		aKey string
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   bool
		want1  bool
	}{
		// TODO: Add test cases.
		{" 1", s1, args{"gibbet nich"}, false, false},
		{" 2", s1, args{"Real"}, false, false},
		{" 3", s1, args{"Wahr"}, true, true},
		{" 4", s1, args{"Zahl"}, false, false},
		{" 5", s1, args{"Zeichenkette"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &TSession{
				sID:    tt.fields.sID,
				sValue: tt.fields.sValue,
			}
			got, got1 := so.GetBool(tt.args.aKey)
			if got != tt.want {
				t.Errorf("TSession.GetBool() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSession.GetBool() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
} // TestTSession_GetBool()

func TestTSession_GetFloat(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
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
		{" 1", s1, args{"gibbet nich"}, 0, false},
		{" 2", s1, args{"Real"}, 12345.6789, true},
		{" 3", s1, args{"Wahr"}, 0, false},
		{" 4", s1, args{"Zahl"}, 0, false},
		{" 5", s1, args{"Zeichenkette"}, 0, false},
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

func TestTSession_GetInt(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
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
		{" 1", s1, args{"gibbet nich"}, 0, false},
		{" 2", s1, args{"Real"}, 0, false},
		{" 3", s1, args{"Wahr"}, 0, false},
		{" 4", s1, args{"Zahl"}, 123456789, true},
		{" 5", s1, args{"Zeichenkette"}, 0, false},
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

func TestTSession_GetString(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
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
		{" 1", s1, args{"gibbet nich"}, "", false},
		{" 2", s1, args{"Real"}, "", false},
		{" 3", s1, args{"Wahr"}, "", false},
		{" 4", s1, args{"Zahl"}, "", false},
		{" 5", s1, args{"Zeichenkette"}, "eine Zeichenkette", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			got, got1 := so.GetString(tt.args.aKey)
			if got != tt.want {
				t.Errorf("TSession.GetString() got = %v,\nwant %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TSession.GetString() got1 = %v,\nwant %v", got1, tt.want1)
			}
		})
	}
} // TestTSession_GetString()

func TestTSession_GetTime(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
	var zero time.Time
	now := time.Now()
	s1.Set("Datum", now)
	type args struct {
		aKey string
	}
	tests := []struct {
		name      string
		fields    TSession
		args      args
		wantRTime time.Time
		wantROK   bool
	}{
		// TODO: Add test cases.
		{" 1", s1, args{"gibbet nich"}, zero, false},
		{" 2", s1, args{"Real"}, zero, false},
		{" 3", s1, args{"Wahr"}, zero, false},
		{" 4", s1, args{"Zahl"}, zero, false},
		{" 5", s1, args{"Zeichenkette"}, zero, false},
		{" 6", s1, args{"Datum"}, now, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &TSession{
				sID:    tt.fields.sID,
				sValue: tt.fields.sValue,
			}
			gotRTime, gotROK := so.GetTime(tt.args.aKey)
			if !reflect.DeepEqual(gotRTime, tt.wantRTime) {
				t.Errorf("TSession.GetTime() gotRTime = %v,\nwant %v", gotRTime, tt.wantRTime)
			}
			if gotROK != tt.wantROK {
				t.Errorf("TSession.GetTime() gotROK = %v,\nwant %v", gotROK, tt.wantROK)
			}
		})
	}
} // TestTSession_GetTime()

func TestTSession_Len(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
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

func TestTSession_request(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
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
		{" 1", s1, args{smLoadSession, "", nil}, w1, 5},
		{" 2", s2, args{smLoadSession, "", nil}, w2, 0},
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

func TestTSession_Set(t *testing.T) {
	sid := initTestSession()
	defer func() {
		chSession <- tShRequest{rType: smTerminate}
	}()
	s1 := TSession{sID: sid}
	w1 := &s1
	type args struct {
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name   string
		fields TSession
		args   args
		want   *TSession
	}{
		// TODO: Add test cases.
		{" 1", s1, args{"testkey1", "value1"}, w1},
		{" 2", s1, args{"testkey2", true}, w1},
		{" 3", s1, args{"Datum", time.Now()}, w1},
		{" 4", s1, args{"Real", 12345.6789}, w1},
		{" 5", s1, args{"Wahr", false}, w1},
		{" 6", s1, args{"Zahl", 987654321}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := &tt.fields
			got := so.Set(tt.args.aKey, tt.args.aValue)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TSession.Set() = %v,\nwant %v", got, tt.want)
			}
			if got.Get(tt.args.aKey) != tt.args.aValue {
				t.Errorf("TSession.Set() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTSession_Set()
