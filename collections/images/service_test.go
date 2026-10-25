package images

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLessByNumber(t *testing.T) {
	names := []string{"10.png", "2.png", "0.png", "100.png", "11.png", "1.png", "9.png", "25843209876543210987654321.png", "readme"}
	sort.SliceStable(names, func(a, b int) bool { return lessByNumber(names[a], names[b]) })
	require.Equal(t, []string{"0.png", "1.png", "2.png", "9.png", "10.png", "11.png", "100.png", "25843209876543210987654321.png", "readme"}, names)
}
