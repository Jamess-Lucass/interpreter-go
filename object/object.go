package object

import "fmt"

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

var _ Object = (*Integer)(nil)

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

var _ Object = (*Boolean)(nil)

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Null struct{}

var _ Object = (*Null)(nil)

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

var _ Object = (*ReturnValue)(nil)

func (n *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func (n *ReturnValue) Inspect() string {
	return n.Value.Inspect()
}
