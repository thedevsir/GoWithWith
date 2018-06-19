# GoWithWith API Build on JSON Web Tokens

A user system API starter. Bring your own front-end.

[![Build Status](https://travis-ci.org/Gommunity/GoWithWith.svg?branch=master)](https://travis-ci.org/Gommunity/GoWithWith)
[![Dependency Status](https://david-dm.org/Gommunity/GoWithWith.svg?style=flat)](https://david-dm.org/Gommunity/GoWithWith)
[![devDependency Status](https://david-dm.org/Gommunity/GoWithWith/dev-status.svg?style=flat)](https://david-dm.org/Gommunity/GoWithWith#info=devDependencies)


## Features (we working on it)

 - Sign up system with verification email
 - Login system with forgot password and reset password
 - Abusive login attempt detection
 - Session management system

## Responsive HTML e-mails

[Hermes](https://github.com/matcornic/hermes)  is the Go part of the great [mailgen](https://github.com/eladnava/mailgen) engine for Node.js. Check their work, it's awesome! It's a package that generates clean, responsive HTML e-mails and associated plain text fallback.

<img src="https://raw.githubusercontent.com/matcornic/hermes/master/screens/default/welcome.png" height="400" /> <img src="https://raw.githubusercontent.com/matcornic/hermes/master/screens/default/reset.png" height="400" /> <img src="https://raw.githubusercontent.com/matcornic/hermes/master/screens/default/receipt.png" height="400" />


## Technology

Fresh is built with the [Echo 3.3](https://echo.labstack.com/) framework. We're
using [MongoDB](http://www.mongodb.org/) as a data store.

## Bring your own front-end

GoWithWith is only a restful JSON API. If you'd like a ready made front-end for clients,
checkout [Hexagenal](https://github.com/escommunity/hexagonal). Or better yet, fork
this repo and build one on top of GoWithWith.

## Requirements

You need [Golang](https://golang.org/) `>=1.5.x` and you'll need a
[MongoDB](http://www.mongodb.org/downloads) `>=2.6` server running.

## Installation

```bash
$ git clone https://github.com/Gommunity/GoWithWith.git
$ cd GoWithWith
// install packages with go tools
```

## Configuration


Simply copy `.env-sample` to `.env` and edit as needed. __Don't commit `.env`
to your repository.__

## Running the app

```bash
$ go run !(*_test).go

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.dev
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:3500
```

Now you should be able to point your browser to http://127.0.0.1:3500/swagger/index.html and
see the documentation page.

## Running in production
It's soon yet !

## Have a question?

Any issues or questions (no matter how basic), open an issue. Please take the
initiative to read relevant documentation and be pro-active with debugging.


## Want to contribute?

Contributions are welcome. If you're changing something non-trivial, you may
want to submit an issue before creating a large pull request.



## License

MIT


## Don't forget

What you build with GoWithWith is more important than GoWithWith. 
