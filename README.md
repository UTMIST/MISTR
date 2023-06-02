# UTMIST Runner (MISTR)

![Minimum go version](https://img.shields.io/badge/go-1.11+-important.svg)
![Go modules version](https://img.shields.io/github/go-mod/go-version/utmist/MISTR/master)

The UofT Machine Intelligence Student Team Runner (MISTR) is a Discord bot serving the [UTMIST Discord server](https://discord.gg/88mSPw8). Our [website (utmist.gitlab.io)](https://utmist.gitlab.io) can be updated via MISTR on our server by triggering jobs for GitLab Pages CI running [utmist.gitlab.io](https://gitlab.com/utmist/utmist.gitlab.io). Other GitLab CI tasks such as batch flushing can also be performed.

### Prerequisites

- [Go](https://golang.org/): minimum 1.11+; recommended 1.13+.

### Dependencies

- [discordgo](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [go-gitlab](https://pkg.go.dev/github.com/xanzy/go-gitlab)

## Details

- Dependencies are managed with [Go modules](https://blog.golang.org/using-go-modules).
- The `.env` file with the variables listed in `copy.env` is required to run this bot.
- The bot is run from `main.go`.
- [Discord Invite Link](https://discordapp.com/oauth2/authorize?client_id=682495255102095391&scope=bot).

### Continuous Deployment (CD)

If you're looking to deploy MISTR continuously on a server or a local machine, run `loop.sh` in the project root directory. Then when you want to restart with new commits, type `mistr! restart` in the appropriate Discord channel.

UTMIST runs MISTR on a Raspberry Pi (of various models) using `loop.sh`.

## Development

- This project is maintained by the [Engineering Department at UTMIST](https://utmist.gitlab.io/team/infrastructure).
- If you're a member of UTMIST and would like to contribute or learn development through this project, you can join our [Discord](https://discord.gg/88mSPw8)) and let us know in _#infrastructure_.

## Acknowledgements

- Adela Hua ([@makurophage](https://www.instagram.com/makurophage/)) for the [MISTR logo](https://gitlab.com/utmist/mistr/-/blob/master/logo.png).
