"""
_.-={ roles.py }=-._

Cog containing get/set role commands for general users with no role privs.
"""


import discord
from discord.ext import commands


class RolesCog(commands.Cog):

    def __init__(self, bot):
        self.bot = bot
        p = discord.Permissions.none()
        # Role names SHOULD be unique; the bot breaks if they aren't.
        self.roles = {str(r): r
                      for r in self.bot.guild.roles
                      if r.permissions <= p}

    @commands.command(name='roles', aliases=['ranks'])
    @commands.guild_only()
    async def list_roles(self, ctx):
        """Lists every enlistable roles with no roles permissions"""
        await ctx.send('```' + '\n'.join(self.roles.keys()) + '```')\

    @commands.command(name='role', aliases=['r', 'rank'])
    @commands.guild_only()
    async def add_roles(self, ctx, *roles: str):
        """Tries to add/remove roles depending on the roles args and whether the
        user already has to role or not. A separate list is made for roles not
        found."""
        auth = ctx.author
        added, removed, error = [], [], []
        output = ''
        for r in roles:
            try:
                role = self.roles[r]
            except KeyError:
                error.append(r)
            else:
                if role in ctx.author.roles:
                    removed.append(r)
                    await auth.remove_roles(role, reason='add_roles invoked')
                else:
                    added.append(r)
                    await ctx.author.add_roles(role, reason='add_roles invoked')

        if added:
            output += f"Roles added: {', '.join(added)}"
        if removed:
            output += f"\nRoles removed: {', '.join(removed)}"
        if error:
            output += f"\nRoles not found: {', '.join(error)}"

        await ctx.send(output)

    @commands.Cog.listener()
    async def on_guild_role_create(self, role):
        """Every time a basic role is added, the bot adds to its role list."""
        if role.guild == self.bot.guild_id and \
                role.permissions <= discord.Permissions.none():
            self.roles[role.name] = role


def setup(bot):
    bot.add_cog(RolesCog(bot))
