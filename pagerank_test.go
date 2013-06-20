package pagerank

import (
  "testing"
  "math"
)

func (f Float) round() Float {
  return Float(math.Floor(float64(f) * 10 + 0.5) / 10)
}

func (f Float) toPercentage() Float {
  tenPow3 := math.Pow(10, 3)
  return Float(100 * (float64(f) * tenPow3) / tenPow3).round()
}

func assertRank(t *testing.T, pageRank Interface, expected map[int]Float) {
  pageRank.Rank(0.85, 0.0001, func(label int, rank Float) {
    if Float(rank).toPercentage() != expected[label] {
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
  assertEqual(t, Float(0.6666666).round(), Float(0.7))
}

func TestRankToPercentage(t *testing.T) {
  assertEqual(t, Float(0.6666666).toPercentage(), Float(66.7))
}

func TestShouldEnterTheBlock(t *testing.T) {
  pageRank := New()
  pageRank.Link(0, 1)

  entered := false
  pageRank.Rank(0.85, 0.0001, func(_ int, _ Float) {
    entered = true
  })

  assert(t, entered)
}

func TestShouldGiveCorrectResultForAOneLinkGraph(t *testing.T) {
  pageRank := New()
  pageRank.Link(0, 1)

  assertRank(t, pageRank, map[int]Float{0: 35.1, 1: 64.9})
}
