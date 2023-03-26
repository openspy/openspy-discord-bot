# bot.py
import os
import redis

import pymongo
import discord

from modules.ServerList import ServerList
from modules.ServerEventHandler import ServerEventHandler
from modules.GuildSettings import GuildSettings
from modules.GameInfo import GameInfo
from commands.ServerListing import ServerListingCommands
from commands.ServerListingEventConfig import ServerListingEventConfigCommands

import asyncio
from aio_pika import connect, IncomingMessage

mongoConnection = pymongo.MongoClient(os.environ.get('MONGODB_URI'))
discordBotCollection = mongoConnection["DiscordBot"]

serverListPool = redis.ConnectionPool.from_url(os.environ.get('REDIS_URL'), db=0)
gameInfoPool = redis.ConnectionPool.from_url(os.environ.get('REDIS_URL'), db=2)

TOKEN = os.getenv('DISCORD_TOKEN')

client = discord.ext.commands.Bot(command_prefix='!', intents=discord.Intents.all())

gameInfo = GameInfo(gameInfoPool)
serverList = ServerList(serverListPool)
guildSettings = GuildSettings(discordBotCollection)
eventHandler = ServerEventHandler(client, serverList, gameInfo, guildSettings)

async def server_event_callback(message: IncomingMessage):
    serverListCache = redis.Redis(connection_pool=serverListPool) #should we be doing this each event?
    event_string = message.body.decode('utf8').split('\\')
    await eventHandler.OnServerEvent(serverListCache, event_string[1], event_string[2])

@client.event
async def on_ready():
    print(f'{client.user} has connected to Discord!')

    bot_activity = discord.Activity(type=discord.ActivityType.custom, name="http://openspy.net")
    await client.change_presence(status=discord.Status.online, activity=bot_activity)

    # Perform connection
    connection = await connect(os.environ.get('AMQP_URL'))

    # Creating a channel
    channel = await connection.channel()

    # Declaring queue
    queue = await channel.declare_queue(exclusive=True)
    await queue.bind("openspy.master", routing_key="server.event")

    await queue.consume(server_event_callback, no_ack=True)

async def reg_commands():
    await client.add_cog(ServerListingCommands(client, serverList, serverListPool, gameInfo))
    await client.add_cog(ServerListingEventConfigCommands(client, guildSettings))

asyncio.run(reg_commands())
client.run(TOKEN)
