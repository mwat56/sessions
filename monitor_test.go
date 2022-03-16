/*
   Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany
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
	sid := newSID()
	list := make(tSessionData)
	list["Zeichenkette"] = "eine Zeichenkette"
	list["Zahl"] = 123456789
	list["Datum"] = time.Now()
	type args struct {
		aSessionDir string
		aSID        string
		aData       tSessionData
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{sdir, sid, list}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goStore(tt.args.aSessionDir, tt.args.aSID, &tt.args.aData)
		})
	}
} // Test_goStore()

func Test_loadSession(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	sid := initTestSession()
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
		{" 1", args{sdir, sid}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadSession(tt.args.aSessionDir, tt.args.aSID); len(*got) != tt.want {
				t.Errorf("loadSession() = %v, want %v", len(*got), tt.want)
			}
		})
	}
} // Test_loadSession()
