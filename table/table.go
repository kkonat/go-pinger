package table

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type validable interface {
	IsValid() bool
}

type column struct {
	name       string
	width      int
	separator  string
	rightAlign bool
	visible    bool
	invalidStr string
}

type Table[Data validable] struct {
	width    int
	header   string
	cols     []column
	rows     []Data
	displNos bool // display table row numbers?
}

// func New[Data validable](rows []Data, displItemNo bool) *Table[Data]
// creates a new Table with Data and displItemNo
func New[Data validable](rows []Data, displayItemNumbers bool) *Table[Data] {
	var (
		dataType Data
		typ      = reflect.TypeOf(dataType)
		newTable = &Table[Data]{
			displNos: displayItemNumbers,
			rows:     rows,
			cols:     make([]column, typ.NumField()),
		}
	)

	for colN := 0; colN < typ.NumField(); colN++ {
		var (
			colType = typ.Field(colN)
			newCol  = column{visible: false, name: colType.Name}
		)
		tags := colType.Tag.Get("td")

		for _, currTag := range strings.Split(tags, ",") {
			key, val := getKVPair(currTag)
			l := len(val)

			switch strings.ToLower(key) {

			case "name":
				if isntEnclosed(val, '\'') {
					panic(fmt.Sprintf("table td: tag - name (%s) must be at least 1 char within ' '", val))
				}

				if newCol.width != 0 && l-2 > newCol.width { // check if col name isn't wider than column
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)",
						val, l-2, newCol.width))
				}
				newCol.name = strings.Trim(val, "'")
				newCol.visible = true

			case "width", "w":
				w, err := strconv.Atoi(val)
				if err != nil {
					panic(fmt.Sprintf("table td: tag - w is not an int: '%s'", currTag))
				}
				if len(newCol.name) != 0 && len(newCol.name) > w {
					panic(fmt.Sprintf("table td: tag - name (%s) is longer (%d) than colum width (%d)", newCol.name, len(newCol.name), w))
				}
				newCol.width = w
				newCol.visible = true

			case "sep", "s":
				if isntEnclosed(val, '\'') {
					panic("table: td: sep must be at least 1 char within ' '")
				}
				newCol.separator = strings.Trim(val, "'")

			case "algnr", "r":
				newCol.rightAlign = true

			case "invld", "i":
				if isntEnclosed(val, '\'') {
					panic("table td: inv - name must be at least 1 char within ' '")
				}
				newCol.invalidStr = strings.Trim(val, "'")
			}
		}
		if newCol.width == 0 && newCol.visible {
			panic(fmt.Sprintf("table: td: Invalid column definition: %s", newCol.name))
		}
		newTable.cols[colN] = newCol
	}
	newTable.width = newTable.buildHeader()
	return newTable
}

func (t *Table[Data]) buildHeader() int {
	hdr := "|"

	if t.displNos {
		hdr += "   #"
	}

	for _, column := range t.cols {
		if column.visible {

			hdr += " "
			if column.separator != "" {
				hdr += column.separator + " "
			}

			if !column.rightAlign { // put column name left of right aligned
				hdr += column.name
			}
			hdr += strings.Repeat(" ", column.width-len(column.name))
			if column.rightAlign {
				hdr += column.name
			}
		}
	}
	hdr += "|"

	t.header = hdr
	return len(t.header)
}

func (t *Table[Data]) GetWidth() int { return t.width }

func (t *Table[Data]) Print() {
	separator := strings.Repeat("-", len(t.header))

	fmt.Println(separator)
	fmt.Println(t.header)
	fmt.Println(separator)

	for i := 0; i < len(t.rows); i++ {
		fmt.Println(t.buildRow(i))
	}

	fmt.Println(separator)
}

func (t *Table[Data]) cellToString(row, col int) string {
	var (
		rowRV = reflect.ValueOf(t.rows[row]) // reflect value (type) of row structure

		vType = rowRV.Field(col).Type().Kind() // column type
		value = rowRV.Field(col)               // column value
		w     = t.cols[col].width              // column width
		valid = rowRV.FieldByName("valid")     // value of the 'valid" field in the row structure
	)

	if !valid.Bool() && t.cols[col].invalidStr != "" {
		return t.cols[col].invalidStr
	}

	str := ""
	switch vType {

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		format := fmt.Sprintf("%%%dd", w)
		str = fmt.Sprintf(format, value.Uint())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		format := fmt.Sprintf("%%%dd", w)
		str = fmt.Sprintf(format, value.Int())

	case reflect.Float32, reflect.Float64:
		format := fmt.Sprintf("%%%d.2f", w)
		str = fmt.Sprintf(format, float32(value.Float()))

	case reflect.String:
		str = value.String()

	case reflect.Bool:
		if value.Bool() {
			str = "x"
		} else {
			str = " "
		}
	}
	return str
}

func (t Table[Data]) buildRow(rowNo int) string {
	var a Data
	typ := reflect.TypeOf(a)

	// start
	line := "| "

	// add line number
	if t.displNos { // add line numbers
		line += fmt.Sprintf("%3d", rowNo)
	}

	for col := 0; col < typ.NumField(); col++ {
		var (
			c      = t.cols[col]
			s, str string
		)

		if c.visible {
			str = t.cellToString(rowNo, col)
			// add separator
			if c.separator != "" {
				line += " " + c.separator + " "
			} else {
				line += " "
			}

			// clip value to width
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

			// add value
			if c.rightAlign {
				line += strings.Repeat(" ", pad)
			}
			line += s
			if !c.rightAlign {
				line += strings.Repeat(" ", pad)
			}
		}
	}
	// end
	line += "|"
	return line
}

func isntEnclosed(s string, c byte) bool {
	return len(s) <= 2 || s[0] != c || s[len(s)-1] != c
}

func getKVPair(tag string) (string, string) {
	var key, val string
	kvPair := strings.Split(strings.TrimSpace(tag), "=")

	if len(kvPair) >= 1 {
		key = strings.TrimSpace(kvPair[0])
	}
	if len(kvPair) == 2 {
		val = strings.TrimSpace(kvPair[1])
	}
	if len(kvPair) > 2 {
		panic(fmt.Sprintf("table td: tag - key=value pair invalid: '%s'", tag))
	}
	return key, val
}
