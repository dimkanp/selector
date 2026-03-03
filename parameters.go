package selector

import "fmt"

type ParamsCollector interface {
	// AddParameter receives parameter to pass with query and returns named (ordered) alias to use in the query
	// to match exact parameter place when simple parameters order in the query is not guaranteed
	AddParameter(p any) string
}

type PgxParamsCollector struct {
	pos    int
	params []any
}

func (p *PgxParamsCollector) AddParameter(param any) string {
	p.params = append(p.params, param)
	p.pos++
	return fmt.Sprintf("$%d", p.pos)
}
