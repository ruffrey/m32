package m32

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func Test_Transpose(t *testing.T) {
	t.Run("it transposes a matrix", func(t *testing.T) {
		original := NewM(3, 3)
		original.W = []float32{
			1, 2, 3,
			4, 5, 6,
			7, 8, 9,
		}
		transposed := original.T()
		fmt.Println("original:", original.W)
		assert.EqualValues(t, []float32{
			1, 4, 7,
			2, 5, 8,
			3, 6, 9,
		}, transposed.W)
	})
}

func Test_Hstack(t *testing.T) {
	t.Run("it concatenates properly", func(t *testing.T) {
		m1 := NewM(5, 2)
		m1.W = []float32{
			1, 2, 3, 4, 5,
			6, 7, 8, 9, 10,
		}
		m2 := NewM(5, 2)
		m2.W = []float32{
			11, 12, 13, 14, 15,
			16, 17, 18, 19, 20,
		}
		stacked := HStack(m1, m2)
		assert.Equal(t, stacked.Cols, 10)
		assert.Equal(t, stacked.Rows, 2)
		assert.EqualValues(t, []float32{
			1, 2, 3, 4, 5, 11, 12, 13, 14, 15,
			6, 7, 8, 9, 10, 16, 17, 18, 19, 20,
		}, stacked.W)
	})
}

func Test_Add(t *testing.T) {
	t.Run("adds with m1 1 row", func(t *testing.T) {
		m1 := NewM(5, 2)
		m1.W = []float32{
			1, 2, 3, 4, 5,
			6, 7, 8, 9, 10,
		}
		m2 := NewM(5, 1)
		m2.W = []float32{
			1, 2, 3, 4, 5,
		}
		sum := Add(m1, m2)
		assert.Equal(t, sum.Cols, 5)
		assert.Equal(t, sum.Rows, 2)
		assert.EqualValues(t, []float32{
			2, 4, 6, 8, 10,
			7, 9, 11, 13, 15,
		}, sum.W)
	})
}
