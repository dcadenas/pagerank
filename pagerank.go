package pagerank

type Float float32;

type Interface interface {
  Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float))
}

type pageRank struct {
  inLinks [][]int
  numberOutLinks []int
  currentAvailableIndex int
  keyToIndex map[int]int
  indexToKey map[int]int
}

func New() *pageRank {
  rg := new(pageRank)
  rg.inLinks = [][]int{}
  rg.numberOutLinks = []int{}
  rg.currentAvailableIndex = -1
  rg.keyToIndex = make(map[int]int)
  rg.indexToKey = make(map[int]int)
  return rg;
}

func (pr *pageRank) Link(from, to int) {
}

func (pr *pageRank) Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float)) {
  resultFunc(1, 1.0)
}
