package main

type pinger struct {
	url     string  `td:"width= 25, name='website', sep='|'"`
	lastVal float64 `td:"width= 9,  name='ping (ms)', algnR, invld='?',sep='>'"`
	min     float64 `td:"width= 8,  invld='?', R"`
	max     float64 `td:"width= 8,  invld='?', R"`
	avg     float64 `td:"width= 8,  invld='?', R"`
	status  string  `td:"name='errors',  w= 60,  sep='|'"`
	count   uint64
	sum     float64
	valid   bool
}

func (p pinger) IsValid() bool { return p.valid }

func (p *pinger) Invalidate(msg string) {
	p.status = msg
	p.valid = false
}
func (p *pinger) SetStatus(s string) {
	p.status = s
}

func (p *pinger) SetVal(v float64) {
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
	p.avg = p.sum / float64(p.count)
	p.status = ""
}
