# Honoroit [![Matrix](https://img.shields.io/matrix/honoroit:etke.cc?logo=matrix&style=for-the-badge)](https://matrix.to/#/#honoroit:etke.cc)[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/honoroit/badges/main/coverage.svg)](https://gitlab.com/etke.cc/honoroit/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/honoroit)](https://goreportcard.com/report/gitlab.com/etke.cc/honoroit) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/honoroit.svg)](https://pkg.go.dev/gitlab.com/etke.cc/honoroit)

> [more about that name](https://finalfantasy.fandom.com/wiki/Honoroit_Banlardois)

A helpdesk bot, used as part of [etke.cc](https://etke.cc) service.

The main idea of that bot is to give you the same abilities as with website chats (like Intercom, jivosite, etc) inside the matrix.

## Features

* End-to-End encryption
* Get a message from any matrix user proxied to a specific room. Any message that user will send in his 1:1 room with Honoroit will be proxied as thread messages
* Chat with that user through the honoroit bot in a thread inside your special room. Any member of that special room can participate in discussion
* When request fulfilled - send a `!ho done` in that thread - thread topic will be renamed and "proxied user" will know that request was closed (bot will leave user's room with special notice)

## Warning

Honoroit utilizes [MSC 3440: Threading](https://github.com/matrix-org/matrix-doc/pull/3440/) which is:

* not stable yet
* not part of the protocol specification yet
* requires a compatible client app (e.g. Element with `Threads` enabled in `Labs`)

Of course, there is fallback mode with reply-to (supported by any matrix server and client):
honoroit will try to find mappings, going recursively through reply-to chain,
but that mode considered as workaround and it's perfomance is poor.

Honoroit used as support service in the [etke.cc](https://etke.cc) and works pretty good in daily operations,
so just keep in mind the limitations above.

## How it looks like

<details>
<summary>Screenshots</summary>

### Step 1: a matrix user (customer) sends a message to Honoroit bot in direct 1:1 chat

![step 1](contrib/screenshots/1.customer sends a message.png)

### Step 2: a new thread created in the backoffice room

![step 2](contrib/screenshots/2.a new thread created in the backoffice room.png)

### Step 3: operator(-s) chat with customer in that thread

![step 3](contrib/screenshots/3.operators chat with customer in that thread.png)

### Step 4: customer sees that like a direct 1:1 chat with honoroit user

![step 4](contrib/screenshots/4.customer sees that like a direct 1:1 chat with honoroit user.png)

### Step 5: operator closes the request

![step 5](contrib/screenshots/5.operator closes the request.png)

### Step 6: customer receives special message and bot leaves the room

![step 6](contrib/screenshots/6.customer receives special message and bot leaves the room.png)

</details>

## TODO

* Unit tests

## Commands

available commands in the threads. Note that all commands should be called with prefix, so `!ho done` will work, but simple `done` will not.

* `done` - close the current request and mark is as done. Customer will receive special message and honoroit bot will leave 1:1 chat with customer. Any new message to the thread will not work and return error.
* `rename` - rename the thread topic title, when you want to change the standard message to something different
* `invite` - invite yourself into the customer 1:1 room


## Configuration

env vars

### mandatory

* **HONOROIT_HOMESERVER** - homeserver url, eg: `https://matrix.example.com`
* **HONOROIT_LOGIN** - user login/localpart, eg: `honoroit`
* **HONOROIT_PASSWORD** - user password
* **HONOROIT_ROOMID** - room ID where threads will be created, eg: `!test:example.com`

### optional

* **HONOROIT_PREFIX** - command prefix
* **HONOROIT_SENTRY** - sentry DSN
* **HONOROIT_LOGLEVEL** - log level
* **HONOROIT_DB_DSN** - database connection string
* **HONOROIT_DB_DIALECT** - database dialect (postgres, sqlite3)
* **HONOROIT_TEXT_GREETINGS** - a message sent to customer on first contact
* **HONOROIT_TEXT_ERROR** - a message sent to customer if something goes wrong
* **HONOROIT_TEXT_EMPTYROOM** - a message sent to backoffice/threads room when customer left his room
* **HONOROIT_TEXT_DONE** - a message sent to customer when request marked as done in the threads room

You can find default values in [config/defaults.go](config/defaults.go)

## Where to get

[docker registry](https://gitlab.com/etke.cc/honoroit/container_registry), [etke.cc](https://etke.cc)
