package table

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type col struct {
	name       string
	width      int
	separator  string
	rightAlign bool
	visible    bool

	invalidStr string
}
type validatable interface {
	IsValid() bool
}

type Table[Data validatable] struct {
	width      int
	header     string
	cols       []col
	rows       []Data
	displayNos bool
}

func New[Data validatable](rows []Data, displItemNo bool) *Table[Data] {
	var (
		a Data
		t = &Table[Data]{
			displayNos: displItemNo,
			rows:       rows,
			cols:       make([]col, reflect.TypeOf(a).NumField()),
		}
		key, val string
		typ      = reflect.TypeOf(a)
	)

	for tagIndx := 0; tagIndx < typ.NumField(); tagIndx++ {
		fldType := typ.Field(tagIndx)
		tags := strings.Split(fldType.Tag.Get("td"), ",")
		col := col{
			visible: false,
			name:    fldType.Name,
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

			valLen := len(val)
			switch key {

			case "name":
				if valLen < 3 && val[0] != '\'' && val[valLen-1] != '\'' {
					panic("table td: tag - name must be at least 1 char within ' '")
				}
				if col.width != 0 && valLen-2 > col.width {
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)",
						val, valLen-2, col.width))
				}
				col.name = val[1 : valLen-1]
				col.visible = true

			case "w":
				w, err := strconv.Atoi(val)
				if err != nil {
					panic(fmt.Sprintf("table td: tag - w is not an int: '%s'", currTag))
				}
				if len(col.name) != 0 && len(col.name) > w {
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)",
						col.name, len(col.name), w))
				}
				col.width = w
				col.visible = true

			case "sep":
				if valLen < 3 {
					panic("table: td: sep must be at least 1 char within ' '")
				}
				col.separator = val[1 : valLen-1]

			case "R":
				fallthrough
			case "right":
				col.rightAlign = true
			case "inv":
				if valLen < 3 && val[0] != '\'' && val[valLen-1] != '\'' {
					panic("table td: inv - name must be at least 1 char within ' '")
				}
				col.invalidStr = val[1 : valLen-1]
			}
		}
		if col.width == 0 && col.visible {
			panic(fmt.Sprintf("table: td: Invalid column definition: %s", tags))
		}
		t.cols[tagIndx] = col
		// fmt.Printf("%#v\n", col)
	}
	t.width = t.buildHeader()
	return t
}

func (t *Table[Data]) buildHeader() int {
	t.header = "|"
	if t.displayNos {
		t.header += "   #"
	}
	for _, c := range t.cols {
		if c.visible {
			if c.separator != "" {
				t.header += " " + c.separator + " "
			} else {
				t.header += " "
			}
			if !c.rightAlign {
				t.header += c.name
			}
			t.header += strings.Repeat(" ", c.width-len(c.name))
			if c.rightAlign {
				t.header += c.name
			}

		}
	}
	t.header += "|"
	return len(t.header)
}

func (t *Table[Data]) MaxWidth() int { return t.width }

func (t Table[Data]) DisplaySep() {
	fmt.Println(strings.Repeat("-", len(t.header)))
}

func (t Table[Data]) DisplayHeader() {
	fmt.Println(t.header)
}
func (t *Table[Data]) varToString(row, col int) string {
	var (
		r           = reflect.ValueOf(t.rows[row])
		format, str = "", ""

		vType = r.Field(col).Type().Name()
		value = r.Field(col)
		valid = r.FieldByName("valid")
		w     = t.cols[col].width
	)
	// fmt.Println(valid, ":", vType, "=", value)
	if !valid.Bool() {
		return t.cols[col].invalidStr
	}

	switch vType {
	case "uint64":
		fallthrough
	case "uint32":
		fallthrough
	case "uint16":
		fallthrough
	case "uint8":
		format = fmt.Sprintf("%%%dd", w)
		str = fmt.Sprintf(format, int64(value.Uint()))
	case "int64":
		fallthrough
	case "int32":
		fallthrough
	case "int16":
		fallthrough
	case "int8":
		format = fmt.Sprintf("%%%dd", w)
		str = fmt.Sprintf(format, int64(value.Int()))
	case "float64":
		fallthrough
	case "float32":
		format = fmt.Sprintf("%%%d.2f", w)
		str = fmt.Sprintf(format, float32(value.Float()))
	case "string":
		str = value.String()
	}
	return str
}

func (t Table[Data]) generateRow(row int) string {
	var a Data
	typ := reflect.TypeOf(a)

	line := "| "

	if t.displayNos { // add line numbers
		line += fmt.Sprintf("%3d", row)
	}

	for col := 0; col < typ.NumField(); col++ {
		var (
			c      = t.cols[col]
			s, str string
		)

		if c.visible {

			str = t.varToString(row, col)

			if c.separator != "" {
				line += " " + c.separator + " "
			} else {
				line += " "
			}

			pad := 0
			l := len(str)
			if l > c.width {
				s = str[:c.width-3] + "..."
			} else {
				s = str
			}
			if l < c.width {
				pad = c.width - l
			}
			if c.rightAlign {
				line += strings.Repeat(" ", pad)
			}
			line += s
			if !c.rightAlign {
				line += strings.Repeat(" ", pad)
			}
		}
	}
	line += "|"
	return line
}

func (t *Table[Data]) Print() {
	t.DisplaySep()
	t.DisplayHeader()
	t.DisplaySep()
	for i := 0; i < len(t.rows); i++ {
		fmt.Println(t.generateRow(i))
	}
	t.DisplaySep()
}
