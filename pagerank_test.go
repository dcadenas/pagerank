package pagerank

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func assertRank(t *testing.T, pageRank Interface, expected map[string]int64) {
	const tolerance = Dot4ONE
	pageRank.Rank(85*Dot2ONE, tolerance, func(label string, rank int64) {
		rankAsPercentage := toPercentage(rank)
		if Abs(rankAsPercentage-expected[label]) > tolerance {
			t.Error("Rank for", label, "should be", expected[label], "but was", rankAsPercentage)
		}
	})
}

func round(f int64) int64 {
	return (f*10 + 5*DotONE) / ONE * ONE / 10
}

func toPercentage(f int64) int64 {
	tenPow3 := int64(1000)
	return round(100 * f * tenPow3 / tenPow3)
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
	assertEqual(t, round(6666666*Dot7ONE), 7*DotONE)
}

func TestRankToPercentage(t *testing.T) {
	assertEqual(t, toPercentage(6666666*Dot7ONE), 667*DotONE)
}

func TestShouldEnterTheBlock(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "1")

	entered := false
	pageRank.Rank(85*Dot2ONE, 1*Dot4ONE, func(_ string, _ int64) {
		entered = true
	})

	assert(t, entered)
}

func TestShouldBePossibleToRecalculateTheRanksAfterANewLinkIsAdded(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "1")
	assertRank(t, pageRank, map[string]int64{"0": 351 * DotONE, "1": 649 * DotONE})
	pageRank.Link("1", "2")
	assertRank(t, pageRank, map[string]int64{"0": 184 * DotONE, "1": 341 * DotONE, "2": 474 * DotONE})
}

func TestShouldBePossibleToClearTheGraph(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "1")
	pageRank.Link("1", "2")
	pageRank.Clear()
	pageRank.Link("0", "1")
	assertRank(t, pageRank, map[string]int64{"0": 351 * DotONE, "1": 649 * DotONE})
}

func TestShouldNotFailWhenCalculatingTheRankOfAnEmptyGraph(t *testing.T) {
	pageRank := New()
	pageRank.Rank(85*Dot2ONE, 00001*Dot4ONE, func(label string, rank int64) {
		t.Error("This should not be seen")
	})
}

func TestShouldReturnCorrectResultsWhenHavingADanglingNode(t *testing.T) {
	pageRank := New()
	//node 2 is a dangling node because it has no outbound links
	pageRank.Link("0", "2")
	pageRank.Link("1", "2")

	expectedRank := map[string]int64{
		"0": 213 * DotONE,
		"1": 213 * DotONE,
		"2": 574 * DotONE,
	}

	assertRank(t, pageRank, expectedRank)
}

func TestShouldNotChangeTheGraphWhenAddingTheSameLinkManyTimes(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "2")
	pageRank.Link("0", "2")
	pageRank.Link("0", "2")
	pageRank.Link("1", "2")
	pageRank.Link("1", "2")

	expectedRank := map[string]int64{
		"0": 213 * DotONE,
		"1": 213 * DotONE,
		"2": 574 * DotONE,
	}

	assertRank(t, pageRank, expectedRank)
}

func TestShouldReturnCorrectResultsForAStarGraph(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "2")
	pageRank.Link("1", "2")
	pageRank.Link("2", "2")

	expectedRank := map[string]int64{
		"0": 5 * ONE,
		"1": 5 * ONE,
		"2": 90 * ONE,
	}

	assertRank(t, pageRank, expectedRank)
}

func TestShouldBeUniformForACircularGraph(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "1")
	pageRank.Link("1", "2")
	pageRank.Link("2", "3")
	pageRank.Link("3", "4")
	pageRank.Link("4", "0")

	expectedRank := map[string]int64{
		"0": 20 * ONE,
		"1": 20 * ONE,
		"2": 20 * ONE,
		"3": 20 * ONE,
		"4": 20 * ONE,
	}

	assertRank(t, pageRank, expectedRank)
}

func TestShouldReturnCorrectResultsForAConvergingGraph(t *testing.T) {
	pageRank := New()
	pageRank.Link("0", "1")
	pageRank.Link("0", "2")
	pageRank.Link("1", "2")
	pageRank.Link("2", "2")

	expectedRank := map[string]int64{
		"0": 5 * ONE,
		"1": 71 * DotONE,
		"2": 879 * DotONE,
	}

	assertRank(t, pageRank, expectedRank)
}

func TestShouldCorrectlyReproduceTheWikipediaExample(t *testing.T) {
	//http://en.wikipedia.org/wiki/File:PageRanks-Example.svg
	pageRank := New()
	pageRank.Link("1", "2")
	pageRank.Link("2", "1")
	pageRank.Link("3", "0")
	pageRank.Link("3", "1")
	pageRank.Link("4", "3")
	pageRank.Link("4", "1")
	pageRank.Link("4", "5")
	pageRank.Link("5", "4")
	pageRank.Link("5", "1")
	pageRank.Link("6", "1")
	pageRank.Link("6", "4")
	pageRank.Link("7", "1")
	pageRank.Link("7", "4")
	pageRank.Link("8", "1")
	pageRank.Link("8", "4")
	pageRank.Link("9", "4")
	pageRank.Link("10", "4")

	expectedRank := map[string]int64{
		"0":  33 * DotONE,  //a
		"1":  384 * DotONE, //b
		"2":  343 * DotONE, //c
		"3":  39 * DotONE,  //d
		"4":  81 * DotONE,  //e
		"5":  39 * DotONE,  //f
		"6":  16 * DotONE,  //g
		"7":  16 * DotONE,  //h
		"8":  16 * DotONE,  //i
		"9":  16 * DotONE,  //j
		"10": 16 * DotONE,  //k
	}

	assertRank(t, pageRank, expectedRank)
}

func BenchmarkOneMillion(b *testing.B) {
	n := 1000000

	pageRank := New()

	rand.Seed(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for from := 0; from < n; from++ {
			for j := 0; j < rand.Intn(400); j++ {
				too := rand.Intn(n)

				to := too
				if too > 800000 {
					to = rand.Intn(3)
				}

				pageRank.Link(strconv.FormatInt(int64(from), 10), strconv.FormatInt(int64(to), 10))
			}
		}
	}

	result := make(map[string]int64, n)
	pageRank.Rank(85*Dot2ONE, 0001*Dot3ONE, func(key string, val int64) {
		result[key] = val
	})

	fmt.Println("5 first values are", result["0"], ",", result["1"], ",", result["2"], ",", result["3"], ",", result["4"])
	pageRank.Clear()
}
