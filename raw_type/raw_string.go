package raw_type

type RawString struct {
	RedisObject
}

func (rs *RawString) Append(val string) int {
	return 0
}

func (rs *RawString) Incr() int {
	return 0
}

func (rs *RawString) Decr() int {
	return 0
}

func (rs *RawString) IncrBy(val int) int {
	return 0
}

func (rs *RawString) DecrBy(val int) int {
	return 0
}

func (rs *RawString) Strlen() int {
	return 0
}
