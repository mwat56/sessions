/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"testing"
)

func TestAddExcludePath(t *testing.T) {
	type args struct {
		aPath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{"certs"}},
		{" 2", args{"/css"}},
		{" 3", args{"fonts"}},
		{" 4", args{"/img"}},
		{" 5", args{"thumb"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddExcludePath(tt.args.aPath)
		})
	}
} // TestAddExcludePath()

func Test_excludeURL(t *testing.T) {
	AddExcludePath("css")
	AddExcludePath("favicon")
	AddExcludePath("img")
	type args struct {
		aURLpath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{" 1", args{"page.html"}, false},
		{" 2", args{"/css/styles.css"}, true},
		{" 3", args{"/img/cover.png"}, true},
		{" 4", args{"favicon.ico"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := excludeURL(tt.args.aURLpath); got != tt.want {
				t.Errorf("excludeURL() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_excludeURL()
