package layers

import (
	"math/big"
	"math/rand"
)

const (
	// DefaultWeight is the weight of a variant whose weight was never set; equal weights mean equal chances.
	DefaultWeight = 100
	// MaxWeight bounds a weight, so a layer's weights scaled in the alias table stay far within int64.
	MaxWeight = 1_000_000
	// maxDrawsPerImage bounds the weighted choice: amount images may take at most this many draws each on
	// average before the generation gives up on finding enough distinct combinations.
	maxDrawsPerImage = 50
)

// layerBases returns the radix of every layer in the layer order: its number of images, plus one for the
// "none" variant of an optional layer, which comes after the images.
func layerBases(projectLayers []Layer) []int {
	bases := make([]int, len(projectLayers))
	for k, layer := range projectLayers {
		bases[k] = layer.LayerImagesCount
		if layer.NoneWeight > 0 {
			bases[k]++
		}
	}
	return bases
}

// layerWeights returns the weights of every layer's variants in the order of layerBases: the images in
// the order of their file names, then the "none" variant of an optional layer.
func layerWeights(projectLayers []Layer) [][]int {
	weights := make([][]int, len(projectLayers))
	for k, layer := range projectLayers {
		for _, name := range layer.ImageNames {
			weights[k] = append(weights[k], layer.Weight(name))
		}
		if layer.NoneWeight > 0 {
			weights[k] = append(weights[k], layer.NoneWeight)
		}
	}
	return weights
}

// combinationCount returns the number of combinations C, the product of the bases. With tens of layers
// it exceeds the largest int64 (28 layers of 7 variants give about 4.6 · 10^23), so it is a big.Int.
func combinationCount(bases []int) *big.Int {
	total := big.NewInt(1)
	for _, base := range bases {
		total.Mul(total, big.NewInt(int64(base)))
	}
	return total
}

// layerIndices returns the variant number (from 1) of every layer for the combination with the given
// number (from 1 to the product of the bases). The number minus one is written in the mixed radix numeral
// system whose bases are the layers' variant counts; the last layer is the lowest digit.
// It takes O(n) operations for n layers instead of counting up to the number one step at a time, O(C);
// a number that fits in int64 is converted in machine integers, a larger one in big.Int.
func layerIndices(combination *big.Int, bases []int) []int {
	indices := make([]int, len(bases))
	if combination.IsInt64() {
		rest := combination.Int64() - 1
		for k := len(bases) - 1; k >= 0; k-- {
			indices[k] = int(rest%int64(bases[k])) + 1
			rest /= int64(bases[k])
		}
		return indices
	}

	rest := new(big.Int).Sub(combination, big.NewInt(1))
	base, digit := new(big.Int), new(big.Int)
	for k := len(bases) - 1; k >= 0; k-- {
		rest.DivMod(rest, base.SetInt64(int64(bases[k])), digit)
		indices[k] = int(digit.Int64()) + 1
	}
	return indices
}

// combinationNumber is the inverse of layerIndices: it reads the variant numbers (from 1) as the digits
// of the mixed radix number, the first layer being the highest digit (Horner's scheme), and adds one.
func combinationNumber(indices, bases []int) *big.Int {
	number := new(big.Int)
	base, digit := new(big.Int), new(big.Int)
	for k, index := range indices {
		number.Mul(number, base.SetInt64(int64(bases[k])))
		number.Add(number, digit.SetInt64(int64(index-1)))
	}
	return number.Add(number, big.NewInt(1))
}

// combinationNumber64 is combinationNumber for the bases whose product fits in int64.
func combinationNumber64(indices, bases []int) int64 {
	var number int64
	for k, index := range indices {
		number = number*int64(bases[k]) + int64(index-1)
	}
	return number + 1
}

// sampleCombinations returns amount distinct combination numbers from 1 to total in random order.
// Robert Floyd's algorithm makes exactly amount draws however close amount is to total, while drawing
// until an unused number turns up needs about total·ln(total) draws to take every combination.
func sampleCombinations(rnd *rand.Rand, amount int, total *big.Int) []*big.Int {
	if total.IsInt64() {
		combinations := make([]*big.Int, amount)
		for i, combination := range sampleCombinations64(rnd, amount, total.Int64()) {
			combinations[i] = big.NewInt(combination)
		}
		return combinations
	}

	chosen := make(map[string]struct{}, amount)
	combinations := make([]*big.Int, 0, amount)
	one := big.NewInt(1)
	j := new(big.Int).Sub(total, big.NewInt(int64(amount)))
	for i := 0; i < amount; i++ {
		j = new(big.Int).Add(j, one)
		combination := new(big.Int).Rand(rnd, j)
		combination.Add(combination, one)
		if _, ok := chosen[string(combination.Bytes())]; ok {
			combination = j
		}
		chosen[string(combination.Bytes())] = struct{}{}
		combinations = append(combinations, combination)
	}

	// the set is uniform, but the order is not: j is taken in increasing order.
	rnd.Shuffle(len(combinations), func(a, b int) {
		combinations[a], combinations[b] = combinations[b], combinations[a]
	})
	return combinations
}

