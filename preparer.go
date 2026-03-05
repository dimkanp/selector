package selector

import (
	"database/sql"
	"fmt"
	"strings"
)

// Preparable is an interface for structure that represents database object to be selected
// with just table columns data or with additional fields
//
// Preparable don't know about selection count, one or more rows
type Preparable interface {
	SetSelectAlias(alias string)
	ScanFields(ctx *Context, selector Selector)
	GetAliasIterator() *AliasIterator
	UseAliasIterator(iterator *AliasIterator)
	ScanFieldNames() []string
	ScanFieldValues() []any
	SelectQuery() string

	Setup(ctx *Context, selector Selector)
}

type Base interface {
	ScanFieldNames() []string
	ScanFieldValues() []any
	TableName() string
}

func Prepare[P Preparable](ctx *Context, s Selector) P {
	var p P
	p.Setup(ctx, s)
	p.ScanFields(ctx, s)

	return p
}

type ScanReady interface {
	// ScanDestinations fieldNames parameter expected to be rows.Columns()
	// and returned slice should contain pointers to relevant places
	// in given order to scan values into
	ScanDestinations(fieldNames []string) []any
}

// ScanAll scans all returned records from rows into T structures
// by column names given from rows.Columns().
// *T should implement ScanReady interface to match scanned field with places for them
func ScanAll[T any, PT interface {
	*T        // PT must be pointer to T
	ScanReady // and implements ScanReady interface
	// ScanAll have two type parameters to divide structure from interface implementation.
	// Interface must be implemented by pointer receiver to allow rows.Scan method store
	// scanned values into the structure.
	// And value type T passed for efficient creation of new T objects without reflection
}](rows *sql.Rows) ([]*T, error) {
	var results []*T
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var x T
		p := PT(&x)
		err = rows.Scan(p.ScanDestinations(columns)...)
		if err != nil {
			return nil, err
		}

		results = append(results, &x)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

type Preparer struct {
	Alias               string
	ScanFieldNamesList  []string
	ScanFieldValuesList []any
	JoinsList           []string
	AliasIterator       *AliasIterator

	WhereClause string
	GroupClause string
	OrderClause string
	Limit       uint32
	Offset      uint32

	base Base
}

func (p *Preparer) Setup(ctx *Context, base Base) {
	p.base = base
	p.AliasIterator = ctx.Iterator
	if len(p.Alias) == 0 {
		p.Alias = p.AliasIterator.NextAlias()
	}
}

func (p *Preparer) SetSelectAlias(alias string) {
	p.Alias = alias
}

func (p *Preparer) ApplyAlias(field string) string {
	if len(p.Alias) == 0 {
		return field
	}

	return fmt.Sprintf("%s.%s", p.Alias, field)
}

func (p *Preparer) GetAliasIterator() *AliasIterator {
	if p.AliasIterator == nil {
		p.AliasIterator = NewAliasIterator()
	}

	return p.AliasIterator
}

func (p *Preparer) UseAliasIterator(iterator *AliasIterator) {
	p.AliasIterator = iterator
}

func (p *Preparer) ScanField(placeToScan any, selectPart string, joinPart string) {
	p.ScanFieldNamesList = append(p.ScanFieldNamesList, selectPart)
	p.ScanFieldValuesList = append(p.ScanFieldValuesList, placeToScan)
	p.JoinsList = append(p.JoinsList, joinPart)
}

func (p *Preparer) ScanFieldNames() (fields []string) {
	baseFields := p.base.ScanFieldNames()
	fields = make([]string, len(baseFields), len(baseFields)+len(p.ScanFieldNamesList))
	copy(fields, baseFields)

	if len(p.Alias) != 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("%s.%s", p.Alias, field)
		}
	}

	return append(fields, p.ScanFieldNamesList...)
}

func (p *Preparer) ScanFieldValues() (values []any) {
	if p.base != nil {
		values = p.base.ScanFieldValues()
	}

	return append(values, p.ScanFieldValuesList...)
}

func (p *Preparer) Where(where string) {
	tableName := p.base.TableName()

	if len(tableName) != 0 {
		p.WhereClause = strings.ReplaceAll(where, tableName+".", p.Alias+".")
	}
}

func (p *Preparer) OrderBy(order string) {
	tableName := p.base.TableName()

	if len(tableName) != 0 {
		p.OrderClause = strings.ReplaceAll(order, tableName+".", p.Alias+".")
	}
}

func (p *Preparer) SelectQuery() string {
	q := fmt.Sprintf(`SELECT %s FROM %s %s`, strings.Join(p.ScanFieldNames(), ", "), p.base.TableName(), p.Alias)

	if len(p.JoinsList) != 0 {
		q += " " + strings.Join(p.JoinsList, "\n\t")
	}

	if len(p.WhereClause) != 0 {
		q += fmt.Sprintf("\nWHERE %s", p.WhereClause)
	}

	if len(p.GroupClause) != 0 {
		q += fmt.Sprintf("\nGROUP BY %s", p.GroupClause)
	}

	if len(p.OrderClause) != 0 {
		q += fmt.Sprintf("\nORDER BY %s", p.OrderClause)
	}

	if p.Limit != 0 {
		q += fmt.Sprintf("\nLIMIT %d", p.Limit)
	}

	if p.Offset != 0 {
		q += fmt.Sprintf("\nOFFSET %d", p.Offset)
	}

	return q
}

func SelectFrom[T Base]() string {
	var base T
	return fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(base.ScanFieldNames(), ", "), base.TableName())
}

func SelectPart(preparable Preparable) string {
	return strings.Join(preparable.ScanFieldNames(), ", ")
}
