package errors

// Block tracks one or more errors.
type Block interface {
	AddError(error)
	HasError() bool
}

// FirstBlock is a Block that stores the first error it receives.
type FirstBlock struct {
	Err error
}

func (c *FirstBlock) AddError(e error) {
	if c.Err == nil {
		c.Err = e
	}
}

func (c *FirstBlock) HasError() bool {
	return c.Err != nil
}

// NullBlock is an empty Block that doesn't store errors.
type NullBlock struct {
}

func (c *NullBlock) AddError(e error) {
}

func (c *NullBlock) HasError() bool {
	return false
}
