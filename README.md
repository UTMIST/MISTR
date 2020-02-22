# Mista

![logo.png](logo.png)

The UofT Machine Intelligence Student Team Assistant (MISTA) is a Discord bot serving the UTMIST Discord Server. It will have three main responsibilities.

- Managing roles and permissions within the server.
- Providing timely information about events and resources for club members.
- Trigger jobs for the GitLab Pages CI running [utmist.gitlab.io](https://gitlab.com/utmist/utmist.gitlab.io).

## Setup/Housekeeping

- The setup for your `pip` packages depends on how you invoke `pip` and whether you use aliases.
  - Run `sh pip3.sh` if you invoke `pip` for Python 3.7+ using `pip3` (particularly for Debian/Ubuntu).
  - Run `sh pip.sh` otherwise.
- `bot.py` is the "main" module controlling `owner.py` and `roles.py`.

## Prerequisites

- [Python](https://www.python.org/), `>=3.7`.
- [Pip (3)](https://pip.pypa.io/).

## Inviting the Bot

Discord invite [link](https://discordapp.com/oauth2/authorize?client_id=679066276047486995&scope=bot&permissions=32).

## Developers

- [Robbert Liu](https://gitlab.com/triglemon)
- [Rupert Wu](https://leglesslamb.gitlab.io)
