package main

type Rank struct {
	name    string
	similar float32
}

type RankSlice []Rank

func (p RankSlice) Len() int {
	return len(p)
}

func (p RankSlice) Less(i, j int) bool {
	return p[i].similar < p[j].similar
}

func (p RankSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
