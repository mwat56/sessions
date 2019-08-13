/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

//lint:file-ignore ST1017 - I prefer Yoda conditions

import "strings"

type (
	tExcludeList []string
)

var (
	excludeList tExcludeList // The zero value of a slice is `nil`.
)

// ExcludePaths appends the `aPath` arguments to the list of
// ignored URL paths.
//
// The given `aPath` arguments are supposed to be the start (beginning)
// of the respective URL to exclude from session handling.
// If an `aPath` argument doesn't start with a slash (`/`) it's
// automatically prepended.
//
//	aPath List of URL paths to skip in session handling.
//	The return value is the current length of the exclude's list.
func ExcludePaths(aPath ...string) int {
	if nil == excludeList { // lazy initialisation
		excludeList = make(tExcludeList, 0, len(aPath)+16)
	}
	for _, path := range aPath {
		if '/' != path[0] {
			path = "/" + path
		}
		excludeList = append(excludeList, path)
	}

	return len(excludeList)
} // ExcludePaths()

// `excludeURL()` returns whether `aURLpath` is one to skip.
//
//	aURLpath The URL path to ckeck for.
func excludeURL(aURLpath string) bool {
	if '/' != aURLpath[0] { // relative paths may omit leading slash
		aURLpath = "/" + aURLpath
	}
	if nil != excludeList {
		for _, skipPath := range excludeList {
			if strings.HasPrefix(aURLpath, skipPath) {
				return true
			}
		}
	}

	return false
} // excludeURL()

/* _EoF_ */
