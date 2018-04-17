package raw_type

type Sds struct {
	RedisObject
}

func (sds *Sds) Append(val string) int {
	return 0
}

func (sds *Sds) Incr() int {
	return 0
}

func (sds *Sds) Decr() int {
	return 0
}

func (sds *Sds) IncrBy(val int) int {
	return 0
}

func (sds *Sds) DecrBy(val int) int {
	return 0
}

func (sds *Sds) Strlen() int {
	return 0
}
