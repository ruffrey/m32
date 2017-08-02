package m32

import (
	"fmt"
	"math"
	"math/rand"
)

/*
M holds a matrix.
*/
type M struct {
	Rows int
	Cols int
	W    []float32
}

/*
NewM makes a zeroed matrix
*/
func NewM(cols, rows int) (m *M) {
	m = &M{
		Rows: rows,
		Cols: cols,
		W:    make([]float32, rows*cols),
	}
	return m
}

/*
NewMLike make a zeroed matrix in the same shape as another matrix.
*/
func NewMLike(m *M) (out *M) {
	out = &M{
		Rows: m.Rows,
		Cols: m.Cols,
		W:    make([]float32, m.Rows*m.Cols),
	}
	return out
}

/*
RandM makes a matrix with random values within a range.
*/
func RandM(cols, rows int, low, high float32) (m *M) {
	m = &M{
		Rows: rows,
		Cols: cols,
	}
	m.W = RandArrayBetween(low, high, rows*cols)
	return m
}

func (m *M) String() string {
	return fmt.Sprintf("{ Rows: %d, Cols: %d }", m.Rows, m.Cols)
}

/*
Clone makes a copy of a matrix
*/
func (m *M) Clone() (out *M) {
	out = &M{
		Rows: m.Rows,
		Cols: m.Cols,
		W:    m.W[:],
	}
	return out
}

/*
RowAt returns a Matrix using row at the specified index.
*/
func (m *M) RowAt(index int) (out *M) {
	out = NewM(m.Cols, 1)
	rowStart := index * m.Cols
	nextRowIndex := (index + 1) * m.Cols
	out.W = m.W[rowStart:nextRowIndex]
	return out
}

/*
PlusEquals modifies the matrix m by substracting m2 in place.
*/
func (m *M) PlusEquals(m2 *M) {
	n := len(m.W)

	for i := 0; i < n; i++ {
		m.W[i] += m2.W[i]
	}
}

/*
T transposes the matrix:
returns a copied matrix that is the transposed version of the original.


	[1, 2          [1, 4, 7,
	 4, 5  becomes  2, 5, 8]
	 7, 8]

*/
func (m *M) T() (out *M) {
	out = &M{
		// swap
		Rows: m.Cols,
		Cols: m.Rows,
	}

	var row int
	var cursor int
	n_1 := len(m.W) - 1
	for row = 0; row < m.Rows; {
		// go down each column and put it in a flat row array
		out.W = append(out.W, m.W[cursor])

		cursor += m.Rows
		if cursor > n_1 {
			row++
			cursor = row
		}
	}

	return out
}

/**
Pow takes every value in the matrix to the `power` and returns it.
*/
func Pow(m *M, power float64) (out *M) {
	out = NewMLike(m)
	n := len(m.W)
	for i := 0; i < n; i++ {
		out.W[i] = float32(math.Pow(float64(m.W[i]), power))
	}
	return out
}

// Times multiplies each array value by the corresponding one and returns the result
func Times(m1, m2 *M) (out *M) {
	out = NewMLike(m1)
	lenM1 := len(m1.W)
	lenM2 := len(m2.W)
	for i := 0; i < lenM1 && i < lenM2; i++ {
		out.W[i] = m1.W[i] * m2.W[i]
	}
	return out
}

/*
Dot two arrays and return the result (dot product)
*/
func Dot(m1, m2 *M) (out *M) {
	// apparently it is not that big a deal in numpy when these are different sizes

	//if m1.Rows != m2.Cols {
	//	fmt.Println("m1=", m1)
	//	fmt.Println("m2=", m2)
	//	panic("Cannot multiply arrays of different sizes")
	//}
	n := m1.Rows
	d := m2.Cols
	out = NewM(n, d)
	m2Len := len(m2.W)
	var cellSum float32 = 0.0

	for row := 0; row < m1.Rows; row++ { // loop over rows of m1
		for col := 0; col < m2.Cols; col++ { // loop over cols of m2
			cellSum = 0.0
			for colCell := 0; colCell < m1.Cols; colCell++ { // dot product loop
				m2Col := m2.Cols*colCell + col
				if m2Col < m2Len {
					m1Row := m1.Cols*row + colCell
					cellSum += m1.W[m1Row] * m2.W[m2Col]
				}
			}
			out.W[d*row+col] = cellSum
		}
	}
	return out
}

/*
Outer produces the outer product of two arrays.

Note: not sure if technically correct. m2.Rows always =1
*/
func Outer(m1, m2 *M) (out *M) {
	out = NewM(m1.Cols, m1.Rows)

	ix := 0
	n := len(m1.W)
	colIter := 0
	m2Len := len(m2.W)
	for ; ix < n; ix++ {
		if colIter < m2Len { // prevent out of bounds
			out.W[ix] = m1.W[ix] * m2.W[colIter]
		} // else zero
		colIter++
		if colIter >= m1.Cols {
			colIter = 0
		}
	}

	return out
	//return Times(m1, m2)
}

