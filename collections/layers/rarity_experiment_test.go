package layers

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The experiments of the thesis on choosing variants by their weights. Like the others they run only with
// GENERATOR_EXPERIMENTS; the weights are those of the procedural art set (scripts/demo-art/render.js writes
// data/demo-art/manifest.json, GENERATOR_ART_MANIFEST overrides the path).

type artLayer struct {
	Name     string `json:"name"`
	None     int    `json:"none"`
	Variants []struct {
		File   string `json:"file"`
		Weight int    `json:"weight"`
	} `json:"variants"`
}

func artLayers(t *testing.T) []Layer {
	path := os.Getenv("GENERATOR_ART_MANIFEST")
	if path == "" {
		path = filepath.Join("..", "..", "data", "demo-art", "manifest.json")
	}
	data, err := os.ReadFile(path)
	require.NoError(t, err, "render the art set first: node scripts/demo-art/render.js")

	var manifest []artLayer
	require.NoError(t, json.Unmarshal(data, &manifest))

	projectLayers := make([]Layer, len(manifest))
	for k, layer := range manifest {
		projectLayers[k] = Layer{Name: layer.Name, NoneWeight: layer.None, Weights: map[string]int{}}
		for _, variant := range layer.Variants {
			projectLayers[k].ImageNames = append(projectLayers[k].ImageNames, variant.File)
			projectLayers[k].Weights[variant.File] = variant.Weight
		}
		projectLayers[k].LayerImagesCount = len(layer.Variants)
	}
	return projectLayers
}

// regularizedGammaQ returns Q(a, x) = Gamma(a, x) / Gamma(a), the upper tail of the gamma distribution, by the
// series for x < a + 1 and by Lentz's continued fraction otherwise (Numerical Recipes, 6.2).
func regularizedGammaQ(a, x float64) float64 {
	if x <= 0 {
		return 1
	}
	lgamma, _ := math.Lgamma(a)
	if x < a+1 {
		sum, term := 1/a, 1/a
		for n := 1; n < 1000; n++ {
			term *= x / (a + float64(n))
			sum += term
			if math.Abs(term) < math.Abs(sum)*1e-15 {
				break
			}
		}
		return 1 - sum*math.Exp(-x+a*math.Log(x)-lgamma)
	}
	const tiny = 1e-300
	b := x + 1 - a
	c, d := 1/tiny, 1/b
	h := d
	for n := 1; n < 1000; n++ {
		an := -float64(n) * (float64(n) - a)
		b += 2
		d = an*d + b
		if math.Abs(d) < tiny {
			d = tiny
		}
		c = b + an/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		d = 1 / d
		delta := d * c
		h *= delta
		if math.Abs(delta-1) < 1e-15 {
			break
		}
	}
	return math.Exp(-x+a*math.Log(x)-lgamma) * h
}

func TestRegularizedGammaQ(t *testing.T) {
	// the chi-square tail: P(chi2(1) > 3.841) = 0.05, P(chi2(10) > 18.307) = 0.05, P(chi2(2) > x) = e^(-x/2).
	require.InDelta(t, 0.05, regularizedGammaQ(0.5, 3.841/2), 1e-4)
	require.InDelta(t, 0.05, regularizedGammaQ(5, 18.307/2), 1e-4)
	require.InDelta(t, math.Exp(-3), regularizedGammaQ(1, 3), 1e-12)
}

// variantCounts counts how many of the combinations take every variant of every layer.
func variantCounts(combinations []*big.Int, bases []int) [][]int {
	counts := make([][]int, len(bases))
	for k, base := range bases {
		counts[k] = make([]int, base)
	}
	for _, combination := range combinations {
		for k, index := range layerIndices(combination, bases) {
			counts[k][index-1]++
		}
	}
	return counts
}

