/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

import (
	"net/http"
)

type (
	// `tLogWriter` embeds a `ResponseWriter`
	tHRefWriter struct {
		http.ResponseWriter // used to construct the HTTP response
		session             *TSession
	}
)

// Write writes the data to the connection as part of an HTTP reply.
//
// Part of the `http.ResponseWriter` interface.
func (hr *tHRefWriter) Write(aData []byte) (int, error) {
	aData = hr.appendSID(aData)

	return hr.ResponseWriter.Write(aData)
} // Write()

// `appendSID()` appends the current session ID to all `a href` tags.
//
// `aData` The web/http response.
func (hr *tHRefWriter) appendSID(aData []byte) []byte {
	// cgi := fmt.Sprintf("%s=%s", sidName, hr.session.ID())
	/*
		There are three cases to consider:
			(a) links to external pages
			(b) links to internal pages w/o CGI arguments
			(c) links to internal pages with CGI arguments

	*/

	return aData
} // appendSID()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

/* _EoF_ */
