# Sessions

[![GoDoc](https://godoc.org/github.com/mwat56/Nele?status.svg)](https://godoc.org/github.com/mwat56/sessions)
[![License](https://img.shields.io/eclipse-marketplace/l/notepad4e.svg)](https://github.com/mwat56/sessions/blob/master/LICENSE)

- [Sessions](#Sessions)
	- [Purpose / About](#Purpose--About)
	- [Installation](#Installation)
	- [Usage](#Usage)
	- [Licence](#Licence)

## Purpose / About

I wanted a session data solution that's user-friendly including privacy-friendly.
When doing some research about saving/retrieving sessions data with `Go` (aka `Golang`) you'll find some slightly different solution which have, however, one detail in common: they all depend on socalled internet `cookies`.
Which is bad.

_Cookies are bad_.

In practice `cookies` are basically an invasion into the user's property (`cookies` claim disk space) and they are, by definition, kind of a surveillance and tracking tool.
Nobody who has their user's best interesst in mind would ever even consider using `cookies`.
Their only advantage is that they are easy to implement – which was kind of the point: ease of implementation.
On the other hand the remote users were considered just a passive and obedient consumer.
It's clear that nowadays – with the user's privacy in mind – using `cookies` is simply an unacceptable technique.
It's also clear that harvesting the user's facilities (including disk space and electricity) should be avoided.

Another flaw you'll find in the literature about user sessions is the fact that it's often primarily considered in connetion with users who are in one way or another `logged in` into the web-server.
But that is only _one_ possible reason for considering some kind of session (data) store and – in my opinion – not even the most important one.
Others may include a navigation history (e.g. breadcrumbs), configuration data, individual preferences, etc.
In any case, session data should be able to transcent certain states of a user's connection.
And – that's a very important point – sessions should never require a user to `log in`.
The conception of the internet (more than a decade before that term became a synonym for commercial enterprise in the earlx 90s of the last centrury) was based on the idea of providing free access to as many data (i.e. knowledge) for as many people as possible.
Beside other things that means: session data should never be used as a kind of lock-in measure.

The ease of use for the application developer was mentioned before (when discussing the inherent badness of `cookies`).
That critique does not, however, mean that it should be difficult to implement session handling.
In fact, a huge part of the time I spent developing this package was spent figuring out an easy way to handle session data.
Now, after it's done, everything seems just 'logical', as if it's obvious how to do things.
That's an impression probably every programmer knows – at least if they do reflect what they are doing.

An other point to consider was to find a solution that's not dependent on third party facilities.
That most prominently excludes external database systems (like e.g. MariaDB and others alike).
The only database system to use is the OS's filesystem, the only system, that is, which doesn't add another possible point of failure.
And it is by definition faster than any other system because every other database system runs _on top_ of the filesystem and therefore adds some amount of overhead.
While it can be challenging and interesting to figure out some smart database structure and sophisticated queries to access and retrieve data such an endeavour is often simply over-the-top.
That's true especially for a task like session handling that has just two tasks to accomplish: point and grab (i.e. reading/loading data) and throw and forget (i.e. writing/storing data).

So, in short: In wanted a system that's unintrusive (i.e. respecting the user's privacy) and doesn't depend on anything but that what is there in any case (i.e. the filesystem).

    //TODO

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/sessions

## Usage

    //TODO

## Licence

        Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                        All rights reserved
                    EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.
