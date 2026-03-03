package selector

type Context struct {
	Iterator        *AliasIterator
	ParamsCollector ParamsCollector
}

type ContextOption func(*Context)

// NewContext default behavior targeted to use with Postgres database
func NewContext(options ...ContextOption) *Context {
	c := &Context{
		Iterator:        NewAliasIterator(),
		ParamsCollector: &PgxParamsCollector{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithIterator(iterator *AliasIterator) ContextOption {
	return func(c *Context) {
		c.Iterator = iterator
	}
}

func WithParamsCollector(collector ParamsCollector) ContextOption {
	return func(c *Context) {
		c.ParamsCollector = collector
	}
}
