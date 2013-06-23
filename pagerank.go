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


func (pr *pageRank) keyAsArrayIndex(key int) int {
  index, ok := pr.keyToIndex[key]

  if !ok {
    pr.currentAvailableIndex++
    index = pr.currentAvailableIndex
    pr.keyToIndex[key] = index
    pr.indexToKey[index] = key
  }

  return index
}

func (pr *pageRank) updateInLinks(fromAsIndex, toAsIndex int) {
  missingSlots := (toAsIndex + 1) - cap(pr.inLinks)

  if missingSlots > 0 {
    pr.inLinks = append(pr.inLinks, make([][]int, missingSlots)...)
  }

  pr.inLinks[toAsIndex] = append(pr.inLinks[toAsIndex], fromAsIndex)
}

func (pr *pageRank) updateNumberOutLinks(fromAsIndex int) {
  setIntSlice(pr.numberOutLinks, fromAsIndex, func(oldValue int)int{
    return oldValue + 1
  })
}

func setIntSlice(slice []int, index int, valueFunc func(int)int) {
  missingSlots := (index + 1) - cap(slice)

  if missingSlots > 0 {
    slice = append(slice, make([]int, missingSlots)...)
  }

  slice[index] = valueFunc(slice[index])
}

func (pr *pageRank) linkWithIndices(fromAsIndex, toAsIndex int) {
  pr.updateInLinks(fromAsIndex, toAsIndex)
  pr.updateNumberOutLinks(fromAsIndex)
}

func (pr *pageRank) Link(from, to int) {
  fromAsIndex := pr.keyAsArrayIndex(from)
  toAsIndex := pr.keyAsArrayIndex(to)

  pr.linkWithIndices(fromAsIndex, toAsIndex)
}

func (pr *pageRank) Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float)) {
  resultFunc(1, 1.0)
}
