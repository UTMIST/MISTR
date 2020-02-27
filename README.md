# UTMIST Assistant (MISTA)

![logo.png](logo.png)

The UofT Machine Intelligence Student Team Assistant (MISTA) is a Discord bot serving the UTMIST Discord Server. It will have three main responsibilities.

- Trigger jobs for the GitLab Pages CI running [utmist.gitlab.io](https://gitlab.com/utmist/utmist.gitlab.io) (**complete**). Documentation to issue this command will follow after we build permissions.
- Managing roles and permissions within the server (**in progress**).
- Providing timely information about events and resources for club members (**in progress**).

## Prerequisites

- [Go](https://golang.org/).
- See [Go-Gitlab Docs](https://godoc.org/github.com/xanzy/go-gitlab)

## Setup/Housekeeping

- Run `sh go-get.sh` to downloade dependencies.
- The `.env` file in the form of `.env.copy` is required to run this bot.
- The bot is entirely run from `main.go`.

## Inviting the Bot

Discord invite [link](https://discordapp.com/oauth2/authorize?client_id=682495255102095391&scope=bot).

## Developers

- [Robbert Liu](https://gitlab.com/triglemon)
- [Rupert Wu](https://leglesslamb.gitlab.io)
