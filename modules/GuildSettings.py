class GuildSettings():
    cacheItems = None
    def __init__(self, dbCtx):
        self.dbCtx = dbCtx
        self.cacheItems = {}

    def GetCachedServerEvent(self, guild, event_type, gamename):
        cache_key = "server_event_{}_{}_{}".format(guild, event_type, gamename)

        if cache_key in self.cacheItems:
            return self.cacheItems[cache_key]

        return None

    def SetCachedServerEvent(self, guild, event_type, gamename, channel_id):
        cache_key = "server_event_{}_{}_{}".format(guild, event_type, gamename)
        self.cacheItems[cache_key] = channel_id

    def DeleteCachedServerEvent(self, guild, event_type, gamename, channel_id):
        cache_key = "server_event_{}_{}_{}".format(guild, event_type, gamename)
        del self.cacheItems[cache_key]

    def GetServerEventChannel(self, guild, event_type, gamename):
        cache_result = self.GetCachedServerEvent(guild, event_type, gamename)
        if cache_result != None:
            return cache_result

        collection = self.dbCtx["server_new_event_channels"]

        match = {"gamename": gamename, "type": event_type,"guild": guild.id}

        result = collection.find_one(match)
        if result == None:
            return None

        self.SetCachedServerEvent(guild, event_type, gamename, result["channel"])
        return result["channel"]
    def SetServerEventChannel(self, guild, channel_id, event_type, gamename):
        self.SetCachedServerEvent(guild, event_type, gamename, channel_id)

        collection = self.dbCtx["server_new_event_channels"]

        match = {"gamename": gamename, "type": event_type,"guild": guild.id}
        update_statement =  { "$set": {"channel": channel_id}}

        collection.update_one(match, update_statement, upsert=True)
    def DeleteChannelEvents(self, guild, channel_id):
        self.cacheItems = {} #should scan for partial matches, but oh well

        collection = self.dbCtx["server_new_event_channels"]
        match = {"guild": guild.id, "channel": channel_id}
        collection.delete_many(match)

    def DeleteChannelGamenameEvents(self, guild, channel_id, gamename):
        self.DeleteChannelGamenameEvents(guild, channel_id, gamename)
        collection = self.dbCtx["server_new_event_channels"]
        match = {"guild": guild.id, "channel": channel_id, "gamename": gamename}
        collection.delete_many(match)
