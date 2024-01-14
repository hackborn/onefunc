package errors

type Block interface {
	AddError(error)
}

type FirstBlock struct {
	Err error
}

func (c *FirstBlock) AddError(e error) {
	if c.Err == nil {
		c.Err = e
	}
}

type NullBlock struct {
}

func (c *NullBlock) AddError(e error) {
}
