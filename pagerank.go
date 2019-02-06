package pagerank

const ONE = int64(10000000)
const DotONE = ONE / 10
const Dot2ONE = ONE / 100
const Dot3ONE = ONE / 1000
const Dot4ONE = ONE / 10000
const Dot7ONE = ONE / 10000000

type Interface interface {
	Rank(followingProb, tolerance int64, resultFunc func(label string, rank int64))
	Link(from, to string)
}

type pageRank struct {
	inLinks               [][]int
	numberOutLinks        []int
	currentAvailableIndex int
	keyToIndex            map[string]int
	indexToKey            map[int]string
}

func New() *pageRank {
	pr := new(pageRank)
	pr.Clear()
	return pr
}

func (pr *pageRank) keyAsArrayIndex(key string) int {
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
	missingSlots := len(pr.keyToIndex) - len(pr.inLinks)

	if missingSlots > 0 {
		pr.inLinks = append(pr.inLinks, make([][]int, missingSlots)...)
	}

	pr.inLinks[toAsIndex] = append(pr.inLinks[toAsIndex], fromAsIndex)
}

func (pr *pageRank) updateNumberOutLinks(fromAsIndex int) {
	missingSlots := len(pr.keyToIndex) - len(pr.numberOutLinks)

	if missingSlots > 0 {
		pr.numberOutLinks = append(pr.numberOutLinks, make([]int, missingSlots)...)
	}

	pr.numberOutLinks[fromAsIndex] += 1
}

func (pr *pageRank) linkWithIndices(fromAsIndex, toAsIndex int) {
	pr.updateInLinks(fromAsIndex, toAsIndex)
	pr.updateNumberOutLinks(fromAsIndex)
}

func (pr *pageRank) Link(from, to string) {
	fromAsIndex := pr.keyAsArrayIndex(from)
	toAsIndex := pr.keyAsArrayIndex(to)

	pr.linkWithIndices(fromAsIndex, toAsIndex)
}

func (pr *pageRank) calculateDanglingNodes() []int {
	danglingNodes := make([]int, 0, len(pr.numberOutLinks))

	for i, numberOutLinksForI := range pr.numberOutLinks {
		if numberOutLinksForI == 0 {
			danglingNodes = append(danglingNodes, i)
		}
	}

	return danglingNodes
}

func (pr *pageRank) step(followingProb, tOverSize int64, p []int64, danglingNodes []int) []int64 {
	innerProduct := int64(0)

	for _, danglingNode := range danglingNodes {
		innerProduct += int64(p[danglingNode])
	}

	innerProductOverSize := innerProduct / int64(len(p))
	vsum := int64(0)
	v := make([]int64, len(p))

	for i, inLinksForI := range pr.inLinks {
		ksum := int64(0)

		for _, index := range inLinksForI {
			ksum += p[index] / int64(pr.numberOutLinks[index])
		}

		v[i] = followingProb*(ksum+innerProductOverSize)/ONE + tOverSize
		vsum += v[i]
	}

	inverseOfSum := ONE * ONE / vsum

	for i := range v {
		v[i] = v[i] * inverseOfSum / ONE
	}

	return v
}

func Abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func calculateChange(p, new_p []int64) int64 {
	acc := int64(0)

	for i, pForI := range p {
		acc += Abs(pForI - new_p[i])
	}

	return acc
}

func (pr *pageRank) Rank(followingProb, tolerance int64, resultFunc func(label string, rank int64)) {
	size := len(pr.keyToIndex)
	if size == 0 {
		return
	}
	inverseOfSize := ONE / int64(size)
	tOverSize := (ONE - followingProb) / int64(size)
	danglingNodes := pr.calculateDanglingNodes()

	p := make([]int64, size)
	for i := range p {
		p[i] = inverseOfSize
	}

	change := 2 * ONE

	for change > tolerance {
		new_p := pr.step(followingProb, tOverSize, p, danglingNodes)
		change = calculateChange(p, new_p)
		p = new_p
	}

	for i, pForI := range p {
		resultFunc(pr.indexToKey[i], pForI)
	}
}

func (pr *pageRank) Clear() {
	pr.inLinks = [][]int{}
	pr.numberOutLinks = []int{}
	pr.currentAvailableIndex = -1
	pr.keyToIndex = make(map[string]int)
	pr.indexToKey = make(map[int]string)
}
