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

func Test_tHRefWriter_appendSID(t *testing.T) {
	sid := "aTestSIDvalue"
	h1 := tHRefWriter{sID: sid}
	d0 := []byte(`bla bla bla`)
	w0 := d0
	d1 := []byte(`Bla bla <a title="Link(1)" href="page1.html">Link(1)</a>`)
	w1 := []byte(`Bla bla <a title="Link(1)" href="page1.html?` + string(sidName) + `=` + sid + `">Link(1)</a>`)
	d2 := []byte(`Bla bla <a title="Link(2)" href="http://example.com/page2.html">Link(2)</a>`)
	w2 := d2
	d3 := []byte(`Bla bla <a title="Link(3)" href="page3.html?k=v">Link(3)</a>`)
	w3 := []byte(`Bla bla <a title="Link(3)" href="page3.html?k=v&` + string(sidName) + `=` + sid + `">Link(3)</a>`)

	d4 := []byte(`Bla bla <a title="Link(4)" href="page4.html?k=v#fragment">Link(4)</a>`)
	w4 := []byte(`Bla bla <a title="Link(4)" href="page4.html?k=v&` + string(sidName) + `=` + sid + `#fragment">Link(4)</a>`)

	d5 := []byte(string(d1) + string(d2) + string(d3))
	w5 := []byte(string(w1) + string(w2) + string(w3))
	type args struct {
		aData []byte
	}
	tests := []struct {
		name   string
		fields tHRefWriter
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
		{" 0", h1, args{d0}, w0},
		{" 1", h1, args{d1}, w1},
		{" 2", h1, args{d2}, w2},
		{" 3", h1, args{d3}, w3},
		{" 4", h1, args{d4}, w4},
		{" 5", h1, args{d5}, w5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := &tt.fields
			if got := hr.appendSID(tt.args.aData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tHRefWriter.appendSID() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_tHRefWriter_appendSID()
