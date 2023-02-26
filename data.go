package main

import "math/rand"

type pinger struct {
	url     string  `td:"width= 20, name='website', sep='|'"`
	lastVal uint64  `td:"width= 9,  name='ping (ms)', algnR, invld='?',sep='>'"`
	min     uint64  `td:"width= 8,  invld='?', R"`
	max     uint64  `td:"width= 8,  invld='?', R"`
	avg     float64 `td:"width= 8,  invld='?', R"`
	status  string  `td:"name='errors',  w= 30, invld='connecting...',  sep='|'"`
	count   uint64
	sum     uint64
	valid   bool
}

func (p pinger) IsValid() bool { return p.valid }
func (p *pinger) setVal(v uint64) {
	p.valid = true
	p.lastVal = v
	if p.count == 0 || p.min > v {
		p.min = v
	}
	if p.count == 0 || p.max < v {
		p.max = v
	}
	p.count++
	p.sum += v
	p.avg = float64(p.sum) / float64(p.count)
	if rand.Intn(1000) > 938 {
		p.status = "error"
	} else {
		p.status = ""
	}
}
