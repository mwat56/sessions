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
