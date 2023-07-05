package tanker

type Bucket interface {
	Fire() bool
	Update() uint32
	Rest() uint32
}

type Tanker interface {
	Check(id string) bool
	Flush()
	Remove(id string)
}
