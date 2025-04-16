package errorcode

const (
	MysqlScanErr                    = middleware*10000 + 1
	MqttConnectTokenError           = middleware*10000 + 2
	RedisGetBytes                   = middleware*10000 + 3
	KafkaProduceTopicPartitionError = middleware*10000 + 4
	KafkaProduce                    = middleware*10000 + 5
	KafkaConsumerSubscribe          = middleware*10000 + 7
	KafkaConsumerReadMessage        = middleware*10000 + 8
	RedisGetInt64                   = middleware*10000 + 9
	MysqlQueryErr                   = middleware*10000 + 10
	MysqlRowsCloseErr               = middleware*10000 + 11
	MysqlRowsScanErr                = middleware*10000 + 12
	MysqlExecErr                    = middleware*10000 + 13
	MysqlLastInsertIdErr            = middleware*10000 + 14
	MysqlTxErr                      = middleware*10000 + 15
	KafkaCommit                     = middleware*10000 + 16
	MysqlSharding                   = middleware*10000 + 17
	MysqlCommit                     = middleware*10000 + 18
	MysqlRollback                   = middleware*10000 + 19
	MysqlShardingTimeRangeUnknown   = middleware*10000 + 20
)

func MiddlewareCode() {
	ErrorCode[MysqlScanErr] = "mysql scan error"
	ErrorCode[MqttConnectTokenError] = "mqtt connect token error"
	ErrorCode[RedisGetBytes] = "redis get bytes error"
	ErrorCode[KafkaProduceTopicPartitionError] = "kafka produce topic partition error"
	ErrorCode[KafkaProduce] = "kafka produce error"
	ErrorCode[KafkaConsumerSubscribe] = "kafka consumer subscribe err"
	ErrorCode[KafkaConsumerReadMessage] = "kafka consumer read message err"
	ErrorCode[RedisGetInt64] = "redis get int64 error"
	ErrorCode[MysqlQueryErr] = "mysql query err"
	ErrorCode[MysqlRowsCloseErr] = "mysql rows close err"
	ErrorCode[MysqlRowsScanErr] = "mysql rows scan err"
	ErrorCode[MysqlExecErr] = "mysql exec err"
	ErrorCode[MysqlLastInsertIdErr] = "mysql exec last insert id err"
	ErrorCode[MysqlTxErr] = "mysql tx err"
	ErrorCode[KafkaCommit] = "kafka commit err"
	ErrorCode[MysqlSharding] = "mysql sharding err"
	ErrorCode[MysqlCommit] = "mysql commit err"
	ErrorCode[MysqlRollback] = "mysql rollback err"
	ErrorCode[MysqlShardingTimeRangeUnknown] = "mysql sharding time range unknown"
}
