# Honoroit [![Matrix](https://img.shields.io/matrix/honoroit:etke.cc?logo=matrix&style=for-the-badge)](https://matrix.to/#/#honoroit:etke.cc) [![Buy me a Coffee](https://shields.io/badge/donate-buy%20me%20a%20coffee-green?logo=buy-me-a-coffee&style=for-the-badge)](https://buymeacoffee.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/honoroit/badges/main/coverage.svg)](https://gitlab.com/etke.cc/honoroit/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/honoroit)](https://goreportcard.com/report/gitlab.com/etke.cc/honoroit) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/honoroit.svg)](https://pkg.go.dev/gitlab.com/etke.cc/honoroit)

> [more about that name](https://finalfantasy.fandom.com/wiki/Honoroit_Banlardois)

A helpdesk bot, used as part of [etke.cc](https://etke.cc) service.

The main idea of that bot is to give you the same abilities as with website chats (like Intercom, jivosite, etc) inside the matrix.

## Features

* Get a message from any matrix user proxied to a specific room. Any message that user will send in his 1:1 room with Honoroit will be proxied as thread messages
* Chat with that user through the honoroit bot in a thread inside your special room. Any member of that special room can participate in discussion
* When request fulfilled - send a `@honoroit done` (with mention) in that thread - thread topic will be renamed and "proxied user" will know that request was closed (bot will leave user's room with special notice)

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

* Email<->Matrix helpdesk
* End-to-End Encryption
* autoleave empty rooms and notify about that (requires persistent store)

## Commands

available commands in the threads. Note that all commands should be called with mention of honoroit, so `@honoroit done` will work, but simple `done` will not.

* `done` - close the current request and mark is as done. Customer will receive special message and honoroit bot will leave 1:1 chat with customer. Any new message to the thread will not work and return error.
* `rename` - rename the thread topic title, when you want to change the standard message to something different


## Configuration

env vars

### mandatory

* **HONOROIT_HOMESERVER** - homeserver url, eg: `https://matrix.example.com`
* **HONOROIT_LOGIN** - user login/localpart, eg: `honoroit`
* **HONOROIT_PASSWORD** - user password
* **HONOROIT_ROOMID** - room ID where threads will be created, eg: `!test:example.com`

### optional

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
