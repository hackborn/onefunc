package msg

type HandlerFunc[T any] func(string, T)

type visitFunc func(any)

type visitSubscriptionsFunc func(*subscriptions)
