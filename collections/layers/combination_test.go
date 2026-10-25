package layers

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"generator/internal/nft"
)

// layerIndicesByIncrement is the former conversion of a combination number, kept as the reference:
// it starts before the first combination and counts up to the number one step at a time, O(C).
func layerIndicesByIncrement(combination int, imageCounts []int) []int {
	indices := make([]int, len(imageCounts))
	last := len(imageCounts) - 1
	for k := 0; k < last; k++ {
		indices[k] = 1
	}

	for i := 0; i < combination; i++ {
		if indices[last] < imageCounts[last] {
			indices[last]++
		} else {
			setNextLevel(indices, imageCounts, last)
		}
	}
	return indices
}

func setNextLevel(indices, imageCounts []int, layer int) {
	indices[layer] = 1
	if indices[layer-1] < imageCounts[layer-1] {
		indices[layer-1]++
	} else {
		setNextLevel(indices, imageCounts, layer-1)
	}
}

func TestLayerIndices(t *testing.T) {
	for _, counts := range [][]int{{7}, {2, 3}, {3, 3, 2}, {4, 1, 5}, {2, 3, 4, 5}} {
		total := 1
		for _, count := range counts {
			total *= count
		}
		require.Equal(t, int64(total), combinationCount(counts).Int64())

		seen := make(map[[4]int]bool)
		for combination := 1; combination <= total; combination++ {
			number := big.NewInt(int64(combination))
			indices := layerIndices(number, counts)
			require.Equal(t, layerIndicesByIncrement(combination, counts), indices, "counts %v, combination %d", counts, combination)
			require.Equal(t, number, combinationNumber(indices, counts), "counts %v, combination %d", counts, combination)

			var key [4]int
			for k, index := range indices {
				require.True(t, index >= 1 && index <= counts[k], "counts %v, combination %d", counts, combination)
				key[k] = index
			}
			require.False(t, seen[key], "counts %v: combination %d repeats an earlier one", counts, combination)
			seen[key] = true
		}
	}
}

func TestLayerIndicesBeyondInt64(t *testing.T) {
	// 28 layers of 7 variants: C = 7^28 ~4.6e23, far beyond the largest int64.
	bases := make([]int, 28)
	for k := range bases {
		bases[k] = 7
	}
	total := combinationCount(bases)
	require.Equal(t, 1, total.Cmp(big.NewInt(math.MaxInt64)))

	last := layerIndices(total, bases)
	for _, index := range last {
		require.Equal(t, 7, index)
	}
	require.Equal(t, total, combinationNumber(last, bases))

	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		number := new(big.Int).Rand(rnd, total)
		number.Add(number, big.NewInt(1))
		require.Equal(t, number, combinationNumber(layerIndices(number, bases), bases))
	}
}

func TestLayerIndicesAroundInt64(t *testing.T) {
	// both ways of conversion agree where numbers stop fitting in int64: 2^63 - 1 and beyond.
	bases := []int{1 << 20, 1 << 20, 1 << 20, 16}
	limit := big.NewInt(math.MaxInt64)
	for _, delta := range []int64{-2, -1, 0, 1, 2} {
		number := new(big.Int).Add(limit, big.NewInt(delta))
		require.Equal(t, number, combinationNumber(layerIndices(number, bases), bases), "2^63-1%+d", delta)
	}

	small := []int{3, 5, 7, 11}
	for number := int64(1); number <= 3*5*7*11; number++ {
		require.Equal(t, number, combinationNumber64(layerIndices(big.NewInt(number), small), small))
	}
}

func TestSampleCombinations(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct{ amount, total int }{{0, 5}, {1, 1}, {1, 10}, {5, 10}, {10, 10}, {1000, 1000}, {300, 100000}} {
		combinations := sampleCombinations(rnd, test.amount, big.NewInt(int64(test.total)))
		require.Len(t, combinations, test.amount)

		seen := make(map[int64]bool)
		for _, combination := range combinations {
			number := combination.Int64()
			require.True(t, number >= 1 && number <= int64(test.total), "%+v: %d", test, number)
			require.False(t, seen[number], "%+v: %d repeats", test, number)
			seen[number] = true
		}
	}

	// every combination of a small set turns up first about equally often.
	const total, runs = 4, 40000
	first := make([]int, total+1)
	for i := 0; i < runs; i++ {
		first[sampleCombinations(rnd, total, big.NewInt(total))[0].Int64()]++
	}
	for combination := 1; combination <= total; combination++ {
		assert.InDelta(t, runs/total, first[combination], runs/total*0.05, "combination %d", combination)
	}

	// numbers beyond int64 are drawn as well.
	huge := new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil)
	for _, combination := range sampleCombinations(rnd, 100, huge) {
		require.True(t, combination.Sign() > 0 && combination.Cmp(huge) <= 0)
	}
}

