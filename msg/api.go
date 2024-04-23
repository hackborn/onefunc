package msg

// Subscribe to the topic with the given function. Answer
// the subscription. Use the subscription to unsubscribe.
// The last message publshed to the topic will be immediately
// sent to the function.
// Topic subscriptions use MQTT rules:
// Pattern is a hierarchy with / separators.
// Wildcard "+" matches a single level in the hierarchy.
// Wildcard "#" matches all remaining levels in the hierarchy.
func Sub[T any](r *Router, pattern string, fn HandlerFunc[T]) Subscription {
	subs := r.sub(pattern, fn)
	retainFn := func(topic string, last any) {
		if lastT, ok := last.(T); ok {
			fn(topic, lastT)
		}
	}
	r.visitRetained(pattern, retainFn)
	return subs
}

func Pub[T any](r *Router, topic string, value T) {
	if r == nil {
		return
	}
	r.readVisit(topic, func(pattern string, subs *routerSubscriptions) {
		for _, _h := range subs.subs {
			if h, ok := _h.(HandlerFunc[T]); ok {
				h(topic, value)
			}
		}
	})
	r.retain(topic, value)
}
