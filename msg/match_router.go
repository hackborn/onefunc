package msg

type MatchRouter[T any] struct {
	Init     initFunc[T]
	Match    MatchFunc
	patterns map[string]*matchEntry[T]
}

func (r *MatchRouter[T]) Edit(pattern string, fn editFunc[T]) {
	r.validate()
	gp, ok := r.patterns[pattern]
	if !ok {
		gp = &matchEntry[T]{}
		if r.Init != nil {
			r.Init(pattern, &gp.data)
		}
		r.patterns[pattern] = gp
	}
	fn(0, &gp.data)
}

// Visit any glob patterns I contain that match the topic.
func (r *MatchRouter[T]) Visit(topic string, fn visitFunc[T]) {
	for k, v := range r.patterns {
		if m := r.Match(k, topic); m {
			fn(k, &v.data)
		}
	}
}

func (r *MatchRouter[T]) validate() {
	if r.patterns == nil {
		r.patterns = make(map[string]*matchEntry[T])
	}
}

type matchEntry[T any] struct {
	data T
}
