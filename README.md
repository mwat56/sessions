# Sessions

[![Golang](https://img.shields.io/badge/Language-Go-green.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/mwat56/sessions?status.svg)](https://godoc.org/github.com/mwat56/sessions/)
[![Go Report](https://goreportcard.com/badge/github.com/mwat56/sessions)](https://goreportcard.com/report/github.com/mwat56/sessions)
[![Issues](https://img.shields.io/github/issues/mwat56/sessions.svg)](https://github.com/mwat56/sessions/issues?q=is%3Aopen+is%3Aissue)
[![Size](https://img.shields.io/github/repo-size/mwat56/sessions.svg)](https://github.com/mwat56/sessions/)
[![Tag](https://img.shields.io/github/tag/mwat56/sessions.svg)](https://github.com/mwat56/sessions/tags)
[![License](https://img.shields.io/github/license/mwat56/sessions.svg)](https://github.com/mwat56/sessions/blob/main/LICENSE)

- [Sessions](#sessions)
	- [Purpose / About](#purpose--about)
	- [Installation](#installation)
	- [Usage](#usage)
	- [Hints](#hints)
		- [Session files](#session-files)
		- [GETter](#getter)
	- [Internals](#internals)
		- [Session name](#session-name)
		- [GC](#gc)
	- [Licence](#licence)

----

## Purpose / About

I wanted a session data solution that's user-friendly – including privacy-friendly.
When doing some research about saving/retrieving sessions data with `Go` (aka `Golang`) you'll find some slightly different solutions which have, however, one detail in common: they all depend on socalled internet `cookies`.
Which is bad.

> _Cookies are (usually) bad_.

In practice `cookies` are basically an invasion into the user's property (`cookies` claim disk space and require additional electricity for processing) and they are, by definition, kind of a surveillance and tracking tool.
Nobody who has their user's best interest in mind would consider using `cookies` under normal circumstance.
Their only advantage is that they are easy to implement – which was kind of the point initially: ease of implementation.
On the other hand the remote users were considered just a passive and obedient consumer – an assumption, however, you can't really make in general:
Since `cookies` are stored on the remote user's computer you don't really have control over that piece of data but instead the remote user's computer (which means the remote user) ultimately controls this data and thus can easily manipulate it.
In other words: _`cookies` are inherently insecure_.
It's clear that nowadays – with data security and the user's privacy in mind – using `cookies` is basically an outdated technique.
It's also clear that harvesting the user's facilities (including disk space and electricity) should be avoided wherever possible.

Additionally using `Cookies` requires you to use `JavaScript` as well (which is another barrier best to be avoided).
In the European Union – with all its currently 27 member countries – `Cookies` are allowed only if the user explicitly agrees; in other words: they may _not be set automatically_ and silently.
And to get the user's consent you'd need JavaScript code which then – _after_ reading the user's reaction – either sets a `Cookie` _or not_.
So if you care for a barrier-free web-presentation and want to respect privacy and data-protection laws you can't use `Cookies`.

Another flaw you'll find in the literature about user sessions is the fact that it's often primarily considered in connection with users who are in one way or another _logged in_ with the web-server.
But that is only _one_ possible reason for considering some kind of session (data) store.
Others may include, for example, individual navigation history (e.g. breadcrumbs), configuration data, personal preferences, etc.
In any case, session data should be able to transcend certain states of a user's connection.
And – that's a very important point – sessions should never _require_ a user to `log in`.
The conception of the internet (more than a decade before that term became a synonym for commercial enterprise in the early 90s of the last century) was based on the idea of providing free access to as many data (i.e. knowledge) for as many people as possible.
Apart from other things that means: session data should never be used as a kind of lock-in or surveillance measure.

Some people argue that cookies are technically necessary because without them you couldn't implement something like shopping carts etc.
However, those people clearly confuse necessity with ease or comfort.

The ease of use for the application developer was mentioned before (when discussing the inherent badness of `cookies`).
That critique does not, however, mean that it should be difficult to implement session handling as such.
In fact, a huge part of the time I spent developing this package was spent figuring out an easy way (for the developer) to handle session data.
Now, after it's done, everything seems just 'logical', as if it's obvious how to do things.
That's an impression probably every programmer knows – at least if they do reflect what they are doing.

An other point to consider was to find a solution that's not dependent on third party facilities.
That most prominently excludes external database systems (like e.g. MariaDB and others alike).
The only database system to use is the OS's filesystem, the only system, that is, which doesn't add another possible layer of failure.
And it is by definition faster than any other system because every other database system runs _on top_ of the filesystem and therefore adds some amount of overhead (either in system calls or memory use or both).

While it can be challenging and interesting to figure out some smart database structure and sophisticated queries to access and retrieve data – such an endeavour is often simply over-the-top.
That's true especially for a job like session handling that has just two tasks to accomplish: point and grab (i.e. loading/reading data) and throw and forget (i.e. writing/storing data).

So, in short: I wanted a system that's as unintrusive as possible (i.e. respecting the user's privacy) and doesn't depend on anything but what is there in any case (i.e. the filesystem).
And, of course, it should be easy to use by the developer.

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/sessions

## Usage

To include the session handling provided by this package you just call the `Wrap()` function as shown here:

	func main() {
		// …
		// system setup etc.

		sessionDir := "./sessions"
		// this value should probably come from some config file
		// or commandline option.

		pageHandler := http.NewServeMux() // or your own page handling provider
		pageHandler.HandleFunc("/", myHandler) // dito

		server := http.Server{
			Addr:    "127.0.0.1:8080",
			Handler: sessions.Wrap(pageHandler, sessionDir),
			//       ^^^^^^^^^^^^^
		}

		if err := server.ListenAndServe(); nil != err {
			log.Fatalf("%s: %v", os.Args[0], err)
		}
	} // main()

Then from inside your pagehandler `myHandler`:

	// …
	// `aRequest` is the `*http.Request` argument passed to your handler func
	mySession := sessions.GetSession(aRequest)
	// …
	myVal := mySession.Get("myKey")
	myVal2 := mySession.Get("myKey2")
	// do something with `myVal`/`myVal2`
	// …
	otherVal := "important session value"
	mySession.Set("otherKey", otherVal)
	someNumber:= 123.456
	mySession.Set("numberKey", someNumber)
	// …

Please note that both, the return value of `mySession.Get()` and the argument value of `mySession.Set()`, are defined/typed as `interface{}`.
That way you can store any value as session data.

## Hints

### Session files

Most web-pages contain more than just HTML markup but other elements as well like _stylesheets_, or _JavaScript_, or _images_.
When any of these elements are requested by the remote users (i.e. their browser) a unique session will be created – which is usually not what you want.
If you have, say, ten images and a stylesheet in your page this library will create 12 unique session IDs, one for the page itself (i.e. the HTML), ten for the images and another one for the stylesheet while you actually need only one (for the HTML page as such).
Since your page handler has to deal with serving all of the page elements you could get rid of the superfluous sessions by destroying them, for example:

	// …
	// in the handler's branch for images
	mySession := sessions.GetSession(aRequest)
	// …
	mySession.Destroy()
	// …

This way there will be no session file created for the unwanted page element.
The same is true if you simply don't use the session instance to store any data (calling `mySession.Set(…)`).

Or – you could just ignore this inconvenience and let the library's internal Garbage Collector take care of the unneeded sessions.
Empty sessions (i.e. sessions with no data added to it) will not be saved to disk.

To further mitigate the problem with creating sessions for e.g. media files there is a function to declare certain URLs to skip in session handling:

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
	func ExcludePaths(aPath ...string) int { … }

The local URLs paths (i.q. without a protocol and host part) passed to that function will be ignored by the session handling. e.g.

	sessions.ExcludePaths("css", "/favicon", "/img/")
	// …

All requested URL paths starting with one of the added exclude paths will be ignored; that means neither will be there a session created nor will the internal session ID be changed.

### GETter

The session object returned by `GetSession()` allows you to store and retrieve any data type.
This is possible by internally using the empty `interface{}` which in consequence means that you seem to loose all type safety.
To help you with retrieving some common data types the session objects provides a few different GETter methods:

* `Get(aKey string) interface{}` returns a seemingly untyped result;
* `GetBool(aKey string) (bool, bool)` returns a `Boolean` value;
* `GetFloat(aKey string) (float64, bool)` returns a `Float` result;
* `GetInt(aKey string) (int64, bool)` returns an `Int` result;
* `GetString(aKey string) (string, bool)` returns a `String` result;
* `GetTime(aKey string) (time.Time, bool)` returns a `Time` result.

The respective second (`bool`) return value signals whether the data associated with the respective `aKey` is indeed of the requested type.
If this is not the case that second return value will be `false` and the first return value will be the zero value of the respective type.

## Internals

The package loads the sessions data (if any) whenever a page is requested and it stores the session data when the page handling is finished (i.e. after the page request was served).
This is done automatically and you don't have to worry about loading/storing (read/write) the session data manually.

### Session name

The session ID (`SID`) is handled automatically as well.
Each ID is valid only for _a single request_ by the remote user and changes for each request.
The `SID` name can be changed by calling

	sessions.SetSIDname(aSID string)

if the default (i.e. `SID`) doesn't satisfy your requirements.
To get the current setting you can call

	sidname := sessions.SIDname()

The `SID` and the one-time-value are appended automatically as an [CGI argument](https://en.wikipedia.org/wiki/Common_Gateway_Interface) to all local `a href="…"` links of the web page sent to the remote user, whereas _local_ means all links without a request scheme like e.g. `https:`.
In other words: the current session ID is only available in the respective page's source code while the ID showing up in the browser's URL-line was valid only when the page was requested.

### GC

The package provides an internal garbage collector (GC) which deletes expired sessions.

Sessions are considered `expired` when they were not touched/updated within the last `SessionTTL()` seconds (_10 minutes_ by default).
This default time can be changed by calling

	sessions.SetSessionTTL(aTTL int)

with `aTTL` seconds as the new time-to-life. –
You can get the current TTL (in seconds) by calling

	ttl := sessions.SessionTTL()

To be on the safe side the GC runs in background with an interval of twice the TTL.

## Licence

        Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany
                        All rights reserved
                    EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.

----
