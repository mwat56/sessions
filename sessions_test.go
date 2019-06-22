/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import "testing"

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

func Test_newSession(t *testing.T) {
	sid := "aTestSID"
	w1 := &TSession{
		sData: tSessionData{},
		sID:   sid,
	}
	type args struct {
		aSID string
	}
	tests := []struct {
		name string
		args args
		want *TSession
	}{
		// TODO: Add test cases.
		{" 1", args{sid}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSession(tt.args.aSID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSession() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // Test_newSession()
