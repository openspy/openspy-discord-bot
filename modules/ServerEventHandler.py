import discord
class ServerEventHandler():
    guildSettings = None
    serverListHelper = None
    gameInfo = None
    def __init__(self, bot, serverListHelper, gameInfo, guildSettings):
        self.guildSettings = guildSettings
        self.serverListHelper = serverListHelper
        self.gameInfo = gameInfo
        self.bot = bot
    def GetServerInfoEmbedMsg(self, event_type, server_info, game_info):
        embedMsg = discord.Embed(title="New Server", description=game_info["description"], color=0x00ff00)
        embedMsg.set_footer(text="Use !serverinfo [server key] for detailed server information")
        embedMsg.add_field(name="Server Key", value=server_info["server_key"], inline=True)
        embedMsg.add_field(name="Hostname", value=server_info["custkeys"]["hostname"], inline=True)
        country_string = "Unknown"
        if "country" in server_info:
            country_string = "{} - :flag_{}:".format(server_info["country"], server_info["country"].lower())
        embedMsg.add_field(name="Country", value=country_string, inline=True)
        embedMsg.add_field(name="Players", value="{}/{}".format(server_info["custkeys"]["numplayers"], server_info["custkeys"]["maxplayers"]), inline=True)
        embedMsg.add_field(name="Mapname", value=server_info["custkeys"]["mapname"], inline=True)
        embedMsg.add_field(name="Gamemode", value=server_info["custkeys"]["gamemode"], inline=True)
        return embedMsg
    async def OnServerEvent(self, serverListCache, event_type, server_key):
        
        if event_type != "new":
            return

        
        server_info = self.serverListHelper.GetServerInfo(server_key, serverListCache)

        if server_info == None:
            return

        game_info = self.gameInfo.GetGameInfoByGameid(server_info["gameid"])
        if game_info == None:
            return
        embedMsg = self.GetServerInfoEmbedMsg(event_type, server_info, game_info)
        for guild in self.bot.guilds:
            all_channel = self.guildSettings.GetServerEventChannel(guild, event_type, "all")
            if all_channel != None:
                channel = self.bot.get_channel(all_channel)
                if channel != None:
                    await channel.send(embed=embedMsg)

            game_channel = self.guildSettings.GetServerEventChannel(guild, event_type, game_info["gamename"])
            if game_channel != None and game_channel != all_channel:
                channel = self.bot.get_channel(game_channel)
                if channel != None:
                    await channel.send(embed=embedMsg)