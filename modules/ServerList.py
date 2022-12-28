import redis
class ServerList():
    serverListPool = None
    def __init__(self, serverListPool):
        self.serverListPool = serverListPool
    def GetServerInfo(self, server_key, serverListCache):
        server_info = {}
        server_custkeys = {}
        server_player_keys = []
        server_team_keys = []
        if not serverListCache.exists(server_key):
            return None
        for hscan_key in serverListCache.hscan_iter(server_key):
            try:
                key = hscan_key[0].decode('utf8', 'replace')
                value = hscan_key[1].decode('utf8', 'replace')
                server_info[key] = value
            except:
                pass

        for hscan_key in serverListCache.hscan_iter(server_key + "custkeys"):
            try:
                key = hscan_key[0].decode('utf8', 'replace')
                value = hscan_key[1].decode('utf8', 'replace')
                server_custkeys[key] = value
            except:
                pass

        for player_key in serverListCache.scan_iter(server_key + "custkeys_player_*"):
            try:
                key = player_key.decode('utf8', 'replace')
                player_keys = {}
                for hscan_key in serverListCache.hscan_iter(key):
                    player_hkey = hscan_key[0].decode('utf8', 'replace')
                    player_hvalue = hscan_key[1].decode('utf8', 'replace')
                    player_keys[player_hkey] = player_hvalue
                server_player_keys.append(player_keys)
            except:
                pass

        for team_key in serverListCache.scan_iter(server_key + "custkeys_team_*"):
            try:
                key = team_key.decode('utf8', 'replace')
                team_keys = {}
                for hscan_key in serverListCache.hscan_iter(key):
                    team_hkey = hscan_key[0].decode('utf8', 'replace')
                    team_hvalue = hscan_key[1].decode('utf8', 'replace')
                    team_keys[team_hkey] = team_hvalue
                server_team_keys.append(team_keys)
            except:
                pass

        
        server_info["custkeys"] = server_custkeys
        server_info["playerkeys"] = server_player_keys
        server_info["teamkeys"] = server_team_keys
        server_info["server_key"] = server_key

        server_info = self.handleVariableRemapping(server_info)
        return server_info
    def FetchServerList(self, gamename):
        serverListCache = redis.Redis(connection_pool=self.serverListPool)

        serverInfoLookupCache = redis.Redis(connection_pool=self.serverListPool)

        servers = []
        for server_key in serverListCache.zscan_iter(gamename):
            key = server_key[0].decode('utf8')
            server_info = self.GetServerInfo(key, serverInfoLookupCache)
            deleted = False
            if server_info != None:
                if "deleted" in server_info:
                    deleted = int(server_info["deleted"]) == 1
                if not deleted:
                    servers.append(server_info)
        return servers
            
    def handleVariableRemapping(self, server_info):
        if "hostname" not in server_info["custkeys"]:
            server_info["custkeys"]["hostname"] = "Unknown"
        if "gamemode" not in server_info["custkeys"]:
            server_info["custkeys"]["gamemode"] = "Unknown"
        if "mapname" not in server_info["custkeys"]:
            server_info["custkeys"]["mapname"] = "Unknown"

        if "numplayers" not in server_info["custkeys"]:
            server_info["custkeys"]["numplayers"] = "-1"
        if "maxplayers" not in server_info["custkeys"]:
            server_info["custkeys"]["maxplayers"] = "-1"

        gameid = int(server_info["gameid"])

        sr2_gameids = [2108,2109,2110,2112,2174]
        if gameid in sr2_gameids:
            if "online_name" in server_info["custkeys"]:
                server_info["custkeys"]["hostname"] = server_info["custkeys"]["online_name"]
            if "open_slots" in server_info["custkeys"] and "maxplayers" in server_info["custkeys"]:
                server_info["custkeys"]["numplayers"] = str(int(server_info["custkeys"]["maxplayers"]) - int(server_info["custkeys"]["open_slots"]))
            del server_info["custkeys"]["session"]
            if "mode" in server_info["custkeys"]:
                mode = int(server_info["custkeys"]["mode"])
                if mode == 3:
                    server_info["custkeys"]["gamemode"] = "Co-op"
                else:
                    ranked = "ranked" in server_info["custkeys"] and int(server_info["custkeys"]["ranked"]) == 1
                    if ranked:
                        server_info["custkeys"]["gamemode"] = "Multiplayer (ranked)"
                    else:
                        server_info["custkeys"]["gamemode"] = "Multiplayer"
        return server_info