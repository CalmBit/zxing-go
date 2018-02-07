/*
 * Copyright 2007 ZXing authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"bytes"
	"errors"

	"github.com/discesoft/zxing-go/core/internal"
)

// BitMatrix encapsulates a matrix of single bits into a compact, space-efficent
// format. An interface is provided to interact with individual bits and
// small or large swathes of bits of varying dimensions.
type BitMatrix struct {
	width, height, rowSize uint32
	bits                   [][]uint32
}

// NewBitMatrixFromDimension calls NewBitMatrix with both the height value
// and width value set to the value of dimension. This produces a perfectly
// square BitMatrix.
// It returns a pointer to a new BitMatrix on success, and an error if either
// of the dimensions are below zero.
func NewBitMatrixFromDimension(dimension uint32) (*BitMatrix, error) {
	return NewBitMatrix(dimension, dimension)
}

// NewBitMatrix constructs a BitMatrix according to the given width and
// height.
// It returns a pointer to a new BitMatrix on success, and an error if either
// of the dimensions are below zero.
func NewBitMatrix(width, height uint32) (*BitMatrix, error) {
	if (width == 0) || (height == 0) {
		return nil, errors.New("Both dimensions must be greater than 0")
	}

	rowSize := ((width + 31) / 32)
	bits := make([][]uint32, height)
	for i := uint32(0); i < height; i++ {
		bits[i] = make([]uint32, rowSize)
	}
	return &BitMatrix{width, height, rowSize, bits}, nil
}

// ParseToBitMatrix takes a 2D slice of booleans and attempts to parse it
// into the form of a BitMatrix. image is the aformentioned slice, which
// will have an attempt made to parse it into a sensible format.
// It returns a pointer to a new BitMatrix on success, or an error on
// some failure.
func ParseToBitMatrix(image [][]bool) (*BitMatrix, error) {
	height := uint32(len(image))
	width := uint32(len(image[0]))

	bits, err := NewBitMatrix(width, height)

	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < height; i++ {
		imageI := image[i]
		for j := uint32(0); j < width; j++ {
			if imageI[j] {
				bits.Set(j, i)
			}
		}
	}
	return bits, nil
}

// ParseStringToBitMatrix takes an arbitrary stringRepresentation
// and makes an attempt to parse the information into a sensible
// BitMatrix value. setString is the value of a "1" within
// stringRepresentation, while unsetString is the value of a "0".
// It returns a pointer to a new BitMatrix on success, or an error
// on some failure.
func ParseStringToBitMatrix(stringRepresentation, setString, unsetString string) (*BitMatrix, error) {
	bits := make([]bool, len(stringRepresentation))

	var bitsPos, rowStartPos, rowLength, nRows, pos int
	rowLength = -1

	runes := []rune(stringRepresentation)
	setRunes := []rune(setString)
	unsetRunes := []rune(unsetString)
	for pos < len(runes) {
		if runes[pos] == '\n' || runes[pos] == '\r' {
			if bitsPos > rowStartPos {
				if rowLength == -1 {
					rowLength = bitsPos - rowStartPos
				} else if (bitsPos - rowStartPos) != rowLength {
					return nil, errors.New("row lengths do not match")
				}
				rowStartPos = bitsPos
				nRows++
			}
			pos++
		} else if internal.CompareRuneSlices(runes[pos:pos+len(setRunes)], setRunes) {
			pos += len(setRunes)
			bits[bitsPos] = true
			bitsPos++
		} else if internal.CompareRuneSlices(runes[pos:pos+len(unsetRunes)], unsetRunes) {
			pos += len(unsetRunes)
			bits[bitsPos] = false
			bitsPos++
		} else {
			return nil, errors.New("illegal character encountered: " + string(runes[pos]))
		}
	}

	if bitsPos > rowStartPos {
		if rowLength == -1 {
			rowLength = bitsPos - rowStartPos
		} else if (bitsPos - rowStartPos) != rowLength {
			return nil, errors.New("row lengths do not match")
		}
		nRows++
	}

	matrix, err := NewBitMatrix(uint32(rowLength), uint32(nRows))

	if err != nil {
		return nil, err
	}

	for i := 0; i < bitsPos; i++ {
		if bits[i] {
			matrix.Set(uint32(i%rowLength), uint32(i/rowLength))
		}
	}
	return matrix, nil
}

// Get returns the value of a specified bit.
// x and y are the horizontal and vertical components of a
// specific location within the BitMatrix, respectively.
// It returns the value of the bit at {x,y}.
func (bm *BitMatrix) Get(x, y uint32) bool {
	return ((bm.bits[y][(x/32)] >> (x & 0x1F)) & 1) != 0
}

// Set activates a specified bit.
// x and y are the horizontal and vertical components of a
// specific location within the BitMatrix, respectively.
func (bm *BitMatrix) Set(x, y uint32) {
	bm.bits[y][(x / 32)] |= (1 << (x & 0x1F))
}

// Unset zeroes out the value of a specified bit.
// x and y are the horizontal and vertical components of a
// specific location within the BitMatrix, respectively.
func (bm *BitMatrix) Unset(x, y uint32) {
	bm.bits[y][(x / 32)] &= ^(1 << (x & 0x1F))
}

// Flip toggles the value of a specified bit.
// x and y are the horizontal and vertical components of a
// specific location within the BitMatrix, respectively.
func (bm *BitMatrix) Flip(x, y uint32) {
	bm.bits[y][(x / 32)] ^= (1 << (x & 0x1F))
}

// Xor functions as an exlusive-or operation against the BitMatrix -
// it expects that mask (another BitMatrix) will be the exact same
// dimensions as the calling BitMatrix.
// If a bit within the mask is set, the same bit within the calling
// BitMatrix will be flipped.
// It returns an error if the mask and calling BitMatrix dimensions are
// mismatched.
func (bm *BitMatrix) Xor(mask *BitMatrix) error {
	if (bm.width != mask.width) || (bm.height != mask.height) || (bm.rowSize != mask.rowSize) {
		return errors.New("input matrix dimensions do not match")
	}

	rowArray := NewBitArray((bm.width / 32) + 1)
	for y := uint32(0); y < bm.height; y++ {
		row := mask.GetRow(y, rowArray).GetBitArray()
		for x := uint32(0); x < bm.rowSize; x++ {
			bm.bits[y][x] ^= row[x]
		}
	}

	return nil
}

// Clear erases all data from the BitMatrix, setting every value to false.
func (bm *BitMatrix) Clear() {
	width := bm.rowSize
	height := bm.height

	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			bm.bits[y][x] = 0
		}
	}
}

// SetRegion takes an entire rectangular region of the BitMatrix and sets
// all values within that region to be true.
// left and top are both inclusive beginning points, while width and height
// are self-explanatory.
func (bm *BitMatrix) SetRegion(left, top, width, height uint32) error {
	if (height == 0) || (width == 0) {
		return errors.New("Height and width must be at least 1")
	}

	right := left + width
	bottom := top + height
	if (bottom > bm.height) || (right > bm.width) {
		return errors.New("The region must fit inside the matrix")
	}
	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			bm.bits[y][x/32] |= (1 << (x & 0x1F))
		}
	}
	return nil
}

// GetRow retrieves a single row, y, from the BitMatrix as a BitArray.
// row is a prallocated bit array, which will be reallocated if its found
// to be too small.
// It returns a pointer to a BitArray, which should be used as opposed to
// row, given the reallocation semantics herein.
func (bm *BitMatrix) GetRow(y uint32, row *BitArray) *BitArray {
	if row == nil || row.size < bm.width {
		row = NewBitArray(bm.width)
	} else {
		row.Clear()
	}

	for x := uint32(0); x < bm.rowSize; x++ {
		row.SetBulk(x*32, bm.bits[y][x])
	}

	return row
}

// SetRow copies data from the parameter `row`, and uses it
// to seed data into the BitMatrix row y.
// It returns an error if there's a problem with this transfer,
// or `nil` if everything went correctly.
func (bm *BitMatrix) SetRow(y uint32, row *BitArray) error {
	return internal.SliceCopyU32(row.bits, bm.bits[y], 0, 0, int(bm.rowSize))
}

// Rotate180 rotates the current BitMatrix 180 degrees.
// It returns an error if there's a problem with the rotation,
// or `nil` if everything went correctly.
func (bm *BitMatrix) Rotate180() error {
	width := bm.width
	height := bm.height
	topRow := NewBitArray(width)
	bottomRow := NewBitArray(width)

	for i := uint32(0); i < (height+1)/2; i++ {
		topRow = bm.GetRow(i, topRow)
		bottomRow = bm.GetRow(height-1-i, bottomRow)
		topRow.Reverse()
		bottomRow.Reverse()

		err := bm.SetRow(i, bottomRow)
		if err != nil {
			return err
		}

		err = bm.SetRow(height-1-i, topRow)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetEnclosingRectangle calculates the set bounds of a BitMatrix.
// It returns a four-element []uint32 formatted as such:
//     []uint32{left, top, width, height}
// If the BitMatrix is completely unset, or white, the result will
// be `nil`.
func (bm *BitMatrix) GetEnclosingRectangle() []uint32 {
	left := int(bm.width)
	top := int(bm.height)
	right := int(-1)
	bottom := int(-1)

	for y := 0; y < int(bm.height); y++ {
		for x32 := 0; x32 < int(bm.rowSize); x32++ {
			theBits := bm.bits[y][x32]
			if theBits != 0 {
				if y < top {
					top = y
				}
				if y > bottom {
					bottom = y
				}
				if (x32 * 32) < left {
					bit := 0
					for (theBits << uint(31-bit)) == 0 {
						bit++
					}
					if ((x32 * 32) + bit) < left {
						left = (x32 * 32) + bit
					}
				}
				if ((x32 * 32) + 31) > right {
					bit := 31
					for (theBits >> uint(bit)) == 0 {
						bit--
					}
					if ((x32 * 32) + bit) > right {
						right = (x32 * 32) + bit
					}
				}
			}
		}
	}

	if (right < left) || (bottom < top) {
		return nil
	}

	return []uint32{uint32(left), uint32(top), uint32(right - left + 1), uint32(bottom - top + 1)}
}

// GetTopLeftOnBit obtains the upper-left corner set bit, i.e.
// the upper left active corner of a BitMatrix.
// It returns a two-element []uint32 formatted as such:
//     []uint32{x,y}
// If the BitMatrix is completely unset, or white, the result will
// be `nil`.
func (bm *BitMatrix) GetTopLeftOnBit() []uint32 {
	yOffset := uint32(0)
	xOffset := uint32(0)
	for yOffset < bm.height-1 {
		for (xOffset < bm.rowSize-1) && (bm.bits[yOffset][xOffset] == 0) {
			xOffset++
		}
		if bm.bits[yOffset][xOffset] == 0 {
			yOffset++
		} else {
			break
		}
	}

	if (yOffset == bm.height-1) && (xOffset == bm.rowSize-1) {
		return nil
	}

	xBitOffset := xOffset * 32

	theBits := bm.bits[yOffset][xOffset]
	bit := uint32(0)

	for theBits<<(31-bit) == 0 {
		bit++
	}

	xBitOffset += bit

	return []uint32{xBitOffset, yOffset}
}

// GetBottomRightOnBit obtains the bottom-right corner set bit, i.e.
// the bottom-right active corner of a BitMatrix.
// It returns a two-element []uint32 formatted as such:
//     []uint32{x,y}
// If the BitMatrix is completely unset, or white, the result will
// be `nil`.
func (bm *BitMatrix) GetBottomRightOnBit() []uint32 {
	yOffset := bm.height - 1
	xOffset := bm.rowSize - 1
	for yOffset >= 0 {
		for (xOffset >= 0) && (bm.bits[yOffset][xOffset] == 0) {
			if xOffset == 0 {
				break
			}
			xOffset--
		}
		if bm.bits[yOffset][xOffset] == 0 {
			if yOffset == 0 {
				break
			}
			yOffset--
		} else {
			break
		}
	}

	if (yOffset == 0) && (xOffset == 0) {
		return nil
	}

	xBitOffset := xOffset * 32

	theBits := bm.bits[yOffset][xOffset]
	bit := uint32(31)

	for (theBits >> bit) == 0 {
		bit--
	}

	xBitOffset += bit

	return []uint32{xBitOffset, yOffset}
}

func (bm *BitMatrix) GetWidth() uint32 {
	return bm.width
}

func (bm *BitMatrix) GetHeight() uint32 {
	return bm.height
}

func (bm *BitMatrix) GetRowSize() uint32 {
	return bm.rowSize
}

func (bm *BitMatrix) Equals(other *BitMatrix) bool {
	if bm.width == other.width && bm.height == other.height && bm.rowSize == other.rowSize {
		for y := uint32(0); y < bm.height; y++ {
			for x := uint32(0); x < bm.rowSize; x++ {
				if bm.bits[y][x] != other.bits[y][x] {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (bm *BitMatrix) String() string {
	return bm.ToString("X ", "  ")
}

func (bm *BitMatrix) ToString(setString, unsetString string) string {
	return bm.buildToString(setString, unsetString, "\n")
}

func (bm *BitMatrix) ToStringWithSeperator(setString, unsetString, lineSeperator string) string {
	return bm.buildToString(setString, unsetString, lineSeperator)
}

func (bm *BitMatrix) buildToString(setString, unsetString, lineSeperator string) string {
	var result bytes.Buffer

	for y := uint32(0); y < bm.height; y++ {
		for x := uint32(0); x < bm.width; x++ {
			if bm.Get(x, y) {
				result.WriteString(setString)
			} else {
				result.WriteString(unsetString)
			}
		}
		result.WriteString(lineSeperator)
	}
	return result.String()
}

func (bm *BitMatrix) Clone() *BitMatrix {
	bits := make([][]uint32, len(bm.bits))

	for i := range bm.bits {
		bits[i] = internal.SliceCloneU32(bm.bits[i])
	}

	return &BitMatrix{bm.width, bm.height, bm.rowSize, bits}
}
