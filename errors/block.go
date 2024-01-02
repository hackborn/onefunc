package errors

type Block interface {
	Add(error)
}

type FirstBlock struct {
	Err error
}

func (c *FirstBlock) Add(e error) {
	if c.Err == nil {
		c.Err = e
	}
}

type NullBlock struct {
}

func (c *NullBlock) Add(e error) {
}
