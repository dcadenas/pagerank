package pagerank

type Float float32;

type Interface interface {
  Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float))
}

type pageRank struct {
}

func New() *pageRank {
  rg := new(pageRank)
  return rg;
}

func (v *pageRank) Link(from, to int) {
}

func (v *pageRank) Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float)) {
  resultFunc(1, 1.0)
}
