package pagerank

import (
  "testing"
  "math"
)

func round(float float64) float64 {
  return math.Floor(float * 10 + 0.5) / 10
}

func rankToPercentage(rank float32) float64 {
  tenPow3 := math.Pow(10, 3)
  return round(100 * (float64(rank) * tenPow3) / tenPow3)
}

func assertRank(t *testing.T, pageRank Interface, expected map[int]float64) {
  pageRank.Rank(0.85, 0.0001, func(label int, rank float32) {
    if rankToPercentage(rank) != expected[label] {
      t.Fail()
    }
  })
}

func assertEqual(t *testing.T, actual, expected interface{}) {
  if actual != expected {
    t.Error("Should be", expected, "but was", actual)
  }
}

func assert(t *testing.T, actual bool) {
  if !actual {
    t.Error("Should be true")
  }
}

func TestRound(t *testing.T) {
  assertEqual(t, round(0.6666666), 0.7)
}

func TestRankToPercentage(t *testing.T) {
  assertEqual(t, rankToPercentage(0.6666666), 66.7)
}

func TestShouldEnterTheBlock(t *testing.T) {
  pageRank := New()
  pageRank.Link(0, 1)

  entered := false
  pageRank.Rank(0.85, 0.0001, func(_ int, _ float32) {
    entered = true
  })

  assert(t, entered)
}

func TestShouldGiveCorrectResultForAOneLinkGraph(t *testing.T) {
  pageRank := New()
  pageRank.Link(0, 1)

  assertRank(t, pageRank, map[int]float64{0: 35.1, 1: 64.9})
}
