/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

type (

	// tSessionData stores the session data.
	tSessionData map[string]interface{}

	// TSession is a `map` based session store
	TSession struct {
		sData tSessionData
		sID   string
	}
)

// Delete removes the session data identified by `aKey`.
func (so *TSession) Delete(aKey string) error {
	// If m is nil or there is no such element, delete is a no-op.
	delete(so.sData, aKey)

	return nil
} // Delete()

// Get returns the session data identified by `aKey`.
//
// If `aKey` doesn't exist the method returns `nil`.
func (so *TSession) Get(aKey string) interface{} {
	if result, ok := so.sData[aKey]; ok {
		return result
	}

	return nil
} // Get()

// SessionID returns the session's ID.
func (so *TSession) SessionID() string {
	return so.sID
} // SessionID()

// Set adds/updates the session data of `aKey` with `aValue`.
//
// This implementation always returns `nil`.
func (so *TSession) Set(aKey string, aValue interface{}) error {
	so.sData[aKey] = aValue

	return nil
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newSession()` returns a new `TMapSession` instance.
func newSession(aSID string) *TSession {
	result := TSession{
		sData: make(tSessionData),
		sID:   aSID,
	}

	return &result
} // newSession()

/* _EoF_ */
