import redis
import discord
from discord.ext import commands
class ServerListingEventConfigCommands(commands.Cog):
    bot = None
    guildSettings = None
    def __init__(self, bot, guildSettings):
        self.bot = bot
        self.guildSettings = guildSettings

    @commands.command()
    @commands.has_permissions(administrator=True)
    @commands.guild_only()
    async def newserverevents(self, ctx, gamename=None):
        if gamename == None:
            self.guildSettings.SetServerEventChannel(ctx.guild, ctx.channel.id, "new", "all")
            await ctx.send("Events for **all** servers are now in this channel")
        else:
            self.guildSettings.SetServerEventChannel(ctx.guild, ctx.channel.id, "new", gamename)
            await ctx.send("Events for **{}** are now in this channel".format(gamename))

    @commands.command()
    @commands.has_permissions(administrator=True)
    @commands.guild_only()
    async def stopevents(self, ctx, gamename=None):
        if gamename == None:
            self.guildSettings.DeleteChannelEvents(ctx.guild, ctx.channel.id)
            await ctx.send("**All** events stoppped in this channel")
        else:
            self.guildSettings.DeleteChannelGamenameEvents(ctx.guild, ctx.channel.id, gamename)
            await ctx.send("Events for **{}** are now stopped in this channel".format(gamename))