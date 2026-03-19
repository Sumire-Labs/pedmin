package player

type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

func (l LoopMode) String() string {
	switch l {
	case LoopTrack:
		return "🔂 Track"
	case LoopQueue:
		return "🔁 Queue"
	default:
		return "Loop Off"
	}
}

func (l LoopMode) Next() LoopMode {
	return (l + 1) % 3
}