// chiSquare returns Pearson's statistic sum((O - E)^2 / E) of the counts against the weights.
func chiSquare(counts, weights []int) float64 {
	total, sum := 0, 0
	for i := range counts {
		total += counts[i]
		sum += weights[i]
	}
	statistic := 0.0
	for i, count := range counts {
		expected := float64(total) * float64(weights[i]) / float64(sum)
		statistic += (float64(count) - expected) * (float64(count) - expected) / expected
	}
	return statistic
}

// TestExperimentFrequencies checks that the variant frequencies of a generated collection follow the
// weights: the chi-square goodness-of-fit test for every layer of the art set, N = 1000 and 10000, 5 runs each.
// Under the null hypothesis about 5% of the tests fall below p = 0.05 by chance.
func TestExperimentFrequencies(t *testing.T) {
	dir := experimentDir(t)
	projectLayers := artLayers(t)
	bases, weights := layerBases(projectLayers), layerWeights(projectLayers)

	var rows, detail [][]string
	for _, amount := range []int{1000, 10000} {
		for repeat := 1; repeat <= experimentRepeats; repeat++ {
			rnd := rand.New(rand.NewSource(int64(amount + repeat)))
			// the art set is not uniform, so chooseCombinations draws by the weights exactly like this.
			require.False(t, isUniform(weights))
			start := time.Now()
			combinations, draws, ok := sampleWeighted(rnd, amount, weights, maxDrawsPerImage*amount+1000)
			require.True(t, ok)
			elapsed := time.Since(start)

			rejected := 0
			for k, counts := range variantCounts(combinations, bases) {
				statistic := chiSquare(counts, weights[k])
				dof := len(counts) - 1
				p := regularizedGammaQ(float64(dof)/2, statistic/2)
				if p < 0.05 {
					rejected++
				}
				rows = append(rows, []string{
					strconv.Itoa(amount), strconv.Itoa(repeat), projectLayers[k].Name, strconv.Itoa(len(counts)),
					strconv.FormatFloat(statistic, 'f', 3, 64), strconv.Itoa(dof), strconv.FormatFloat(p, 'f', 4, 64),
					strconv.FormatFloat(float64(elapsed.Microseconds())/1000, 'f', 2, 64), strconv.Itoa(draws),
				})

				if amount == 10000 && repeat == 1 {
					sum := 0
					for _, weight := range weights[k] {
						sum += weight
					}
					for i, count := range counts {
						name := NoneName
						if i < len(projectLayers[k].ImageNames) {
							name = projectLayers[k].ImageNames[i]
						}
						detail = append(detail, []string{
							projectLayers[k].Name, name, strconv.Itoa(weights[k][i]),
							strconv.FormatFloat(float64(weights[k][i])/float64(sum), 'f', 5, 64),
							strconv.FormatFloat(float64(count)/float64(amount), 'f', 5, 64), strconv.Itoa(count),
						})
					}
				}
			}
			t.Logf("N=%d #%d: chosen in %v with %d draws, %d of %d layers with p < 0.05", amount, repeat, elapsed, draws, rejected, len(bases))
		}
	}

	writeCSV(t, filepath.Join(dir, "frequencies.csv"),
		[]string{"amount", "repeat", "layer", "variants", "chi2", "dof", "p_value", "choice_ms", "draws"}, rows)
	writeCSV(t, filepath.Join(dir, "frequencies-detail.csv"),
		[]string{"layer", "variant", "weight", "expected_share", "observed_share", "count"}, detail)
}

// zipfWeights returns weights proportional to 1/i^s for m variants (s = 0 - equal weights), scaled to
// integers with the rarest one at least 1.
func zipfWeights(m int, s float64) []int {
	weights := make([]int, m)
	for i := range weights {
		weights[i] = int(math.Round(1000 / math.Pow(float64(i+1), s)))
		if weights[i] < 1 {
			weights[i] = 1
		}
	}
	return weights
}

func sumOf(values []int) int {
	sum := 0
	for _, value := range values {
		sum += value
	}
	return sum
}

