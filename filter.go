/*
   Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides functions to ignore certain URL paths
 * from session handling.
*/

import (
	"strings"
)

type (
	tExcludeList []string
)

var (
	soFilterExcludeList tExcludeList // The zero value of a slice is `nil`.
)

// ExcludePaths appends the `aPath` argument(s) to the list of
// URL paths to ignore.
//
// The given `aPath` arguments are supposed to be the start (beginning)
// of the respective URL to exclude from session handling.
// If an `aPath` argument doesn't start with a slash (`/`) it's
// automatically prepended.
//
//	aPath List of URL paths to skip in session handling.
//	The return value is the current length of the exclude list.
func ExcludePaths(aPath ...string) int {
	if nil == soFilterExcludeList { // lazy initialisation
		soFilterExcludeList = make(tExcludeList, 0, len(aPath)+8)
	}
	for _, path := range aPath {
		if '/' != path[0] {
			path = "/" + path
		}
		soFilterExcludeList = append(soFilterExcludeList, path)
	}

	return len(soFilterExcludeList)
} // ExcludePaths()

// `excludeURL()` returns whether `aURLpath` is one to skip.
//
//	aURLpath The URL path to ckeck for.
func excludeURL(aURLpath string) bool {
	if nil == soFilterExcludeList {
		return false
	}
	if '/' != aURLpath[0] { // relative paths may omit leading slash
		aURLpath = "/" + aURLpath
	}
	for _, skipPath := range soFilterExcludeList {
		if strings.HasPrefix(aURLpath, skipPath) {
			return true
		}
	}

	return false
} // excludeURL()

/* _EoF_ */
