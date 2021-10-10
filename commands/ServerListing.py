import redis
import discord
from discord.ext import commands
class ServerListingCommands(commands.Cog):
    serverListHelper = None
    serverListPool = None
    gameInfo = None
    def __init__(self, bot, serverListHelper, serverListPool, gameInfo):
        self.serverListHelper = serverListHelper
        self.serverListPool = serverListPool
        self.gameInfo = gameInfo
        self.bot = bot
    
    def embedServer(self, server_info, embedMsg, field_index):
        embedMsg.insert_field_at(index=field_index,name="Server Key", value=server_info["server_key"], inline=False)
        field_index += 1
        embedMsg.insert_field_at(index=field_index,name="Country", value=server_info["country"], inline=True)
        field_index += 1

        display_keys = ["hostname", "numplayers", "maxplayers", "mapname"]
        for key in display_keys:
            if key in server_info["custkeys"]:
                embedMsg.insert_field_at(index=field_index,name=key, value=server_info["custkeys"][key], inline=True)
                field_index += 1

        return field_index + 2

    @commands.command()
    async def servers(self, ctx, gamename):
        if not self.gameInfo.GamenameExists(gamename):
            await ctx.send("Gamename **{}** does not exist.".format(gamename))
            return

        servers = self.serverListHelper.FetchServerList(gamename)

        desc = "Server Info - {} - Total servers: {}\n".format(gamename, len(servers))
        embedMsg = discord.Embed(title="Server List",description=desc, color=0x0000ff, url="http://beta.openspy.net/en/server-list/{}".format(gamename))
        embedMsg.set_footer(text="Use !serverinfo [server key] for detailed server information")

        MAX_SERVERS_PER_MESSAGE = 4
        server_index = 0
        for server in servers:
            embedMsg.add_field(name="Server Key", value=server["server_key"], inline=True)
            embedMsg.add_field(name="Hostname", value=server["custkeys"]["hostname"], inline=True)
            country_string = "Unknown"
            if "country" in server:
                country_string = "{} - :flag_{}:".format(server["country"], server["country"].lower())
            embedMsg.add_field(name="Country", value=country_string, inline=True)
            embedMsg.add_field(name="Players", value="{}/{}".format(server["custkeys"]["numplayers"], server["custkeys"]["maxplayers"]), inline=True)
            embedMsg.add_field(name="Mapname", value=server["custkeys"]["mapname"], inline=True)
            embedMsg.add_field(name="Gamemode", value=server["custkeys"]["gamemode"], inline=True)

            server_index += 1

            if server_index >= MAX_SERVERS_PER_MESSAGE:
                server_index = 0
                await ctx.send(embed=embedMsg)
                embedMsg.clear_fields()

        if server_index != 0 or len(servers) == 0:
            await ctx.send(embed=embedMsg)


    @commands.command()
    async def serverinfo(self, ctx, server_key):
        serverListCache = redis.Redis(connection_pool=self.serverListPool)
        server_info = self.serverListHelper.GetServerInfo(server_key, serverListCache)

        embedMsg = discord.Embed(title="Server Info",description="Server Values", color=0x0000ff)
        embedMsg.add_field(name="Server Key", value=server_info["server_key"], inline=True)

        if "num_probes" in server_info and int(server_info["num_probes"]) > 0:
            ports_forwarded = "allow_unsolicited_udp" in server_info and bool(server_info["allow_unsolicited_udp"])
            embedMsg.add_field(name="Port Forwarded", value=str(ports_forwarded), inline=True)

        country_string = "Unknown"
        if "country" in server_info:
            country_string = "{} - :flag_{}:".format(server_info["country"], server_info["country"].lower())
        embedMsg.add_field(name="Country", value=country_string, inline=True)

        hide_keys = ["publicip", "publicport", "statechanged", "hostport", "localport", "localip0", "localip1", "localip2", "localip3", "localip4", "localip5", "localip6", "localip7", "queryid", "final"]

        for key, value in server_info["custkeys"].items():
            if not key in hide_keys:
                    if len(key) > 0 and len(value) > 0:
                        embedMsg.add_field(name=key, value=value, inline=True)
        await ctx.send(embed=embedMsg)

        if len(server_info["teamkeys"]) > 0:
            embedMsg = discord.Embed(title="Server Info - Team",description="Team Values", color=0x0000ff)
            team_index = 1
            for team_info in server_info["teamkeys"]:
                embedMsg.add_field(name="Team Index", value="{}".format(team_index), inline=False)
                for key, value in team_info.items():
                    if len(key) > 0 and len(value) > 0:
                        embedMsg.add_field(name=key, value=value, inline=True)                
                team_index += 1
            await ctx.send(embed=embedMsg)

        if len(server_info["playerkeys"]) > 0:
            embedMsg = discord.Embed(title="Server Info - Players",description="Player Values", color=0x0000ff)
            player_index = 1
            for player_info in server_info["playerkeys"]:
                embedMsg.add_field(name="Player Index", value="{}".format(player_index), inline=False)
                for key, value in player_info.items():
                    if len(key) > 0 and len(value) > 0:
                        embedMsg.add_field(name=key, value=value, inline=True)                
                player_index += 1
            await ctx.send(embed=embedMsg)