// sampleCombinations64 is Floyd's algorithm in machine integers for total that fits in int64.
func sampleCombinations64(rnd *rand.Rand, amount int, total int64) []int64 {
	chosen := make(map[int64]struct{}, amount)
	combinations := make([]int64, 0, amount)
	for j := total - int64(amount) + 1; j <= total; j++ {
		combination := rnd.Int63n(j) + 1
		if _, ok := chosen[combination]; ok {
			combination = j
		}
		chosen[combination] = struct{}{}
		combinations = append(combinations, combination)
	}

	rnd.Shuffle(len(combinations), func(a, b int) {
		combinations[a], combinations[b] = combinations[b], combinations[a]
	})
	return combinations
}

// aliasTable draws a variant with probability proportional to its weight in O(1) time: Walker's alias
// method in Vose's construction, kept in integers so the probabilities are exact. Column i is taken
// uniformly, then the variant i itself if a draw from [0, sum) is below threshold[i], otherwise alias[i].
type aliasTable struct {
	threshold []int64
	alias     []int
	sum       int64
}

// newAliasTable builds the table for positive weights in O(m) time for m variants.
func newAliasTable(weights []int) aliasTable {
	m := len(weights)
	table := aliasTable{threshold: make([]int64, m), alias: make([]int, m)}
	scaled := make([]int64, m)
	for _, weight := range weights {
		table.sum += int64(weight)
	}

	// every column holds sum/m of the probability mass, scaled by m·sum to stay in integers.
	var small, large []int
	for i, weight := range weights {
		scaled[i] = int64(weight) * int64(m)
		if scaled[i] < table.sum {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}
	for len(small) > 0 && len(large) > 0 {
		less, more := small[len(small)-1], large[len(large)-1]
		small = small[:len(small)-1]
		table.threshold[less], table.alias[less] = scaled[less], more
		// the rest of the column of less is filled from more.
		scaled[more] -= table.sum - scaled[less]
		if scaled[more] < table.sum {
			large = large[:len(large)-1]
			small = append(small, more)
		}
	}
	for _, i := range append(small, large...) {
		table.threshold[i], table.alias[i] = table.sum, i
	}
	return table
}

// draw returns a variant index from 0 to m-1.
func (table aliasTable) draw(rnd *rand.Rand) int {
	i := rnd.Intn(len(table.threshold))
	if rnd.Int63n(table.sum) < table.threshold[i] {
		return i
	}
	return table.alias[i]
}

// isUniform reports whether every layer's variants have equal weights, so that every combination is
// equally likely and Floyd's algorithm gives the same distribution in exactly amount draws.
func isUniform(weights [][]int) bool {
	for _, layer := range weights {
		for _, weight := range layer {
			if weight != layer[0] {
				return false
			}
		}
	}
	return true
}

// sampleWeighted returns amount distinct combination numbers whose layer variants are drawn independently
// by their weights, so the probability of a combination is the product of its variants' probabilities.
// A combination drawn before is drawn again; it gives up after maxDraws draws and also returns the number
// of draws made.
func sampleWeighted(rnd *rand.Rand, amount int, weights [][]int, maxDraws int) ([]*big.Int, int, bool) {
	tables := make([]aliasTable, len(weights))
	bases := make([]int, len(weights))
	for k, layer := range weights {
		tables[k] = newAliasTable(layer)
		bases[k] = len(layer)
	}

	// numbers that fit in int64 are compared as machine integers, larger ones by their bytes.
	small := combinationCount(bases).IsInt64()
	chosen := make(map[interface{}]struct{}, amount)
	combinations := make([]*big.Int, 0, amount)
	indices := make([]int, len(weights))
	draws := 0
	for len(combinations) < amount {
		if draws == maxDraws {
			return combinations, draws, false
		}
		draws++
		for k, table := range tables {
			indices[k] = table.draw(rnd) + 1
		}

		var key interface{}
		var combination *big.Int
		if small {
			number := combinationNumber64(indices, bases)
			key = number
			if _, ok := chosen[key]; ok {
				continue
			}
			combination = big.NewInt(number)
		} else {
			combination = combinationNumber(indices, bases)
			key = string(combination.Bytes())
			if _, ok := chosen[key]; ok {
				continue
			}
		}
		chosen[key] = struct{}{}
		combinations = append(combinations, combination)
	}
	return combinations, draws, true
}
