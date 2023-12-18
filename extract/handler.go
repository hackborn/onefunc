package extract

type Handler interface {
	// Handle the pair in some fashion. The basic system
	// is designed to work without errors; if a client needs
	// any error handling, it should track that in an internal state.
	Handle(name string, value any)
}

type Chainer interface {
	Handler
	GetNext() Handler
	SetNext(Handler)
}

type Chain struct {
	Next Handler
}

func (c *Chain) GetNext() Handler {
	return c.Next
}

func (c *Chain) SetNext(h Handler) {
	c.Next = h
}

func WithMap(h Handler, m map[string]string) Handler {
	return &remapHandler{Chain: Chain{Next: h}, m: m, filter: false}
}

func WithFilterMap(h Handler, m map[string]string) Handler {
	return &remapHandler{Chain: Chain{Next: h}, m: m, filter: true}
}

func Map(m map[string]string) Chainer {
	return &remapHandler{m: m, filter: false}
}

func FilterMap(m map[string]string) Chainer {
	return &remapHandler{m: m, filter: true}
}

func NewChain(chainers ...Chainer) Handler {
	var first Chainer = nil
	last := first
	for _, c := range chainers {
		if first == nil {
			first = c
			last = c
		} else {
			last.SetNext(c)
			last = c
		}
	}
	return first
}

func tail(h Handler) Chainer {
	var last Chainer = nil
	for c, ok := h.(Chainer); ok; c, ok = h.(Chainer) {
		last = c
		h = c.GetNext()
	}
	return last
}

type remapHandler struct {
	Chain
	m      map[string]string
	filter bool
}

func (h *remapHandler) Handle(name string, value any) {
	if h.Next == nil {
		return
	}
	if n, ok := h.m[name]; ok {
		h.Next.Handle(n, value)
	} else if !h.filter {
		h.Next.Handle(name, value)
	}
}

type sliceHandler struct {
	result     []any
	assignment string
	connector  string
}

func (h *sliceHandler) Handle(name string, value any) {
	if h.connector != "" && len(h.result) > 0 {
		h.result = append(h.result, h.connector)
	}
	h.result = append(h.result, name)
	if h.assignment != "" {
		h.result = append(h.result, h.assignment)
	}
	h.result = append(h.result, value)
}
