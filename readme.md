## Pinger ##
* code to test several golang topics and techniques
-----

  - custom tags for print formating & reflection
    The idea is so you could add the following custom tags to a struct, so it can be nicely printed as a formatted table
    Tags are: (all start with td: - for Table Data)
    -   width or w - is the column width, text not fitting into the column will be abbreviated with '...'
    -   name or n is the alternative name for the column, if not given the struct name will be used
    -   sep or s - optional separator prited on the left side of the column
    -   algnR or R or r - indicates the column is to be right-adjusted
    -   invld or i - alt text to be printed in each column, if the whole row is invalid

```go
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
```
prints this
```
-----------------------------------------------------------------------------------------------------
|   # | website              > ping (ms)      min      max      avg | errors                        |
-----------------------------------------------------------------------------------------------------
|   0 | http://www.onet.pl   >      1320     1320     1320  1320.00 |                               |
|   1 | http://www.gazeta.pl >       139      139      139   139.00 |                               |
|   2 |                      >         ?        ?        ?        ? | connecting...                 |
|   3 | http://www.google... >       673      673      673   673.00 |                               |
|   4 |                      >         ?        ?        ?        ? | connecting...                 |
|   5 | http://www.icm.ed... >       548      548      548   548.00 |                               |
-----------------------------------------------------------------------------------------------------
```
  used reflection to print formatted table

  - generics 
  table.Table is a generic struct, which be instantiated on any type, which implements 'validable' interface
  (the 'validatable' interface only requires one function: IsValid() bool to be available)
  
  - concurrency and multiprocess communication (contex, cancel, channels, waitgroups, thread-safe logger)

  there are several concurrent processes running
    - one performig periodic measurements
    - another periodically displaying table with data
    - etc.
  
  The idea was to ensure graceful shutdown of all these processes if
    - program lifetime expires
    - ctrl-C is pressed


