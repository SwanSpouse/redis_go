package raw_type

type RawInt struct {
	RedisObject
}

func (ri *RawInt) Append(val string) int {
	return 0
}

func (ri *RawInt) Incr() int {
	return 0
}

func (ri *RawInt) Decr() int {
	return 0
}

func (ri *RawInt) IncrBy(val int) int {
	return 0
}

func (ri *RawInt) DecrBy(val int) int {
	return 0
}

func (ri *RawInt) Strlen() int {
	return 0
}
