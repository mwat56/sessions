# Sessions

[![GoDoc](https://godoc.org/github.com/mwat56/sessions?status.svg)](https://godoc.org/github.com/mwat56/sessions)
[![License](https://img.shields.io/eclipse-marketplace/l/notepad4e.svg)](https://github.com/mwat56/sessions/blob/master/LICENSE)

- [Sessions](#Sessions)
	- [Purpose / About](#Purpose--About)
	- [Installation](#Installation)
	- [Usage](#Usage)
	- [Internals](#Internals)
	- [Licence](#Licence)

## Purpose / About

I wanted a session data solution that's user-friendly – including privacy-friendly.
When doing some research about saving/retrieving sessions data with `Go` (aka `Golang`) you'll find some slightly different solutions which have, however, one detail in common: they all depend on socalled internet `cookies`.
Which is bad.

_Cookies are bad_.

In practice `cookies` are basically an invasion into the user's property (`cookies` claim disk space and require additional electricity for processing) and they are, by definition, kind of a surveillance and tracking tool.
Nobody who has their user's best interesst in mind would consider using `cookies`.
Their only advantage is that they are easy to implement – which was kind of the point: ease of implementation.
On the other hand the remote users were considered just a passive and obedient consumer – an essumption you can't really make in general:
Since `cookies` are stored on the remote user's computer you don't really have control over that piece of data but instead the remote user's computer (which means the remote user) ultimately controls this data and thus can easily manipulate it.
In other words: _`cookies` are inherently insecure_.
It's clear that nowadays – with data security and the user's privacy in mind – using `cookies` is just an outdated technique.
It's also clear that harvesting the user's facilities (including disk space and electricity) should be avoided.

Another flaw you'll find in the literature about user sessions is the fact that it's often primarily considered in connection with users who are in one way or another _logged in_ with the web-server.
But that is only _one_ possible reason for considering some kind of session (data) store.
Others may include, for example, individual navigation history (e.g. breadcrumbs), configuration data, personal preferences, etc.
In any case, session data should be able to transcent certain states of a user's connection.
And – that's a very important point – sessions should never _require_ a user to `log in`.
The conception of the internet (more than a decade before that term became a synonym for commercial enterprise in the early 90s of the last centrury) was based on the idea of providing free access to as many data (i.e. knowledge) for as many people as possible.
Apart from other things that means: session data should never be used as a kind of lock-in or surveillance measure.

The ease of use for the application developer was mentioned before (when discussing the inherent badness of `cookies`).
That critique does not, however, mean that it should be difficult to implement session handling as such.
In fact, a huge part of the time I spent developing this package was spent figuring out an easy way to handle session data.
Now, after it's done, everything seems just 'logical', as if it's obvious how to do things.
That's an impression probably every programmer knows – at least if they do reflect what they are doing.

An other point to consider was to find a solution that's not dependent on third party facilities.
That most prominently excludes external database systems (like e.g. MariaDB and others alike).
The only database system to use is the OS's filesystem, the only system, that is, which doesn't add another possible point of failure.
And it is by definition faster than any other system because every other database system runs _on top_ of the filesystem and therefore adds some amount of overhead.

While it can be challenging and interesting to figure out some smart database structure and sophisticated queries to access and retrieve data such an endeavour is often simply over-the-top.
That's true especially for a job like session handling that has just two tasks to accomplish: point and grab (i.e. reading/loading data) and throw and forget (i.e. writing/storing data).

So, in short: I wanted a system that's as unintrusive as possible (i.e. respecting the user's privacy) and doesn't depend on anything but what is there in any case (i.e. the filesystem).
And, of course, it should be easy to use by the developer.

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/sessions

## Usage

To include the session handling provided by this package you just call the `Wrap()` function as shown here:

	func main() {
		// ...
		// system setup etc.

		sessionDir := "./sessions"
		// this value should probably come from some config file
		// or commandline option.

		pageHandler := http.NewServeMux() // or you own page handling provider
		pageHandler.HandleFunc("/", myHandler)

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

	// ...

	// `aRequest` is the `*http.Request` argument passed to your handler func
	mySession := sessions.GetSession(aRequest)
	// ...
	myVal := mySession.Get("myKey")
	myVal2 := mySession.Get("myKey2")
	// do something with `myVal`/`myVal2`
	// ...
	otherVal := "important session value"
	mySession.Set("otherKey", otherVal)
	someNumber:= 123.456
	mySession.Set("numberKey", someNumber)
	// ...

Please note that both, the return value of `mySession.Get()` and the argument value of `mySession.Set()`, are defined/typed as `interface{}`.
That way you can store any value as session data.

## Internals

The package loads the sessions data (if any) whenever a page is requested and it stores the session data when the page handling is finished.
This is done automatically and you don't have to worry about loading/storing (read/write) of the session data.

The session ID (`SID`) is handled automatically as well.
Each ID is valid only for a single request by the remote user and changes for each request.
The name of the `SID` can be changed by calling

	sessions.SetSIDname(aSID string)

if the default (i.e. `SID`) doesn't satisfy your requirements.
To get the current setting you can call

	sid := sessions.SIDname()

The `SID` and the one-time-value are appended automatically as an [CGI argument](https://en.wikipedia.org/wiki/Common_Gateway_Interface) to all local `a href=""` links of the web page sent to the remote user, whereas _local_ means all links without a request scheme like `http:` or `https:`.

The package provides an internal garbage collector which deletes all expired sessions.
`Expired` are sessions when they were not touched/updated within the last _10 minutes_.
This default time can be changed by calling

	sessions.SetSessionTTL(aTTL int)

with `aTTL` seconds as the new time-to-life.
You can get the current TTL by calling

	ttl := sessions.SessionTTL()

## Licence

        Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                        All rights reserved
                    EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.
