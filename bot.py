"""
_.-={ UTDISC }=-._

This is the main bot module. When run as __main__, it will attempt to load every
cog in initial_extensions and will safely error if it isn't found.

Prefix: >
"""


import discord
from discord.ext import commands


with open('token') as t:
    token = t.read()

initial_extensions = ['cogs.roles'
                      'cogs.owner']


class UTBot(commands.Bot):
    """
    guild: UTMist discord server
        type: Optional[discord.Guild]
    """
    def __init__(self):
        super().__init__(command_prefix='>',
                         description="UTMIST's personal bot")
        self.guild_id = 673778422291628033


bot = UTBot()

if __name__ == '__main__':
    for extension in initial_extensions:
        bot.load_extension(extension)

# Event decorator defines lists of actions for preset events like "on_ready"
# Anything with a bot or Commands decorator should not return anything.


@bot.check
def check_commands(ctx):
    return ctx.guild.id == bot.guild_id


@bot.event
async def on_ready():
    """
    After initializing period, bot is "ready" and will print some data and
    change the bot's activity.
    """
    print(f'\n{bot.user.name=}\n{bot.user.id=}\n{discord.__version__=}\n')
    await bot.change_presence(activity=discord.Game('loafing around'))
    print(f'Logged in :^)')


bot.run(token, bot=True, reconnect=True)
