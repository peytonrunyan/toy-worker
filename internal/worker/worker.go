package worker

type Consumer interface {
	Consume()
}

type Producer interface {
	Produce()
}
