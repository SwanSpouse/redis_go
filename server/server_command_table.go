package server

import (
	"fmt"

	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/handlers"
)

/*
 * This is the meaning of the flags:
 *
 * w: write command (may modify the key space).
 *    写入命令，可能会修改 key space
 *
 * r: read command  (will never modify the key space).
 *    读命令，不修改 key space
 *
 * m: may increase memory usage once called. Don't allow if out of memory.
 *    可能会占用大量内存的命令，调用时对内存占用进行检查
 *
 * a: admin command, like SAVE or SHUTDOWN.
 *    管理员使用的命令
 *
 * p: Pub/Sub related command.
 *    发送/订阅相关的命令
 *
 * f: force replication of this command, regarless of server.dirty.
 *    强制同步这个命令，无视 server.dirty
 *
 * s: command not allowed in scripts.
 *    不允许在脚本中使用的命令
 *
 * R: random command. Command is not deterministic, that is, the same command
 *    with the same arguments, with the same key space, may have different
 *    results. For instance SPOP and RANDOMKEY are two random commands.
 *    随机命令，对于同样数据集的同一个命令调用，得出的结果可能是不相同的。
 *
 * S: Sort command output array if called from script, so that the output
 *    is deterministic.
 *    如果命令在脚本中执行，那么对输出进行排序，从而让输出变得确定起来。
 *
 * l: Allow command while loading the database.
 *    允许在载入数据库时执行的命令
 *
 * t: Allow command while a slave has stale data but is not allowed to
 *    server this data. Normally no command is accepted in this condition
 *    but just a few.
 *    允许在附属节点包含过期数据时执行的命令
 *
 * M: Do not automatically propagate the command on MONITOR.
 *    不要自动将此命令发送到 MONITOR
 */
