package common

type Num interface {
	int | int32 | int8 | int16 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}
