package main

type Book struct {
	name  string
	score float32
}

type BookSlice []Book

func (p BookSlice) Len() int {
	return len(p)
}

func (p BookSlice) Less(i, j int) bool {
	return p[i].score < p[j].score
}

func (p BookSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
