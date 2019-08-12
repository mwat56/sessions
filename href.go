/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
)

type (
	// `tHRefWriter` embeds a `ResponseWriter`
	tHRefWriter struct {
		http.ResponseWriter // used to construct the HTTP response
		sID                 string
	}

	// `tBoolLookup` is a simple binary lookup table
	tBoolLookup map[bool]string
)

var (
	// RegEx to match complete link tags
	hrefRE = regexp.MustCompile(`(?si)(<a[^>]*href=")([^"#]+)([^"]*"[^>]*>)`)

	// lookup table for appending CGI argument
	lookupCGIchar = tBoolLookup{true: "&", false: "?"}

	// check whether an URL starts with a scheme
	schemeRE = regexp.MustCompile(`^\w+:`)
)

// `appendSID()` appends the current session ID to all local `a href` tags.
//
// `aData` The web/http response.
func (hr *tHRefWriter) appendSID(aData []byte) []byte {
	so := &TSession{sID: hr.sID}
	if so.EmptySession() {
		return aData
	}
	linkMatches := hrefRE.FindAllSubmatch(aData, -1)
	if nil == linkMatches {
		return aData
	}
	cgi := fmt.Sprintf("%s=%s", sidName, hr.sID)
	/*
		There are three cases to consider:
			(a) links to external pages (ignored)
			(b) links to internal pages w/o CGI arguments
			(c) links to internal pages with CGI arguments
	*/
	for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
		if 0 == len(linkMatches[cnt][2]) { // the URL to check
			continue
		}
		if schemeRE.Match(linkMatches[cnt][2]) {
			continue // skip links to external sites
		}
		if excludeURL(string(linkMatches[cnt][2])) {
			continue // skip excluded URLs
		}
		repl := fmt.Sprintf("%s%s%s%s%s",
			linkMatches[cnt][1],
			linkMatches[cnt][2],
			lookupCGIchar[0 < bytes.IndexRune(linkMatches[cnt][2], '?')],
			cgi,
			linkMatches[cnt][3])
		aData = bytes.ReplaceAll(aData, linkMatches[cnt][0], []byte(repl))
	}

	return aData
} // appendSID()

// Write writes the data to the connection as part of an HTTP reply.
//
// Part of the `http.ResponseWriter` interface.
func (hr *tHRefWriter) Write(aData []byte) (int, error) {
	return hr.ResponseWriter.Write(hr.appendSID(aData))
} // Write()

/* _EoF_ */
