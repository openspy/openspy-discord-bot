import redis
class GameInfo():
    gameInfoPool = None
    def __init__(self, gameInfoPool):
        self.gameInfoPool = gameInfoPool
    def GetGameInfoByGameid(self, gameid):
        game_info = {}

        gameInfoCache = redis.Redis(connection_pool=self.gameInfoPool)

        gameid_key = "gameid_{}".format(gameid)
        game_key = gameInfoCache.get(gameid_key)

        game_info["id"] = int(gameid)
        game_info["gamename"] = gameInfoCache.hget(game_key, "gamename").decode('utf8')
        game_info["secretkey"] = gameInfoCache.hget(game_key, "secretkey").decode('utf8')
        game_info["description"] = gameInfoCache.hget(game_key, "description").decode('utf8')
        game_info["queryport"] = gameInfoCache.hget(game_key, "queryport").decode('utf8')
        return game_info