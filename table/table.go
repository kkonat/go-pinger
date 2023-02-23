package table

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Col struct {
	Name    string
	Width   int
	Sep     string
	R       bool
	visible bool
}

type Data interface {
	GenCells() *[]string
}

type Table struct {
	Rows        int
	Width       int
	header      string
	Cols        []Col
	data        []Data
	displItemNo bool
}

func New() *Table {
	t := &Table{}
	t.Cols = make([]Col, 0, 10)
	return t
}

// func (t *Table) Init(row any)
// defines table columns by interpreting table row in struct "td:" tags
func (t *Table) Init(rows []Data, displItemNo bool) {
	var key, val string

	t.data = rows
	t.displItemNo = displItemNo
	if len(rows) == 0 {
		panic("Table is empty")
	}
	typ := reflect.TypeOf(rows[0])

	for tagIndx := 0; tagIndx < typ.NumField(); tagIndx++ {
		fldType := typ.Field(tagIndx)
		tags := strings.Split(fldType.Tag.Get("td"), ",")
		col := Col{
			visible: false,
			Name:    fldType.Name,
		}

		for _, currTag := range tags {
			currTag = strings.TrimSpace(currTag)
			kv := strings.Split(currTag, "=")

			if len(kv) == 2 {
				key, val = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			} else if len(kv) == 1 {
				key, val = strings.TrimSpace(kv[0]), ""
			} else {
				panic(fmt.Sprintf("table td: tag - key=value pair invalid: '%s'", currTag))
			}

			switch key {

			case "name":
				if len(val) < 3 && val[0] != '\'' && val[len(val)-1] != '\'' {
					panic("table td: tag - name must be at least 1 char within ' '")
				}
				if col.Width != 0 && len(val)-2 > col.Width {
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)",
						val, len(val)-2, col.Width))
				}
				col.Name = val[1 : len(val)-1]
				col.visible = true

			case "w":
				w, err := strconv.Atoi(val)
				if err != nil {
					panic(fmt.Sprintf("table td: tag - w is not an int: '%s'", currTag))
				}
				if len(col.Name) != 0 && len(col.Name) > w {
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)",
						col.Name, len(col.Name), w))
				}
				col.Width = w
				col.visible = true
			case "sep":
				if len(val) < 3 {
					panic("table: td: sep must be at least 1 char within ' '")
				}
				col.Sep = val[1 : len(val)-1]

			case "R":
				fallthrough
			case "right":
				col.R = true

			}
		}
		if col.Width == 0 && col.visible {
			panic(fmt.Sprintf("table: td: Invalid column definition: %s", tags))
		}
		t.Cols = append(t.Cols, col)
		// fmt.Printf("%#v\n", col)
	}
	t.Width = t.buildHeader()
}

func (t *Table) buildHeader() int {
	t.header = "|"
	if t.displItemNo {
		t.header += "   #"
	}
	for _, c := range t.Cols {
		if c.visible {
			if c.Sep != "" {
				t.header += " " + c.Sep + " "
			} else {
				t.header += " "
			}
			if !c.R {
				t.header += c.Name
			}
			t.header += strings.Repeat(" ", c.Width-len(c.Name))
			if c.R {
				t.header += c.Name
			}

		}
	}
	t.header += "|"
	return len(t.header)
}
func (t Table) DisplaySep() {
	fmt.Println(strings.Repeat("-", len(t.header)))
}
func (t Table) DisplayHeader() {
	fmt.Println(t.header)
}

func (t Table) genRow(idx int) string {
	cells := t.data[idx].GenCells()

	row := "| "
	if t.displItemNo {
		row += fmt.Sprintf("%3d", idx)
	}
	col := 0
	for i := 0; i < len(*cells); i++ {
		c := t.Cols[col]
		if c.visible {
			s := (*cells)[i]

			if c.Sep != "" {
				row += " " + c.Sep + " "
			} else {
				row += " "
			}

			pad := 0
			l := len((*cells)[i])
			if l > c.Width {
				s = (*cells)[i][:c.Width-3] + "..."
			}
			if l < c.Width {
				pad = c.Width - l
			}
			if c.R {
				row += strings.Repeat(" ", pad)
			}
			row += s
			if !c.R {
				row += strings.Repeat(" ", pad)
			}
		}

		col++
	}
	row += "|"
	return row
}

func (t Table) Print() {
	t.DisplaySep()
	t.DisplayHeader()
	t.DisplaySep()
	for i := 0; i < len(t.data); i++ {
		fmt.Println(t.genRow(i))
	}
	t.DisplaySep()
}