func TestAliasTable(t *testing.T) {
	for _, weights := range [][]int{{1}, {5, 5}, {1, 2, 3, 4}, {100, 1, 1, 1, 1, 1}, {7, 300, 2, 2, 50, 1000, 13}, {MaxWeight, 1, MaxWeight}} {
		table := newAliasTable(weights)
		m := int64(len(weights))

		// the exact probability of variant i is (1/m) * sum over columns of its share of the column,
		// so m * sum * P(i) is an integer that must equal m * weight_i.
		mass := make([]int64, len(weights))
		for column := range weights {
			require.True(t, table.threshold[column] >= 0 && table.threshold[column] <= table.sum, "%v", weights)
			mass[column] += table.threshold[column]
			mass[table.alias[column]] += table.sum - table.threshold[column]
		}
		for i, weight := range weights {
			require.Equal(t, int64(weight)*m, mass[i], "weights %v, variant %d", weights, i)
		}
	}

	// and the draws follow the weights.
	weights := []int{60, 25, 10, 4, 1}
	table := newAliasTable(weights)
	rnd := rand.New(rand.NewSource(2))
	const runs = 200000
	counts := make([]int, len(weights))
	for i := 0; i < runs; i++ {
		counts[table.draw(rnd)]++
	}
	for i, weight := range weights {
		expected := float64(runs*weight) / 100
		assert.InDelta(t, expected, counts[i], 5*math.Sqrt(expected), "variant %d", i)
	}
}

func TestSampleWeighted(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))
	weights := [][]int{{70, 20, 10}, {1, 1}, {50, 30, 15, 5}, {90, 10}}
	bases := []int{3, 2, 4, 2}

	combinations, draws, ok := sampleWeighted(rnd, 48, weights, 1_000_000)
	require.True(t, ok)
	require.Len(t, combinations, 48)
	require.True(t, draws >= 48)
	seen := make(map[string]bool)
	for _, combination := range combinations {
		require.True(t, combination.Sign() > 0 && combination.Cmp(combinationCount(bases)) <= 0)
		require.False(t, seen[combination.String()], "%s repeats", combination)
		seen[combination.String()] = true
	}

	// with C >> N repeats are rare, and the variant frequencies follow the weights.
	even := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	wide := [][]int{{60, 25, 10, 4, 1}, even, even, even, even, even}
	const amount = 20000
	combinations, _, ok = sampleWeighted(rnd, amount, wide, 10*amount)
	require.True(t, ok)
	counts := make([]int, 5)
	for _, combination := range combinations {
		counts[layerIndices(combination, []int{5, 10, 10, 10, 10, 10})[0]-1]++
	}
	for i, weight := range wide[0] {
		expected := float64(amount*weight) / 100
		assert.InDelta(t, expected, counts[i], 5*math.Sqrt(expected), "variant %d", i)
	}

	// it gives up when the draws run out: all 4 combinations, one of them with probability 1/10000.
	_, draws, ok = sampleWeighted(rnd, 4, [][]int{{9999, 1}, {1, 1}}, 500)
	require.False(t, ok)
	require.Equal(t, 500, draws)
}

func TestChooseCombinations(t *testing.T) {
	rnd := rand.New(rand.NewSource(4))
	layer := func(noneWeight int, weights ...int) Layer {
		layer := Layer{NoneWeight: noneWeight, Weights: map[string]int{}}
		for i, weight := range weights {
			name := string(rune('a'+i)) + ".png"
			layer.ImageNames = append(layer.ImageNames, name)
			layer.Weights[name] = weight
		}
		layer.LayerImagesCount = len(layer.ImageNames)
		return layer
	}

	// equal weights: every combination, including the absent optional layer, by Floyd's algorithm.
	uniform := []Layer{layer(0, 100, 100), layer(100, 100, 100)}
	require.Equal(t, []int{2, 3}, layerBases(uniform))
	combinations, err := chooseCombinations(rnd, 6, uniform)
	require.NoError(t, err)
	require.Len(t, combinations, 6)

	// a rare optional layer is mostly absent.
	weighted := []Layer{layer(0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), layer(0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1), layer(900, 50, 50)}
	combinations, err = chooseCombinations(rnd, 50, weighted)
	require.NoError(t, err)
	absent := 0
	for _, combination := range combinations {
		if layerIndices(combination, layerBases(weighted))[2] == 3 {
			absent++
		}
	}
	require.True(t, absent > 30, "absent in %d of 50", absent)

	// too skewed for the size: the client gets an error instead of a stuck generation.
	_, err = chooseCombinations(rnd, 4, []Layer{layer(0, MaxWeight, 1), layer(0, MaxWeight, 1)})
	require.True(t, ErrRareCombinations.Has(err))
}

func TestAttributes(t *testing.T) {
	projectLayers := []Layer{
		{Name: "Background", ImageNames: []string{"Aurora.png", "Sunset.png"}},
		{Name: "Hat", ImageNames: []string{"Crown.png"}, NoneWeight: 10},
		{Name: "Eyes", ImageNames: []string{"Happy.png", "Sleepy.png"}},
	}
	require.Equal(t, []nft.Attribute{
		{TraitType: "Background", Value: "Sunset"},
		{TraitType: "Hat", Value: NoneName},
		{TraitType: "Eyes", Value: "Happy"},
	}, attributes(projectLayers, []int{2, 2, 1}))
}
