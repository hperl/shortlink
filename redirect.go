package main

type redirect struct {
	From string
	To   string
}

type byFrom []*redirect

func (f byFrom) Len() int           { return len(f) }
func (f byFrom) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byFrom) Less(i, j int) bool { return f[i].From < f[j].From }
