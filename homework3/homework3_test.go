package homework3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	input          [][]int
	expectedOutput []int
}

var data = []testData{
	{
		[][]int{
			{1},
		},
		[]int{1},
	},
	{
		[][]int{
			{1, 2, 3},
		},
		[]int{1, 2, 3},
	},
	{
		[][]int{
			{1},
			{4},
			{7},
		},
		[]int{1, 4, 7},
	},
	{
		[][]int{
			{1, 2},
			{4, 5},
		},
		[]int{1, 2, 5, 4},
	},
	{
		[][]int{
			{1, 2},
			{4, 5},
			{7, 8},
		},
		[]int{1, 2, 5, 8, 7, 4},
	},
	{
		[][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		[]int{1, 2, 3, 6, 9, 8, 7, 4, 5},
	},
	{
		[][]int{
			{1, 2, 3},
			{4, 5, 6, 9},
			{7, 8, 9},
		},
		nil,
	},
	{
		[][]int{},
		nil,
	},
}

func TestMain(t *testing.T) {
	for _, d := range data {
		output, _ := FlattenMatrixClockwise(d.input)
		assert.Equal(t, d.expectedOutput, output, "input: "+fmt.Sprint(d.input))
	}
}
