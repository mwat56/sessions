/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"testing"
)

func TestExcludePaths(t *testing.T) {
	excludeList = nil // make sure to start with a fresh list
	type args struct {
		aPath []string
	}
	x1 := []string{"certs/", "css/", "/favicon", "fonts/"}
	l1 := len(x1)
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", args{x1}, l1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExcludePaths(tt.args.aPath...); got != tt.want {
				t.Errorf("AddExcludePaths() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestExcludePaths()

func Test_excludeURL(t *testing.T) {
	excludeList = nil // make sure to start with a fresh list
	ExcludePaths("css/", "/favicon", "/img/")
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
		{" 5", args{"/img_cover.png"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := excludeURL(tt.args.aURLpath); got != tt.want {
				t.Errorf("excludeURL() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_excludeURL()
