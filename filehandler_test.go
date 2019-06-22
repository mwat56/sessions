/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"reflect"
	"testing"
)

func TestTFileSessionHandler_setSessionDir(t *testing.T) {
	fh1 := TSessionHandler{}
	type args struct {
		aSavePath string
	}
	tests := []struct {
		name    string
		fields  TSessionHandler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, args{"./sessions"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &tt.fields
			if err := fs.initSessionDir(tt.args.aSavePath); (err != nil) != tt.wantErr {
				t.Errorf("TFileSessionHandler.initSessionDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTFileSessionHandler_setSessionDir()

func TestTFileSessionHandler_store(t *testing.T) {
	fh1, _ := newFilehandler("./sessions")
	s1 := newSession("aTestSID")

	type args struct {
		aSession *TSession
	}
	tests := []struct {
		name    string
		fields  *TSessionHandler // fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, args{s1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fields
			if err := fs.store(tt.args.aSession); (err != nil) != tt.wantErr {
				t.Errorf("TFileSessionHandler.store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTFileSessionHandler_store()

func TestTFileSessionHandler_load(t *testing.T) {
	fh1, _ := newFilehandler("./sessions")
	sid := "aTestSID"
	s1 := newSession(sid)
	w1 := &TSession{
		sData: tSessionData{},
		sID:   sid,
	}
	type args struct {
		aSID string
	}
	tests := []struct {
		name    string
		fields  *TSessionHandler // fields
		args    args
		want    *TSession
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, args{s1.SessionID()}, w1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fields
			got, err := fs.load(tt.args.aSID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TFileSessionHandler.load() error = %v,\nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TFileSessionHandler.load() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTFileSessionHandler_load()

func TestTFileSessionHandler_GC(t *testing.T) {
	fh1, _ := newFilehandler("./sessions")
	sid := "aTestSID"
	s1 := newSession(sid)
	fh1.store(s1)
	tests := []struct {
		name    string
		fields  *TSessionHandler // fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fields
			if err := fs.GC(); (err != nil) != tt.wantErr {
				t.Errorf("TFileSessionHandler.GC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTFileSessionHandler_GC()
