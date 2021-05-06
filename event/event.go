package event

type GetEvent struct {
	Pid           uint32
	Pad           uint32
	Start_time_ns uint64
	Duration      uint64
	Klen          int32
	Key           [128]byte
}

type SetEvent struct {
	Pid           uint32
	Pad           uint32
	Start_time_ns uint64
	Duration      uint64
	Klen          int32
	Key           [128]byte
}
