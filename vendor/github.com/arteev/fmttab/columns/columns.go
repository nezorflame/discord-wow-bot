package columns

import (
	"errors"
)

const (
	//WidthAuto auto sizing of width column
	WidthAuto = 0
	//AlignLeft align text along the left edge
	AlignLeft = Align(false)
	//AlignRight align text along the right edge
	AlignRight = Align(true)
)

//Errors
var (
	ErrorAlreadyExists = errors.New("Column already exists")
)

//A Align text alignment in column of the table
type Align bool

//A Column type of table columns
type Column struct {
	MaxLen  int
	Caption string
	Name    string
	Width   int
	Aling   Align
	Visible bool
}

//A Columns array of the columns
type Columns []*Column

//IsAutoSize returns auto sizing of width column
func (t Column) IsAutoSize() bool {
	return t.Width == WidthAuto
}

//Len returns count columns
func (c *Columns) Len() int {
	if c == nil {
		return 0
	}
	return len(*c)
}

//FindByName returns columns by name if exists or nil
func (c *Columns) FindByName(name string) *Column {
	for i := range *c {
		if (*c)[i].Name == name {
			return (*c)[i]
		}
	}
	return nil
}

//NewColumn append new column in list with check by name of column
func (c *Columns) NewColumn(name, caption string, width int, aling Align) (*Column, error) {
	if c.FindByName(name) != nil {
		return nil, ErrorAlreadyExists
	}
	column := &Column{
		Name:    name,
		Caption: caption,
		Width:   width,
		Aling:   aling,
		Visible: true,
	}
	*c = append(*c, column)
	return column, nil
}

//Add append column with check exists
func (c *Columns) Add(col *Column) error {
	for i := range *c {
		if (*c)[i] == col {
			return ErrorAlreadyExists
		}
	}
	*c = append(*c, col)
	return nil
}

//ColumnsVisible returns count visible columns
func (c *Columns) ColumnsVisible() (res Columns) {
	for i, col := range *c {
		if col.Visible {
			res = append(res, (*c)[i])
		}
	}
	return
}
