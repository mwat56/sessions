/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package sessions

type (

	// TMapSession is a `map` bases session store
	TMapSession struct {
		dmData tSessionData
		dmID   string
	}
)

// Data returns the complete set of session data.
func (dm *TMapSession) Data() *tSessionData {
	return &dm.dmData
} // Data()

// Delete removes the session data identified by `aKey`.
func (dm *TMapSession) Delete(aKey string) error {
	// If m is nil or there is no such element, delete is a no-op.
	delete(dm.dmData, aKey)

	return nil
} // Delete()

// Get returns the session data identified by `aKey`.
//
// If `aKey` doesn't exist the method returns `nil`.
func (dm *TMapSession) Get(aKey string) interface{} {
	if result, ok := dm.dmData[aKey]; ok {
		return result
	}

	return nil
} // Get()

// SessionID returns the session's ID.
func (dm *TMapSession) SessionID() string {
	return dm.dmID
} // SessionID()

// Set adds/updates the session data of `aKey` with `aValue`.
//
// This implementation always returns `nil`.
func (dm *TMapSession) Set(aKey string, aValue interface{}) error {
	dm.dmData[aKey] = aValue

	return nil
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newSession()` returns a new `TMapSession` instance.
func newSession(aSID string) *TMapSession {
	result := TMapSession{
		dmData: make(tSessionData),
		dmID:   aSID,
	}

	return &result
} // newSession()

/* _EoF_ */
