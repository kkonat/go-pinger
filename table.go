package main

import (
	"fmt"
	w "pinger/window"
)

type tableItem struct {
	count   uint64
	lastVal uint64
	min     uint64
	max     uint64
	sum     uint64
	status  string
}

func (d *tableItem) update(v uint64) {
	if d.count == 0 {
		d.min = v
		d.max = v
	} else {
		if d.min > v {
			d.min = v
		} else if d.max < v {
			d.max = v
		}
	}
	d.lastVal = v
	d.sum += v
	d.count++
}

type table struct {
	data []tableItem
}

func NewTable() *table {
	t := &table{}
	t.data = make([]tableItem, len(urls))
	return t
}

func (t table) displayHeader() {
	w.ClearScreen()
	w.Gotoxy(0, 0)
	fmt.Printf("#\turl\t\t\t\tping\tmin\tmax\tavg\terrors")
	w.Gotoxy(0, 2)
	fmt.Printf("--------------------------------------------------------------------------------")
}
func (t table) displayRow(row int) {
	w.Gotoxy(0, row+3)
	ti := t.data[row]
	if ti.count != 0 {
		fmt.Printf("%d\t%s\t\t%d  \t%d  \t%d  \t%.2f  \t%s", row+1, urls[row], ti.lastVal, ti.min, ti.max, (float64(ti.sum) / float64(ti.count)), ti.status)
	} else {
		fmt.Printf("%d\t%s\t\t?  \t?  \t?  \t?  \t%s", row+1, urls[row], ti.status)
	}

}

const clipLen = 80

func (t *table) update(meas measurement) {

	if meas.status != "OK" {

		if len(meas.status) <= clipLen {
			t.data[meas.urlId].status = meas.status
		} else {
			t.data[meas.urlId].status = meas.status[:clipLen] + "..."
		}

	} else {
		t.data[meas.urlId].update(meas.v)
	}
}
