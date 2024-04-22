package msg

type HandlerFunc[T any] func(string, T)

type initFunc[T any] func(pattern string, data *T)

type editFunc[T any] func(int64, *T)

type visitFunc[T any] func(topic string, data *T)

type MatchFunc func(pattern, topic string) bool
