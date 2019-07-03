/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"path/filepath"
	"testing"
	"time"
)

func Test_goStore(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	sid := "aTestSID"
	list := make(tSessionData)
	s1 := &TSession{
		sData: &list,
		sID:   sid,
	}
	s1.Set("Zeichenkette", "eine Zeichenkette").
		Set("Zahl", 123456789).
		Set("Datum", time.Now())
	type args struct {
		aSessionDir string
		aSession    *TSession
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{sdir, s1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goStore(tt.args.aSessionDir, tt.args.aSession)
		})
	}
} // Test_goStore()

func Test_loadSession(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	sid := "aTestSID"
	type args struct {
		aSessionDir string
		aSID        string
	}
	tests := []struct {
		name string
		args args
		want int //*tSessionData
	}{
		// TODO: Add test cases.
		{" 1", args{sdir, sid}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadSession(tt.args.aSessionDir, tt.args.aSID); len(*got) != tt.want {
				t.Errorf("loadSession() = %v, want %v", len(*got), tt.want)
			}
		})
	}
} // Test_loadSession()
