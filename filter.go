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

// AddExcludePath appends `aPath` to the list of ignored URL paths.
//
// The given `aPath` is supposed to be the start (beginning) of the URL to
// exclude from session handling.
// If `aPath` doesn't start with a slash (`/`) it's automatically prepended.
//
//	aPath An URL path to skip in session handling.
func AddExcludePath(aPath string) {
	if 0 == len(aPath) {
		return
	}
	if '/' != aPath[0] {
		aPath = "/" + aPath
	}
	if nil == excludeList { // lazy initialisation
		excludeList = make(tExcludeList, 1, 16)
		excludeList[0] = aPath
	} else {
		excludeList = append(excludeList, aPath)
	}
} // AddExcludePath()

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
