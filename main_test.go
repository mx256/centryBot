package main

import "testing"

func Test_SaveToDb(t *testing.T) {
	InitDb()
	m := map[int]int64{
		1111: 1111111,
		2222: 2222222,
		3333: 333333,
		4444: 4545545,
		375:  98273982,
		3985: 903840,
		3487: 487932,
		3432: 948939,
	}
	for key, val := range m {
		SaveToDb(key, val)
		read, err := SearchInDb(key)
		if read != m[key] || err != nil {
			t.Error()
		}
	}
}
