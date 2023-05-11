package test

type UseCase[IN any, Out any] struct {
	Input        IN
	ExceptOutput Out
}
