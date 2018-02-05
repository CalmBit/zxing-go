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
	"errors"
	"math"
)

type BitMatrix struct {
	width, height, rowSize uint32
	bits                   [][]uint32
}

func NewBitMatrixFromDimension(dimension uint32) (*BitMatrix, error) {
	return NewBitMatrix(dimension, dimension)
}

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

func Parse(image [][]bool) (*BitMatrix, error) {
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

func (bm *BitMatrix) Get(x, y uint32) bool {
	return ((bm.bits[y][(x/32)] >> (x & 0x1F)) & 1) != 0
}

func (bm *BitMatrix) Set(x, y uint32) {
	bm.bits[y][(x / 32)] |= (1 << (x & 0x1F))
}

func (bm *BitMatrix) Unset(x, y uint32) {
	bm.bits[y][(x / 32)] &= ^(1 << (x & 0x1F))
}

func (bm *BitMatrix) Flip(x, y uint32) {
	bm.bits[y][(x / 32)] ^= (1 << (x & 0x1F))
}

func (bm *BitMatrix) Xor(mask BitMatrix) error {
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

func (bm *BitMatrix) Clear() {
	width := bm.rowSize
	height := bm.height

	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			bm.bits[y][x] = 0
		}
	}
}

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
			bm.bits[y][x] |= (1 << (x & 0x1F))
		}
	}
	return nil
}

func (bm *BitMatrix) GetRow(y uint32, row *BitArray) *BitArray {
	if row.size < bm.width {
		row = NewBitArray(bm.width)
	} else {
		row.Clear()
	}

	for x := uint32(0); x < bm.rowSize; x++ {
		row.SetBulk(x*32, bm.bits[y][x])
	}

	return row
}

func (bm *BitMatrix) SetRow(y int, row BitArray) {
	max := uint32(math.Min(float64(bm.rowSize), float64(len(row.bits))))

	for x := uint32(0); x < max; x++ {
		bm.bits[y][x] = row.bits[x]
	}
}
