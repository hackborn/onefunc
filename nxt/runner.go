package nxt

func newPipeline() Pipeline {
	return Pipeline{}
}

type Pipeline struct {
	heads []Handler
}

func (p Pipeline) Run(event any) {
	r := newRunner(p)
	r.run(event)
}

func newRunner(p Pipeline) *runner {
	id := runnerCounter.Add(1)
	r := &runner{id: id, heads: p.heads}
	return r
}

type runner struct {
	id    int64
	heads []Handler
}

func (r *runner) run(event any) {
	args := HandlerArgs{runnerId: r.id}
	for _, h := range r.heads {
		h(args, event)
	}
}
