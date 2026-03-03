package selector

import (
	"fmt"
	"sync/atomic"
)

func NewAliasIterator() *AliasIterator {
	return &AliasIterator{}
}

type AliasIterator struct {
	a atomic.Int32
}

func (p *AliasIterator) NextAlias() string {
	return fmt.Sprintf("s_p%d", p.a.Add(1))
}
