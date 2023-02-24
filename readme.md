## Pinger ##
* code to test several golang topics and techniques
-----

  - custom tags for print formating
    The idea is so you could add the following custom tags to a struct, so it can be nicely printed as a formatted table
    Tags are: (all start with td: - for Table Data)
    -   w is the column width, text not fitting into the column will be abbreviated with '...'
    -   name is the alternative name for the column, if not given the struct name will be used
    -   sep - optional separator prited on the left side of the column
    -   R - indicates the column is to be right-adjusted
    -   inv - alt text to be printed in each column, if the whole row is invalid

```go
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
```

  - generics
  table.Table is generic struct, can be instantiated on any type, which implements 'validatable' interface
  (the 'validatable' interface only requires one function: IsValid() bool to be available)
  
  - concurrency and multiprocess communication (contex, cancel, channels, waitgroups, thread-safe logger)

  there are several concurrent processes running
    - one performig periodic measurements
    - another periodically displaying table with data
    - etc.
  
  The idea was to ensure graceful shutdown of all these processes if
    - program lifetime expires
    - ctrl-C is pressed


