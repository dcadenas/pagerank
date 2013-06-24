package pagerank

import "fmt"
import "math"

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
  missingSlots := (fromAsIndex + 1) - cap(pr.numberOutLinks)

  if missingSlots > 0 {
    pr.numberOutLinks = append(pr.numberOutLinks, make([]int, missingSlots)...)
  }

  pr.numberOutLinks[fromAsIndex] += 1
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

func (pr *pageRank) calculateDanglingNodes() []int {
  danglingNodes := make([]int, 0, len(pr.numberOutLinks))

  for i := range pr.numberOutLinks {
    if pr.numberOutLinks[i] == 0 {
      danglingNodes = append(danglingNodes, i)
    }
  }

  return danglingNodes
}

func (pr *pageRank) step(followingProb, tOverSize Float, p []Float, danglingNodes []int) []Float {
  innerProduct := Float(0)

  for i := range danglingNodes {
    innerProduct += p[danglingNodes[i]]
  }

  innerProductOverSize := innerProduct / Float(len(p))
  vsum := Float(0)
  v := make([]Float, len(p))

  for i := range pr.inLinks {
    ksum := Float(0)
    inLinksForI := pr.inLinks[i]

    for j := range inLinksForI {
      index := inLinksForI[j]
      fmt.Println(p[index])
      ksum += p[index] / Float(pr.numberOutLinks[index])
    }

    v[i] = followingProb * (ksum + innerProductOverSize) + tOverSize
    vsum += v[i]
  }

  inverseOfSum := 1.0 / vsum

  for i := range v {
    v[i] *= inverseOfSum
  }

  return v
}

func calculateChange(p , new_p []Float) Float {
  acc := Float(0)

  for i := range p {
    acc += Float(math.Abs(float64(p[i] - new_p[i])))
  }

  return acc
}

func (pr *pageRank) Rank(followingProb, tolerance Float, resultFunc func(label int, rank Float)) {
  size := len(pr.keyToIndex)
  inverseOfSize := 1.0 / Float(size)
  tOverSize := (1.0 - followingProb) / Float(size)
  danglingNodes := pr.calculateDanglingNodes()

  p := make([]Float, size)
  for i := range p {
    p[i] = inverseOfSize
  }

  change := Float(2.0)

  for change > tolerance {
    new_p := pr.step(followingProb, tOverSize, p, danglingNodes)
    change = calculateChange(p, new_p)
    p = new_p
  }

  for i := range p {
    resultFunc(pr.indexToKey[i], p[i])
  }
}
