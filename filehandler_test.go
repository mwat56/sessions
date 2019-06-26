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

func Test_tSessionHandler_setSessionDir(t *testing.T) {
	fh1 := tSessionHandler{}
	type args struct {
		aSavePath string
	}
	tests := []struct {
		name    string
		fields  tSessionHandler
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
				t.Errorf("tSessionHandler.initSessionDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tSessionHandler_setSessionDir()

func Test_tSessionHandler_store(t *testing.T) {
	fh1, _ := newSessionHandler("./sessions")
	s1 := newSession("aTestSID")
	s1.Set("Zeichenkette", "eine Zeichenkette")
	s1.Set("Zahl", 123456789)
	s1.Set("Datum", time.Now())

	type args struct {
		aSession *TSession
	}
	tests := []struct {
		name    string
		fields  *tSessionHandler // fields
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
				t.Errorf("tSessionHandler.store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tSessionHandler_store()

func Test_tSessionHandler_load(t *testing.T) {
	fh1, _ := newSessionHandler("./sessions")
	sid := "aTestSID"
	// s1 := newSession(sid)
	// s1.Set("Zeichenkette", "eine Zeichenkette")
	// s1.Set("Zahl", 123456789)
	// s1.Set("Datum", time.Now())
	// w1 := &TSession{
	// 	sData: s1.sData, // &tSessionData{},
	// 	sID:   sid,
	// }
	type args struct {
		aSID string
	}
	tests := []struct {
		name    string
		fields  *tSessionHandler // fields
		args    args
		want    int //*TSession
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, args{sid}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fields
			got, err := fs.load(tt.args.aSID)
			if (err != nil) != tt.wantErr {
				t.Errorf("tSessionHandler.load() error = %v,\nwantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			if got.Len() != tt.want {
				t.Errorf("tSessionHandler.load() = %v,\nwant %v", got.Len(), tt.want)
			}
		})
	}
} // Test_tSessionHandler_load()

func Test_tSessionHandler_GC(t *testing.T) {
	fh1, _ := newSessionHandler("./sessions")
	sid := "gcTestSID"
	s1 := newSession(sid)
	fh1.store(s1)
	tests := []struct {
		name    string
		fields  *tSessionHandler // fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fh1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fields
			if err := fs.GC(); (err != nil) != tt.wantErr {
				t.Errorf("tSessionHandler.GC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // Test_tSessionHandler_GC()