// totalVariation returns the total variation distance 1/2 * sum(|observed - expected|) between the variant shares.
func totalVariation(counts, weights []int) float64 {
	total, sum := 0, 0
	for i := range counts {
		total += counts[i]
		sum += weights[i]
	}
	distance := 0.0
	for i := range counts {
		distance += math.Abs(float64(counts[i])/float64(total) - float64(weights[i])/float64(sum))
	}
	return distance / 2
}

// TestExperimentUniqueness measures the price of uniqueness under weights: the draws needed to choose
// N distinct combinations out of C = 6^4 = 1296 and the bias of the variant frequencies away from the
// weights as N/C grows, for equal weights and Zipf weights 1/i^s with s = 0.5, 1 and 2.
func TestExperimentUniqueness(t *testing.T) {
	dir := experimentDir(t)
	const layersCount, variants = 4, 6
	total := 1
	for k := 0; k < layersCount; k++ {
		total *= variants
	}

	var rows [][]string
	for _, s := range []float64{0, 0.5, 1, 2} {
		weights := make([][]int, layersCount)
		for k := range weights {
			weights[k] = zipfWeights(variants, s)
		}
		bases := make([]int, layersCount)
		for k := range bases {
			bases[k] = variants
		}

		for _, ratio := range []float64{0.05, 0.1, 0.25, 0.5, 0.75, 0.9, 1} {
			amount := int(math.Round(ratio * float64(total)))
			for repeat := 1; repeat <= experimentRepeats; repeat++ {
				rnd := rand.New(rand.NewSource(int64(repeat)))
				maxDraws := maxDrawsPerImage*amount + 1000
				combinations, draws, ok := sampleWeighted(rnd, amount, weights, maxDraws)

				// the layers have equal weights, so their variant counts are pooled (4N choices) to cut the
				// sampling noise; the frequencies are those of what was chosen, even when the draws ran out.
				pooled := make([]int, variants)
				for _, counts := range variantCounts(combinations, bases) {
					for i, count := range counts {
						pooled[i] += count
					}
				}
				distance := totalVariation(pooled, weights[0])
				topShare := float64(pooled[0]) / float64(layersCount*len(combinations))
				expectedTop := float64(weights[0][0]) / float64(sumOf(weights[0]))
				rows = append(rows, []string{
					strconv.FormatFloat(s, 'f', 1, 64), strconv.Itoa(total), strconv.FormatFloat(ratio, 'f', 2, 64),
					strconv.Itoa(amount), strconv.Itoa(repeat), strconv.Itoa(draws), strconv.Itoa(len(combinations)),
					strconv.FormatBool(ok), strconv.FormatFloat(distance, 'f', 4, 64),
					strconv.FormatFloat(topShare, 'f', 4, 64), strconv.FormatFloat(expectedTop, 'f', 4, 64),
				})
				t.Logf("s=%.1f N/C=%.2f N=%d #%d: %d draws (%.2f per image), %v, TV=%.4f, top %.3f of %.3f",
					s, ratio, amount, repeat, draws, float64(draws)/float64(amount), ok, distance, topShare, expectedTop)
			}
		}
	}

	writeCSV(t, filepath.Join(dir, "uniqueness.csv"),
		[]string{"zipf_s", "combinations", "ratio", "amount", "repeat", "draws", "chosen", "ok", "tv_distance", "top_share", "expected_top_share"}, rows)
}

// TestExperimentArtSet prints the size of the art set for the thesis: layers, variants and C.
func TestExperimentArtSet(t *testing.T) {
	experimentDir(t)
	projectLayers := artLayers(t)
	variants := 0
	for _, layer := range projectLayers {
		variants += layer.LayerImagesCount
	}
	total := combinationCount(layerBases(projectLayers))
	t.Logf("%d layers, %d variants, C = %s (%s)", len(projectLayers), variants, total, fmt.Sprintf("%.3g", new(big.Float).SetInt(total)))
}
