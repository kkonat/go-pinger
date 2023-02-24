package main

type pinger struct {
	url     string  `td:"w= 30, name='website',     sep='|'"`
	lastVal uint64  `td:"w= 9, name='ping (ms)', R, inv='?',sep='>'"`
	min     uint64  `td:"w= 8, inv='?', R"`
	max     uint64  `td:"w= 8, inv='?', R"`
	avg     float64 `td:"w= 8, inv='?', R"`
	status  string  `td:"name='errors',  w= 60, inv='connecting...',  sep='|'"`
	count   uint64
	sum     uint64
	valid   bool
}

func (p pinger) IsValid() bool { return p.valid }