func (srv *Server) populateCommandTable() {
	connectionHandler := new(handlers.ConnectionHandler)
	stringHandler := new(handlers.StringHandler)
	keyHandler := new(handlers.KeyHandler)
	listHandler := new(handlers.ListHandler)
	hashHandler := new(handlers.HashHandler)
	setHandler := new(handlers.SetHandler)
	sortedSetHandler := new(handlers.SortedSetHandler)

	// connection command
	srv.commandTable[handlers.RedisConnectionCommandPing] = client.NewCommand(handlers.RedisConnectionCommandPing, 1, "r", connectionHandler.Ping)
	srv.commandTable[handlers.RedisConnectionCommandAuth] = client.NewCommand(handlers.RedisConnectionCommandAuth, 2, "rs", connectionHandler.Auth)
	srv.commandTable[handlers.RedisConnectionCommandSelect] = client.NewCommand(handlers.RedisConnectionCommandSelect, 2, "r", connectionHandler.CmdSelect)
	srv.commandTable[handlers.RedisConnectionCommandEcho] = client.NewCommand(handlers.RedisConnectionCommandEcho, 2, "r", connectionHandler.Echo)
	srv.commandTable[handlers.RedisConnectionCommandQuit] = client.NewCommand(handlers.RedisConnectionCommandQuit, 1, "r", connectionHandler.Quit)

	// key command
	srv.commandTable[handlers.RedisKeyCommandDel] = client.NewCommand(handlers.RedisKeyCommandDel, -2, "w", keyHandler.Del)
	srv.commandTable[handlers.RedisKeyCommandObject] = client.NewCommand(handlers.RedisKeyCommandObject, -2, "r", keyHandler.Object)
	srv.commandTable[handlers.RedisKeyCommandType] = client.NewCommand(handlers.RedisKeyCommandType, 2, "r", keyHandler.Type)
	srv.commandTable[handlers.RedisKeyCommandExists] = client.NewCommand(handlers.RedisKeyCommandExists, 2, "r", keyHandler.Exists)
	srv.commandTable[handlers.RedisKeyCommandDump] = client.NewCommand(handlers.RedisKeyCommandDump, 2, "ar", nil)
	srv.commandTable[handlers.RedisKeyCommandExpire] = client.NewCommand(handlers.RedisKeyCommandExpire, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandExpireAt] = client.NewCommand(handlers.RedisKeyCommandExpireAt, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandKeys] = client.NewCommand(handlers.RedisKeyCommandKeys, 2, "rS", nil)
	srv.commandTable[handlers.RedisKeyCommandMigrate] = client.NewCommand(handlers.RedisKeyCommandMigrate, -6, "aw", nil)
	srv.commandTable[handlers.RedisKeyCommandMove] = client.NewCommand(handlers.RedisKeyCommandMove, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandPersist] = client.NewCommand(handlers.RedisKeyCommandPersist, 2, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandPExpire] = client.NewCommand(handlers.RedisKeyCommandPExpire, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandPExpireAt] = client.NewCommand(handlers.RedisKeyCommandPExpireAt, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandPTTL] = client.NewCommand(handlers.RedisKeyCommandPTTL, 2, "r", nil)
	srv.commandTable[handlers.RedisKeyCommandRandomKey] = client.NewCommand(handlers.RedisKeyCommandRandomKey, 1, "rR", keyHandler.RandomKey)
	srv.commandTable[handlers.RedisKeyCommandRename] = client.NewCommand(handlers.RedisKeyCommandRename, 3, "w", keyHandler.Rename)
	srv.commandTable[handlers.RedisKeyCommandRenameNx] = client.NewCommand(handlers.RedisKeyCommandRenameNx, 3, "w", nil)
	srv.commandTable[handlers.RedisKeyCommandRestore] = client.NewCommand(handlers.RedisKeyCommandRestore, -4, "awm", nil)
	srv.commandTable[handlers.RedisKeyCommandSort] = client.NewCommand(handlers.RedisKeyCommandSort, -2, "wm", nil)
	srv.commandTable[handlers.RedisKeyCommandTTL] = client.NewCommand(handlers.RedisKeyCommandTTL, 2, "r", nil)
	srv.commandTable[handlers.RedisKeyCommandScan] = client.NewCommand(handlers.RedisKeyCommandScan, 2, "r", nil)

	// string command
	srv.commandTable[handlers.RedisStringCommandAppend] = client.NewCommand(handlers.RedisStringCommandAppend, 3, "wm", stringHandler.Append)
	srv.commandTable[handlers.RedisStringCommandSet] = client.NewCommand(handlers.RedisStringCommandSet, 3, "wm", stringHandler.Set)
	srv.commandTable[handlers.RedisStringCommandMSet] = client.NewCommand(handlers.RedisStringCommandMSet, -3, "wm", stringHandler.MSet)
	srv.commandTable[handlers.RedisStringCommandMSetNx] = client.NewCommand(handlers.RedisStringCommandMSetNx, -3, "wm", stringHandler.MSetNx)
	srv.commandTable[handlers.RedisStringCommandSetNx] = client.NewCommand(handlers.RedisStringCommandSetNx, 3, "wm", stringHandler.SetNx)
	srv.commandTable[handlers.RedisStringCommandGet] = client.NewCommand(handlers.RedisStringCommandGet, 2, "r", stringHandler.Get)
	srv.commandTable[handlers.RedisStringCommandMGet] = client.NewCommand(handlers.RedisStringCommandMGet, -2, "r", stringHandler.MGet)
	srv.commandTable[handlers.RedisStringCommandGetSet] = client.NewCommand(handlers.RedisStringCommandGetSet, 3, "wm", stringHandler.GetSet)
	srv.commandTable[handlers.RedisStringCommandIncr] = client.NewCommand(handlers.RedisStringCommandIncr, 2, "wm", stringHandler.Incr)
	srv.commandTable[handlers.RedisStringCommandIncrBy] = client.NewCommand(handlers.RedisStringCommandIncrBy, 3, "wm", stringHandler.IncrBy)
	srv.commandTable[handlers.RedisStringCommandIncrByFloat] = client.NewCommand(handlers.RedisStringCommandIncrByFloat, 3, "wm", stringHandler.IncrByFloat)
	srv.commandTable[handlers.RedisStringCommandDecr] = client.NewCommand(handlers.RedisStringCommandDecr, 2, "wm", stringHandler.Decr)
	srv.commandTable[handlers.RedisStringCommandDecrBy] = client.NewCommand(handlers.RedisStringCommandDecrBy, 3, "wm", stringHandler.DecrBy)
	srv.commandTable[handlers.RedisStringCommandStrLen] = client.NewCommand(handlers.RedisStringCommandStrLen, 2, "r", stringHandler.Strlen)

	// list command
	srv.commandTable[handlers.RedisListCommandLIndex] = client.NewCommand(handlers.RedisListCommandLIndex, 3, "r", listHandler.LIndex)
	srv.commandTable[handlers.RedisListCommandLInsert] = client.NewCommand(handlers.RedisListCommandLInsert, 5, "wm", listHandler.LInsert)
	srv.commandTable[handlers.RedisListCommandLLen] = client.NewCommand(handlers.RedisListCommandLLen, 2, "r", listHandler.LLen)
	srv.commandTable[handlers.RedisListCommandLPop] = client.NewCommand(handlers.RedisListCommandLPop, 2, "w", listHandler.LPop)
	srv.commandTable[handlers.RedisListCommandLPush] = client.NewCommand(handlers.RedisListCommandLPush, -3, "wm", listHandler.LPush)
	srv.commandTable[handlers.RedisListCommandLPushX] = client.NewCommand(handlers.RedisListCommandLPushX, 3, "wm", nil)
	srv.commandTable[handlers.RedisListCommandLRange] = client.NewCommand(handlers.RedisListCommandLRange, 4, "r", listHandler.LRange)
	srv.commandTable[handlers.RedisListCommandLRem] = client.NewCommand(handlers.RedisListCommandLRem, 4, "w", listHandler.LRem)
	srv.commandTable[handlers.RedisListCommandLSet] = client.NewCommand(handlers.RedisListCommandLSet, 4, "wm", listHandler.LSet)
	srv.commandTable[handlers.RedisListCommandLTrim] = client.NewCommand(handlers.RedisListCommandLTrim, 4, "w", listHandler.LTrim)
	srv.commandTable[handlers.RedisListCommandRPop] = client.NewCommand(handlers.RedisListCommandRPop, 2, "w", listHandler.RPop)
	srv.commandTable[handlers.RedisListCommandRPopLPush] = client.NewCommand(handlers.RedisListCommandRPopLPush, 3, "wm", nil)
	srv.commandTable[handlers.RedisListCommandRPush] = client.NewCommand(handlers.RedisListCommandRPush, -3, "wm", listHandler.RPush)
	srv.commandTable[handlers.RedisListCommandRpushX] = client.NewCommand(handlers.RedisListCommandRpushX, 3, "wm", nil)
	srv.commandTable[handlers.RedisListCommandLDebug] = client.NewCommand(handlers.RedisListCommandLDebug, 2, "r", listHandler.Debug)

	// hash command
	srv.commandTable[handlers.RedisHashCommandHDel] = client.NewCommand(handlers.RedisHashCommandHDel, -3, "w", hashHandler.HDel)
	srv.commandTable[handlers.RedisHashCommandHExists] = client.NewCommand(handlers.RedisHashCommandHExists, 3, "r", hashHandler.HExists)
	srv.commandTable[handlers.RedisHashCommandHGet] = client.NewCommand(handlers.RedisHashCommandHGet, 3, "r", hashHandler.HGet)
	srv.commandTable[handlers.RedisHashCommandHGetAll] = client.NewCommand(handlers.RedisHashCommandHGetAll, 2, "r", hashHandler.HGetAll)
	srv.commandTable[handlers.RedisHashCommandHIncrBy] = client.NewCommand(handlers.RedisHashCommandHIncrBy, 4, "wm", hashHandler.HIncrBy)
	srv.commandTable[handlers.RedisHashCommandHIncrByFloat] = client.NewCommand(handlers.RedisHashCommandHIncrByFloat, 4, "wm", hashHandler.HIncrByFloat)
	srv.commandTable[handlers.RedisHashCommandHKeys] = client.NewCommand(handlers.RedisHashCommandHKeys, 2, "rS", hashHandler.HKeys)
	srv.commandTable[handlers.RedisHashCommandHLen] = client.NewCommand(handlers.RedisHashCommandHLen, 2, "r", hashHandler.HLen)
	srv.commandTable[handlers.RedisHashCommandHMGet] = client.NewCommand(handlers.RedisHashCommandHMGet, -3, "r", hashHandler.HMGet)
	srv.commandTable[handlers.RedisHashCommandHMSet] = client.NewCommand(handlers.RedisHashCommandHMSet, -4, "wm", hashHandler.HMSet)
	srv.commandTable[handlers.RedisHashCommandHSet] = client.NewCommand(handlers.RedisHashCommandHSet, 4, "wm", hashHandler.HSet)
	srv.commandTable[handlers.RedisHashCommandHSetNX] = client.NewCommand(handlers.RedisHashCommandHSetNX, 4, "wm", hashHandler.HSetNX)
	srv.commandTable[handlers.RedisHashCommandHVals] = client.NewCommand(handlers.RedisHashCommandHVals, 2, "rS", hashHandler.HVals)
	srv.commandTable[handlers.RedisHashCommandHScan] = client.NewCommand(handlers.RedisHashCommandHScan, 3, "r", nil)
	srv.commandTable[handlers.RedisHashCommandHStrLen] = client.NewCommand(handlers.RedisHashCommandHStrLen, 3, "r", hashHandler.HStrLen)
	srv.commandTable[handlers.RedisHashCommandHDebug] = client.NewCommand(handlers.RedisHashCommandHDebug, 2, "r", hashHandler.HDebug)

	// set command
	srv.commandTable[handlers.RedisSetCommandSADD] = client.NewCommand(handlers.RedisSetCommandSADD, -3, "wm", setHandler.SAdd)
	srv.commandTable[handlers.RedisSetCommandSCARD] = client.NewCommand(handlers.RedisSetCommandSCARD, 2, "r", setHandler.SCard)
	srv.commandTable[handlers.RedisSetCommandSDIFF] = client.NewCommand(handlers.RedisSetCommandSDIFF, -2, "rS", setHandler.SDiff)
	srv.commandTable[handlers.RedisSetCommandSDIFFSTORE] = client.NewCommand(handlers.RedisSetCommandSDIFFSTORE, -3, "wm", setHandler.SDiffStore)
	srv.commandTable[handlers.RedisSetCommandSINTER] = client.NewCommand(handlers.RedisSetCommandSINTER, -2, "rS", setHandler.SInter)
	srv.commandTable[handlers.RedisSetCommandSINTERSTORE] = client.NewCommand(handlers.RedisSetCommandSINTERSTORE, -3, "wm", setHandler.SInterStore)
	srv.commandTable[handlers.RedisSetCommandSISMEMBER] = client.NewCommand(handlers.RedisSetCommandSISMEMBER, 3, "r", setHandler.SIsMember)
	srv.commandTable[handlers.RedisSetCommandSMEMBERS] = client.NewCommand(handlers.RedisSetCommandSMEMBERS, 2, "rS", setHandler.SMembers)
	srv.commandTable[handlers.RedisSetCommandSMOVE] = client.NewCommand(handlers.RedisSetCommandSMOVE, 4, "w", setHandler.SMove)
	srv.commandTable[handlers.RedisSetCommandSPOP] = client.NewCommand(handlers.RedisSetCommandSPOP, 2, "wRs", setHandler.SPop)
	srv.commandTable[handlers.RedisSetCommandSRANDMEMBER] = client.NewCommand(handlers.RedisSetCommandSRANDMEMBER, -2, "rR", setHandler.SRandMember)
	srv.commandTable[handlers.RedisSetCommandSREM] = client.NewCommand(handlers.RedisSetCommandSREM, -3, "r", setHandler.SRem)
	srv.commandTable[handlers.RedisSetCommandSUNION] = client.NewCommand(handlers.RedisSetCommandSUNION, -2, "rS", nil)
	srv.commandTable[handlers.RedisSetCommandSUNIONSTORE] = client.NewCommand(handlers.RedisSetCommandSUNIONSTORE, -3, "wm", nil)
	srv.commandTable[handlers.RedisSetCommandSSCAN] = client.NewCommand(handlers.RedisSetCommandSSCAN, 2, "rS", nil)

	// sorted set command
	srv.commandTable[handlers.RedisSortedSetCommandZAdd] = client.NewCommand(handlers.RedisSortedSetCommandZAdd, -4, "wm", sortedSetHandler.ZAdd)
	srv.commandTable[handlers.RedisSortedSetCommandZCard] = client.NewCommand(handlers.RedisSortedSetCommandZCard, 2, "r", sortedSetHandler.ZCard)
	srv.commandTable[handlers.RedisSortedSetCommandZCount] = client.NewCommand(handlers.RedisSortedSetCommandZCount, 4, "r", sortedSetHandler.ZCount)
	srv.commandTable[handlers.RedisSortedSetCommandZIncrBy] = client.NewCommand(handlers.RedisSortedSetCommandZIncrBy, 4, "wm", sortedSetHandler.ZIncrBy)
	srv.commandTable[handlers.RedisSortedSetCommandZRange] = client.NewCommand(handlers.RedisSortedSetCommandZRange, -4, "r", sortedSetHandler.ZRange)
	srv.commandTable[handlers.RedisSortedSetCommandZRangeByScore] = client.NewCommand(handlers.RedisSortedSetCommandZRangeByScore, -4, "r", sortedSetHandler.ZRangeByScore)
	srv.commandTable[handlers.RedisSortedSetCommandZRank] = client.NewCommand(handlers.RedisSortedSetCommandZRank, 3, "r", sortedSetHandler.ZRank)
	srv.commandTable[handlers.RedisSortedSetCommandZRem] = client.NewCommand(handlers.RedisSortedSetCommandZRem, -3, "w", sortedSetHandler.ZRem)
	srv.commandTable[handlers.RedisSortedSetCommandZRemRangeByRank] = client.NewCommand(handlers.RedisSortedSetCommandZRemRangeByRank, 4, "w", sortedSetHandler.ZRemRangeByRank)
	srv.commandTable[handlers.RedisSortedSetCommandZRemRangeByScore] = client.NewCommand(handlers.RedisSortedSetCommandZRemRangeByScore, 4, "w", sortedSetHandler.ZRemRangeByScore)
	srv.commandTable[handlers.RedisSortedSetCommandZRevRange] = client.NewCommand(handlers.RedisSortedSetCommandZRevRange, -4, "r", sortedSetHandler.ZRevRange)
	srv.commandTable[handlers.RedisSortedSetCommandZRevRangeByScore] = client.NewCommand(handlers.RedisSortedSetCommandZRevRangeByScore, -4, "r", sortedSetHandler.ZRevRangeByScore)
	srv.commandTable[handlers.RedisSortedSetCommandZRevRank] = client.NewCommand(handlers.RedisSortedSetCommandZRevRank, 3, "r", sortedSetHandler.ZRevRank)
	srv.commandTable[handlers.RedisSortedSetCommandZScore] = client.NewCommand(handlers.RedisSortedSetCommandZScore, 3, "r", sortedSetHandler.ZScore)
	srv.commandTable[handlers.RedisSortedSetCommandZUnionStore] = client.NewCommand(handlers.RedisSortedSetCommandZUnionStore, -4, "wm", nil)
	srv.commandTable[handlers.RedisSortedSetCommandZInterStore] = client.NewCommand(handlers.RedisSortedSetCommandZInterStore, -4, "wm", nil)
	srv.commandTable[handlers.RedisSortedSetCommandZScan] = client.NewCommand(handlers.RedisSortedSetCommandZScan, 2, "r", nil)

	// server command
	srv.commandTable[RedisServerCommandBGSRewriteAof] = client.NewCommand(RedisServerCommandBGSRewriteAof, 1, "ar", nil)
	srv.commandTable[RedisServerCommandBGSave] = client.NewCommand(RedisServerCommandBGSave, 1, "ar", srv.BgSave)
	srv.commandTable[RedisServerCommandClient] = client.NewCommand(RedisServerCommandClient, -2, "ar", nil)
	srv.commandTable[RedisServerCommandConfig] = client.NewCommand(RedisServerCommandConfig, -2, "ar", nil)
	srv.commandTable[RedisServerCommandDBSize] = client.NewCommand(RedisServerCommandDBSize, 1, "r", nil)
	srv.commandTable[RedisServerCommandDebug] = client.NewCommand(RedisServerCommandDebug, -2, "as", nil)
	srv.commandTable[RedisServerCommandFlushAll] = client.NewCommand(RedisServerCommandFlushAll, 1, "w", srv.FlushAll)
	srv.commandTable[RedisServerCommandFlushDB] = client.NewCommand(RedisServerCommandFlushDB, 1, "w", srv.FlushDB)
	srv.commandTable[RedisServerCommandInfo] = client.NewCommand(RedisServerCommandInfo, -1, "rlt", nil)
	srv.commandTable[RedisServerCommandLastSave] = client.NewCommand(RedisServerCommandLastSave, 1, "r", nil)
	srv.commandTable[RedisServerCommandMonitor] = client.NewCommand(RedisServerCommandMonitor, 1, "ars", nil)
	srv.commandTable[RedisServerCommandPSync] = client.NewCommand(RedisServerCommandPSync, 1, "ars", nil)
	srv.commandTable[RedisServerCommandShutDown] = client.NewCommand(RedisServerCommandShutDown, -1, "ar", nil)
	srv.commandTable[RedisServerCommandSave] = client.NewCommand(RedisServerCommandSave, 1, "ars", srv.Save)
	srv.commandTable[RedisServerCommandSlaveOf] = client.NewCommand(RedisServerCommandSlaveOf, 3, "ast", nil)
	srv.commandTable[RedisServerCommandSlowLog] = client.NewCommand(RedisServerCommandSlowLog, -2, "r", nil)
	srv.commandTable[RedisServerCommandSync] = client.NewCommand(RedisServerCommandSync, 1, "ars", nil)
	srv.commandTable[RedisServerCommandTime] = client.NewCommand(RedisServerCommandTime, 1, "rR", nil)
	srv.commandTable[RedisServerCommandAofDebug] = client.NewCommand(RedisServerCommandAofDebug, 1, "r", srv.AofDebug)
	srv.commandTable[RedisServerCommandAofFlush] = client.NewCommand(RedisServerCommandAofFlush, 1, "r", srv.AofFlush)

	srv.commandTable[RedisServerCommandCommand] = client.NewCommand(RedisServerCommandCommand, 1, "r", srv.Command)
	srv.commandTable[RedisServerCommandExit] = client.NewCommand(RedisServerCommandExit, 1, "r", srv.Exit)

	// pub sub command
	srv.commandTable[RedisPubSubCommandPSubscribe] = client.NewCommand(RedisPubSubCommandPSubscribe, -2, "rpslt", srv.PSubscribe)
	srv.commandTable[RedisPubSubCommandPublish] = client.NewCommand(RedisPubSubCommandPublish, 3, "pflt", srv.Publish)
	srv.commandTable[RedisPubSubCommandPUnsubscribe] = client.NewCommand(RedisPubSubCommandPUnsubscribe, -1, "rpslt", srv.PUnsubscribe)
	srv.commandTable[RedisPubSubCommandSubscribe] = client.NewCommand(RedisPubSubCommandSubscribe, -2, "rpslt", srv.Subscribe)
	srv.commandTable[RedisPubSubCommandUnsubscribe] = client.NewCommand(RedisPubSubCommandUnsubscribe, -1, "rpslt", srv.Unsubscribe)
	srv.commandTable[RedisPubSubCommandPubSub] = client.NewCommand(RedisPubSubCommandPubSub, -2, "r", nil)

	// debug command
	srv.commandTable[RedisDebugCommandRuntimeStat] = client.NewCommand(RedisDebugCommandRuntimeStat, 1, "r", srv.RuntimeStat)

	// 计算command flags
	for _, cmd := range srv.commandTable {
		for _, flag := range cmd.SFlags {
			switch flag {
			case 'w':
				cmd.Flags |= client.RedisCmdWrite
				break
			case 'r':
				cmd.Flags |= client.RedisCmdReadOnly
				break
			case 'm':
				cmd.Flags |= client.RedisCmdDenyOom
				break
			case 'a':
				cmd.Flags |= client.RedisCmdAdmin
				break
			case 'p':
				cmd.Flags |= client.RedisCmdPubSub
				break
			case 'f':
				cmd.Flags |= client.RedisCmdForceReplication
				break
			case 's':
				cmd.Flags |= client.RedisCmdNoScript
				break
			case 'R':
				cmd.Flags |= client.RedisCMdRandom
				break
			case 'S':
				cmd.Flags |= client.RedisCmdSortForScript
				break
			case 'l':
				cmd.Flags |= client.RedisCmdLoading
				break
			case 't':
				cmd.Flags |= client.RedisCmdStable
				break
			case 'M':
				cmd.Flags |= client.RedisCmdSkipMonitor
				break
			default:
				panic(fmt.Sprintf("Unsupported command flag:%v", cmd.SFlags))
			}
		}
	}
}
