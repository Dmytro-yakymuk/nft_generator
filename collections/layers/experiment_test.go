package layers

import (
	"encoding/csv"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The experiments of the thesis (section 3) on choosing and converting combination numbers.
// They are skipped unless GENERATOR_EXPERIMENTS names the directory for the CSV results:
//
//	GENERATOR_EXPERIMENTS=$PWD/experiments/results go test -run Experiment -v -timeout 1h ./collections/layers/
const experimentRepeats = 5

func experimentDir(t *testing.T) string {
	dir := os.Getenv("GENERATOR_EXPERIMENTS")
	if dir == "" {
		t.Skip("set GENERATOR_EXPERIMENTS to the results directory to run the experiment")
	}
	require.NoError(t, os.MkdirAll(dir, 0o755))
	return dir
}

func writeCSV(t *testing.T, path string, header []string, rows [][]string) {
	file, err := os.Create(path)
	require.NoError(t, err)
	defer func() { require.NoError(t, file.Close()) }()

	writer := csv.NewWriter(file)
	require.NoError(t, writer.Write(header))
	require.NoError(t, writer.WriteAll(rows))
	t.Logf("%s: %d rows", path, len(rows))
}

// incrementSamples keeps a repeat of the step-by-step conversion near 4e9 steps (~0.75 ns a step),
// with 3 to 10000 numbers.
func incrementSamples(total int) int {
	samples := 4_000_000_000 / total
	if samples < 3 {
		return 3
	}
	if samples > 10000 {
		return 10000
	}
	return samples
}

// conversionSink keeps the compiler from dropping the measured calls.
var conversionSink interface{}

// TestExperimentConversion compares the step-by-step conversion of a combination number, O(C), with
// the direct one, O(n), for n = 2...10 layers of 10 images each, i.e. C = 10^n combinations.
func TestExperimentConversion(t *testing.T) {
	dir := experimentDir(t)

	var rows [][]string
	for layersCount := 2; layersCount <= 10; layersCount++ {
		counts := make([]int, layersCount)
		total := 1
		for k := range counts {
			counts[k] = 10
			total *= 10
		}

		for _, method := range []struct {
			name    string
			convert func(int, []int) []int
			samples int
		}{
			{"increment", layerIndicesByIncrement, incrementSamples(total)},
			{"direct", func(number int, counts []int) []int { return layerIndices(big.NewInt(int64(number)), counts) }, 100000},
		} {
			// evenly spaced numbers: the mean work of the step-by-step conversion is exactly C/2.
			numbers := make([]int, method.samples)
			for j := range numbers {
				numbers[j] = int(float64(total)*(float64(j)+0.5)/float64(method.samples)) + 1
			}

			for repeat := 1; repeat <= experimentRepeats; repeat++ {
				runtime.GC()
				start := time.Now()
				for _, number := range numbers {
					conversionSink = method.convert(number, counts)
				}
				perConversion := float64(time.Since(start).Nanoseconds()) / float64(len(numbers))

				rows = append(rows, []string{
					strconv.Itoa(layersCount), "10", strconv.Itoa(total), method.name, strconv.Itoa(repeat),
					strconv.Itoa(len(numbers)), strconv.FormatFloat(perConversion, 'f', 1, 64),
				})
				t.Logf("n=%d C=%d %s #%d: %.1f ns", layersCount, total, method.name, repeat, perConversion)
			}
		}
	}

	writeCSV(t, filepath.Join(dir, "conversion.csv"),
		[]string{"layers", "images_per_layer", "combinations", "method", "repeat", "samples", "ns_per_conversion"}, rows)
}

// sampleByReseed is the former choice of combinations: the global source was reseeded with the current
// second before every image, so within a second the sequence started over and the i-th image needed
// about i draws; then draws were repeated until an unused number turned up. It gives up at the deadline.
func sampleByReseed(amount, total int, deadline time.Time) (draws int, ok bool) {
	used := make(map[int]struct{}, amount)
	for i := 1; i <= amount; i++ {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			draws++
			number := rnd.Intn(total) + 1
			if _, ok := used[number]; !ok {
				used[number] = struct{}{}
				break
			}
			if draws%1_000_000 == 0 && time.Now().After(deadline) {
				return draws, false
			}
		}
	}
	return draws, true
}

// sampleByRejection draws until an unused number turns up, with the source seeded once.
func sampleByRejection(rnd *rand.Rand, amount, total int) (draws int) {
	used := make(map[int]struct{}, amount)
	for len(used) < amount {
		draws++
		used[rnd.Intn(total)+1] = struct{}{}
	}
	return draws
}

// TestExperimentSampling compares the choice of amount distinct combinations out of C = 100000:
// the former reseeding, drawing until an unused number with one seed, and Floyd's algorithm.
func TestExperimentSampling(t *testing.T) {
	dir := experimentDir(t)
	const total, limit = 100000, time.Minute

	var rows [][]string
	timedOut := false
	for _, amount := range []int{1000, 5000, 10000, 25000, 50000, 75000, 90000, 100000} {
		for _, method := range []string{"reseed", "rejection", "floyd"} {
			if method == "reseed" && timedOut {
				rows = append(rows, []string{strconv.Itoa(total), strconv.Itoa(amount), method, "1", "", "", "skipped"})
				continue
			}

			for repeat := 1; repeat <= experimentRepeats; repeat++ {
				rnd := rand.New(rand.NewSource(int64(repeat)))
				runtime.GC()
				start := time.Now()
				draws, status := amount, "ok"
				switch method {
				case "reseed":
					var ok bool
					if draws, ok = sampleByReseed(amount, total, start.Add(limit)); !ok {
						status, timedOut = "timeout", true
					}
				case "rejection":
					draws = sampleByRejection(rnd, amount, total)
				case "floyd":
					conversionSink = sampleCombinations(rnd, amount, big.NewInt(total))
					// one draw per number in Floyd's loop and amount-1 more in rand.Shuffle of the result.
					draws = 2*amount - 1
				}
				elapsed := float64(time.Since(start).Microseconds()) / 1000

				rows = append(rows, []string{
					strconv.Itoa(total), strconv.Itoa(amount), method, strconv.Itoa(repeat),
					strconv.Itoa(draws), strconv.FormatFloat(elapsed, 'f', 3, 64), status,
				})
				t.Logf("C=%d N=%d %s #%d: %d draws, %.3f ms %s", total, amount, method, repeat, draws, elapsed, status)
				if status != "ok" {
					break
				}
			}
		}
	}

	writeCSV(t, filepath.Join(dir, "sampling.csv"),
		[]string{"combinations", "amount", "method", "repeat", "draws", "ms", "status"}, rows)
}
