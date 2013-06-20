package pagerank

type Interface interface {
  Rank(followingProb, tolerance float32, resultFunc func(label int, rank float32))
}

type pageRank struct {
}

func New() *pageRank {
  rg := new(pageRank)
  return rg;
}

func (v *pageRank) Link(from, to int) {
}

func (v *pageRank) Rank(followingProb, tolerance float32, resultFunc func(label int, rank float32)) {
  resultFunc(1, 1.0)
}
