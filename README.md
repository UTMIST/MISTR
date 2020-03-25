# UTMIST Assistant (MISTA)

![UTMIST Logo](logo.png)

The UofT Machine Intelligence Student Team Assistant (MISTA) is a Discord bot serving the [UTMIST Discord server](https://discord.gg/88mSPw8).

## Overview

- MISTA can manage server roles and permissions in our [Discord server](https://discord.gg/88mSPw8) using message reacts.
- Our [website (utmist.gitlab.io)](https://utmist.gitlab.io) can be updated via MISTA on our server by triggering jobs for GitLab Pages CI running [utmist.gitlab.io](https://gitlab.com/utmist/utmist.gitlab.io).
- MISTA also has Event, Prejct, and Media interfaces to provide resources and information on our server.

### Prerequisites

- [Go](https://golang.org/).

### Dependencies

- [discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [go-gitlab](https://pkg.go.dev/github.com/xanzy/go-gitlab)

## Details

- Run `sh go-get.sh` to download dependencies.
- The `.env` file in the form of `.env.copy` is required to run this bot.
- The bot is entirely run from `main.go`.
- [Discord Invite Link](https://discordapp.com/oauth2/authorize?client_id=682495255102095391&scope=bot).

## Development

- This project is maintained by the [Infrastructure Department at UTMIST](https://utmist.gitlab.io/team/infrastructure).
  - [Robbert Liu](https://github.com/triglemon), Infrastructure Developer & MISTA Lead.
  - [Robert (Rupert) Wu](https://leglesslamb.gitlab.io), VP Infrastructure.
- If you're a member of UTMIST and would like to contribute or learn development through this project, you can join our [Discord](https://discord.gg/88mSPw8)) and let us know in _#infrastructure_.

## Acknowledgements

- Adela Hua ([@makurophage](https://www.instagram.com/makurophage/)) for the [MISTA logo](https://gitlab.com/utmist/mista/-/blob/master/logo.png).
