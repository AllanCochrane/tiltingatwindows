Tiltingatwindmills
==================

A small Go app to show user's tweets and find shared followers

Requirements
------------

[Go](http://golang.org)

[Twitter API credentials](https://apps.twitter.com/) (TODO: This might not be the right link ...)

Installation
------------

Install the packages that this application relies on:

``` shell

export GOPATH=<wherever you wish Go libs to reside>

go get gopkg.in/gin-gonic/gin.v1
go get github.com/ChimeraCoder/anaconda

```

Set the required environment variables with the credentials as
supplied when you create a Twitter app account:

``` shell

export TW_ACCESS_TOKEN_SECRET=<twitter secret access token>
export TW_ACCESS_TOKEN=<twitter access token>
export TW_CONSUMER_SECRET=<twitter secret consumer key>
export TW_CONSUMER_KEY=<twitter consumer key>

```

Invocation
----------

Use `go run` to start the app:

``` shell

go run tiltingatwindmils.go

```

The server runs on localhost port 8000, documentation is at:

http://localhost:8000/docs
