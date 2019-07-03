/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"path/filepath"
	"testing"
)

func Test_doRequest(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	go sessionMonitor(sdir, chSession)
	defer func() {
		chSession <- tShRequest{req: shTerminate}
	}()
	sid := "aTestSID2"

	type args struct {
		aSID     string
		aRequest tShLookupType
	}
	tests := []struct {
		name string
		args args
		want int //*TSession
	}{
		// TODO: Add test cases.
		{" 1", args{sid, shLoadSession}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doRequest(tt.args.aSID, tt.args.aRequest); got.Len() != tt.want {
				t.Errorf("doRequest() = %v,\nwant %v", got.Len(), tt.want)
			}
		})
	}
} // Test_doRequest()

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

func TestNewSession(t *testing.T) {
	sdir, _ := filepath.Abs("./sessions")
	go sessionMonitor(sdir, chSession)
	defer func() {
		chSession <- tShRequest{req: shTerminate}
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