/*
MultiplyValue multiplies the `value` by every item in the array and returns the result.
*/
func MultiplyValue(value float32, m1 *M) (out *M) {
	out = NewM(m1.Cols, m1.Rows)
	n := len(m1.W)
	for i := 0; i < n; i++ {
		out.W[i] = value * m1.W[i]
	}
	return out
}

/*
Add two matrices and return the result.
*/
func Add(m1, m2 *M) (out *M) {
	out = NewM(m1.Cols, m1.Rows)

	m2Len := len(m2.W)

	if m2.Rows == 1 { // multiply this row by all rows of the other array
		ix := 0
		n := len(m1.W)
		colIter := 0
		for ; ix < n; ix++ {
			if colIter < m2Len { // prevent out of bounds
				out.W[ix] = m1.W[ix] + m2.W[colIter]
			} else {
				out.W[ix] = m1.W[ix]
			}
			colIter++
			if colIter >= m1.Cols {
				colIter = 0
			}
		}

		return out
	}

	var cursor int
	var m2Cursor int
	for row := 0; row < m1.Rows; row++ {
		cursor = row * m1.Cols
		m2Cursor = row * m2.Cols
		for col := 0; col < m1.Cols; col++ {
			if m2Cursor < m2.Cols {
				out.W[cursor] = m1.W[cursor] + m2.W[m2Cursor]
			} else {
				out.W[cursor] = m1.W[cursor]
			}
			cursor++
			m2Cursor++
		}
	}

	return out
}

/*
Subtract m2 from m1 and return the result.
*/
func Subtract(m1, m2 *M) (out *M) {
	out = NewM(m1.Cols, m1.Rows)
	n := len(m1.W)
	// assuming all arrays are the same size
	for i := 0; i < n; i++ {
		out.W[i] = m1.W[i] - m2.W[i]
	}

	return out
}

/*
SubtractFromAll subtracts a value from every array item and returns the result.
*/
func SubtractFromAll(m *M, subtract float32) (out *M) {
	out = NewM(m.Cols, m.Rows)
	n := len(m.W)

	for i := 0; i < n; i++ {
		out.W[i] = m.W[i] - subtract
	}

	return out
}

/*
HStack horizontally stacks two arrays? concatenates them? see numpy hstack.
*/
func HStack(m1, m2 *M) (out *M) {
	if m1.Rows == 0 && m2.Rows == 0 {
		out = NewM(m1.Cols+m2.Cols, 1)
		out.W = append(m1.W, m2.W...)
		return out
	}
	out = NewM(m1.Cols+m2.Cols, m1.Rows)

	iter := 0
	for row := 0; row < out.Rows; row++ {
		iter = row * (m1.Cols * 2)
		for col := 0; col < m1.Cols; col++ {
			out.W[iter] = m1.W[iter-(m1.Cols*row)]
			out.W[iter+m1.Cols] = m2.W[iter-(m2.Cols*row)]
			iter++
		}
	}

	return out
}

/*
Sigmoid calculates the sigmoid for every element of an array.

The m[i] value is squashed to between 0 and 1.
*/
func Sigmoid(m *M) (out *M) {
	out = NewMLike(m)
	n := len(m.W)
	for i := 0; i < n; i++ {
		// sigmoid calculation
		sig := float32(1.0 / (1.0 + math.Exp(-float64(m.W[i]))))
		out.W[i] = sig
	}

	return out
}

/*
SigmoidDerivative calculates the derivative
*/
func SigmoidDerivative(m *M) (out *M) {
	out = NewMLike(m)
	n := len(m.W)

	for i := 0; i < n; i++ {
		out.W[i] = m.W[i] * (1 - m.W[i])
	}

	return out
}

/*
Tanh is the hyperbolic tangent for every element of the array.

The m.W[i] value is squashed to between -1 and 1.
*/
func Tanh(m *M) (out *M) {
	out = NewMLike(m)
	n := len(m.W)

	for i := 0; i < n; i++ {
		out.W[i] = float32(math.Tanh(float64(m.W[i])))
	}

	return out
}

var _one float32 = 1

/*
TanhDerivative calculates the derivative
*/
func TanhDerivative(m *M) (out *M) {
	out = NewMLike(m)
	n := len(m.W)

	for i := 0; i < n; i++ {
		out.W[i] = _one - float32(math.Pow(float64(m.W[i]), 2))
	}

	return out
}

/*
Softplus is an alternative to Softmax, which might be better or might not.

It certainly is faster.
*/
func Softplus(m *M) (out *M) {
	out = NewMLike(m)
	n := len(m.W)

	for i := 0; i < n; i++ {
		out.W[i] = float32(math.Max(0, float64(m.W[i])))
	}

	return out
}

// RandArray makes a random float32 array
func RandArray(length int) (f32 []float32) {
	f32 = make([]float32, length)
	for i := 0; i < length; i++ {
		f32[i] = rand.Float32()
	}
	return f32
}

// RandArrayBetween makes a random float32 array where values are in a range
func RandArrayBetween(low, high float32, length int) (f32 []float32) {
	for i := 0; i < length; i++ {
		f32 = append(f32, Randf(low, high))
	}
	return f32
}

/*
Randf makes random float32 in a range
*/
func Randf(low float32, high float32) float32 {
	return rand.Float32()*(high-low) + low
}
