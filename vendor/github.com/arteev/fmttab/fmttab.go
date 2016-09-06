package fmttab

import (
	"strconv"
	"unicode/utf8"

	"github.com/arteev/fmttab/columns"
)

//A Border of table
type Border int

//A BorderKind type of element on the border of the table
type BorderKind int

const (
	//WidthAuto auto sizing of width column
	WidthAuto = columns.WidthAuto
	//BorderNone table without borders
	BorderNone = Border(0)
	//BorderThin table with a thin border
	BorderThin = Border(1)
	//BorderDouble table with a double border
	BorderDouble = Border(2)

	//BorderSimple table with a simple border
	BorderSimple = Border(3)
	//AlignLeft align text along the left edge
	AlignLeft = columns.AlignLeft
	//AlignRight align text along the right edge
	AlignRight = columns.AlignRight
)

//The concrete type of the object on the border of the table
const (
	BKLeftTop BorderKind = iota
	BKRighttop
	BKRightBottom
	BKLeftBottom
	BKLeftToRight
	BKRightToLeft
	BKTopToBottom
	BKBottomToTop
	BKBottomCross
	BKHorizontal
	BKVertical
	BKHorizontalBorder
	BKVerticalBorder
)

//Trimend - end of line after trimming
var Trimend = ".."
var trimlen = utf8.RuneCountInString(Trimend)

//Borders predefined border types
var Borders = map[Border]map[BorderKind]string{
	BorderNone: map[BorderKind]string{
		BKVertical: " ",
	},
	BorderSimple: map[BorderKind]string{
		BKBottomCross: "+",
		BKHorizontal:  "-",
		BKVertical:    "|",
	},
	BorderThin: map[BorderKind]string{
		BKLeftTop:          "\u250c",
		BKRighttop:         "\u2510",
		BKRightBottom:      "\u2518",
		BKLeftBottom:       "\u2514",
		BKLeftToRight:      "\u251c",
		BKRightToLeft:      "\u2524",
		BKTopToBottom:      "\u252c",
		BKBottomToTop:      "\u2534",
		BKBottomCross:      "\u253c",
		BKHorizontal:       "\u2500",
		BKVertical:         "\u2502",
		BKHorizontalBorder: "\u2500",
		BKVerticalBorder:   "\u2502",
	},
	BorderDouble: map[BorderKind]string{
		BKLeftTop:          "\u2554",
		BKRighttop:         "\u2557",
		BKRightBottom:      "\u255d",
		BKLeftBottom:       "\u255a",
		BKLeftToRight:      "\u255f",
		BKRightToLeft:      "\u2562",
		BKTopToBottom:      "\u2564",
		BKBottomToTop:      "\u2567",
		BKBottomCross:      "\u253c",
		BKHorizontal:       "\u2500",
		BKVertical:         "\u2502",
		BKHorizontalBorder: "\u2550",
		BKVerticalBorder:   "\u2551",
	},
}

//A DataGetter functional type for table data
type DataGetter func() (bool, map[string]interface{})

//A Table is the repository for the columns, the data that are used for printing the table
type Table struct {
	dataget         DataGetter
	border          Border
	caption         string
	autoSize        int
	CloseEachColumn bool
	Columns         columns.Columns
	Data            []map[string]interface{}
	VisibleHeader   bool
	masks           map[*columns.Column]string
	columnsvisible  columns.Columns
}

// A trimEnds supplements the text with special characters by limiting the length of the text column width
func trimEnds(val string, max int) string {
	if utf8.RuneCountInString(val) <= max {
		return val
	}
	if trimlen < max {
		return val[:max-trimlen] + Trimend
		//return string([]rune(val)[:(max-trimlen)]) + end
	}
	return Trimend[:max]
}

//GetMaskFormat returns a pattern string for formatting text in table column alignment
func (t *Table) GetMaskFormat(c *columns.Column) string {
	if c.Aling == AlignLeft {
		return "%-" + strconv.Itoa(t.getWidth(c)) + "v"
	}
	return "%" + strconv.Itoa(t.getWidth(c)) + "v"
}

//must be calculated before call
func (t *Table) getWidth(c *columns.Column) int {
	if c.IsAutoSize() || t.autoSize > 0 {
		return c.MaxLen
	}
	return c.Width

}

//AddColumn adds a column to the table
func (t *Table) AddColumn(name string, width int, aling columns.Align) *Table {
	_, err := t.Columns.NewColumn(name, name, width, aling)
	if err != nil {
		//fix it
		panic(err)
	}
	return t
}

//AppendData adds the data to the table
func (t *Table) AppendData(rec map[string]interface{}) *Table {
	t.Data = append(t.Data, rec)
	return t
}

//ClearData removes data from a table
func (t *Table) ClearData() *Table {
	t.Data = nil
	return t
}

//AutoSize fit columns
func (t *Table) AutoSize(enabled bool, destWidth int) {
	if enabled {
		t.autoSize = destWidth
	} else {
		t.autoSize = 0
	}
}

//CountData the amount of data in the table
func (t *Table) CountData() int {
	return len(t.Data)
}

//SetBorder - set  type of border table
func (t *Table) SetBorder(b Border) {
	t.border = b
}

//GetBorder - get current border
func (t Table) GetBorder() Border {
	return t.border
}

//New creates a Table object. DataGetter can be nil
func New(caption string, border Border, datagetter DataGetter) *Table {
	return &Table{
		caption:       caption,
		border:        border,
		dataget:       datagetter,
		VisibleHeader: true,
	}
}
