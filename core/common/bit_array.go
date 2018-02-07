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
	"math/bits"

	"github.com/discesoft/zxing-go/core/internal"
)

type BitArray struct {
	bits []uint32
	size uint32
}

func NewEmptyBitArray() *BitArray {
	return &BitArray{make([]uint32, 1), 0}
}

func NewBitArray(size uint32) *BitArray {
	return &BitArray{make([]uint32, ((size + 31) / 32)), size}
}

func (ba *BitArray) GetSize() uint32 {
	return ba.size
}

func (ba *BitArray) GetSizeInBytes() uint32 {
	return (ba.size + 7) / 8
}

func (ba *BitArray) EnsureCapacity(size uint32) {
	if size > uint32(len(ba.bits)*32) {
		newBits := make([]uint32, size)
		for i := range ba.bits {
			newBits[i] = ba.bits[i]
		}
		ba.bits = newBits
	}
}

func (ba *BitArray) Get(i uint32) bool {
	return (ba.bits[(i/32)] & (1 << (i & 0x1F))) != 0
}

func (ba *BitArray) Set(i uint32) {
	ba.bits[(i / 32)] |= (1 << (i & 0x1F))
}

func (ba *BitArray) Flip(i uint32) {
	ba.bits[(i / 32)] ^= (1 << (i & 0x1F))
}

func (ba *BitArray) GetNextSet(from uint32) uint32 {
	if from >= ba.size {
		return ba.size
	}

	bitsOffset := (from / 32)
	currentBits := ba.bits[bitsOffset]
	currentBits &= ^((1 << (from & 0x1F)) - 1)
	for {
		if currentBits != 0 {
			break
		}

		bitsOffset++
		if bitsOffset == uint32(len(ba.bits)) {
			return ba.size
		}

		currentBits = ba.bits[bitsOffset]
	}

	result := (bitsOffset * 32) + uint32(bits.TrailingZeros32(currentBits))

	if result > ba.size {
		return ba.size
	}
	return result
}

func (ba *BitArray) GetNextUnset(from uint32) uint32 {
	if from >= ba.size {
		return ba.size
	}

	bitsOffset := (from / 32)
	currentBits := ^(ba.bits[bitsOffset])
	currentBits &= ^((1 << (from & 0x3F)) - 1)
	for {
		if currentBits != 0 {
			break
		}

		bitsOffset++
		if bitsOffset == uint32(len(ba.bits)) {
			return ba.size
		}

		currentBits = ^(ba.bits[bitsOffset])
	}

	result := (bitsOffset * 32) + uint32(bits.TrailingZeros32(currentBits))

	if result > ba.size {
		return ba.size
	}
	return result
}

func (ba *BitArray) SetBulk(i, newBits uint32) {
	ba.bits[(i / 32)] = newBits
}

func (ba *BitArray) SetRange(start, end uint32) error {
	if end < start || end > ba.size {
		return errors.New("Invalid argument passed to SetRange on BitArray")
	}

	if end == start {
		return nil
	}

	end--
	firstInt := start / 32
	lastInt := end / 32
	for i := firstInt; i <= lastInt; i++ {
		var firstBit, lastBit uint32

		if i > firstInt {
			firstBit = 0
		} else {
			firstBit = (start & 0x1F)
		}

		if i < lastInt {
			lastBit = 31
		} else {
			lastBit = (end & 0x1F)
		}

		mask := uint32((2 << lastBit) - (1 << firstBit))
		ba.bits[i] |= mask
	}

	return nil
}

func (ba *BitArray) Clear() {
	max := len(ba.bits)
	for i := 0; i < max; i++ {
		ba.bits[i] = 0
	}
}

func (ba *BitArray) IsRange(start, end uint32, value bool) (bool, error) {
	if end < start || end > ba.size {
		return false, errors.New("Invalid argument passed to IsRange on BitArray")
	}

	if end == start {
		return true, nil
	}
	end--

	firstInt := start / 32
	lastInt := end / 32
	for i := firstInt; i <= lastInt; i++ {
		var firstBit, lastBit uint32

		if i > firstInt {
			firstBit = 0
		} else {
			firstBit = (start & 0x1F)
		}

		if i < lastInt {
			lastBit = 31
		} else {
			lastBit = (end & 0x1F)
		}

		mask := uint32((2 << lastBit) - (1 << firstBit))

		if (value && ((ba.bits[i] & mask) != mask)) || (!value && ((ba.bits[i] & mask) != 0)) {
			return false, nil
		}
	}

	return true, nil
}

func (ba *BitArray) AppendBit(bit bool) {
	ba.EnsureCapacity(ba.size + 1)
	if bit {
		ba.bits[(ba.size / 32)] |= (1 << (ba.size & 0x1F))
	}
	ba.size++
}

func (ba *BitArray) AppendBits(value, numBits uint32) error {
	if numBits > 32 {
		return errors.New("Num bits must be between 0 and 32")
	}
	ba.EnsureCapacity(ba.size + numBits)
	for numBitsLeft := numBits; numBitsLeft > 0; numBitsLeft-- {
		ba.AppendBit(((value >> (numBitsLeft - 1)) & 0x01) == 1)
	}

	return nil
}

func (ba *BitArray) AppendBitArray(other BitArray) {
	otherSize := other.size
	ba.EnsureCapacity(ba.size + otherSize)
	for i := uint32(0); i < otherSize; i++ {
		ba.AppendBit(other.Get(i))
	}
}

func (ba *BitArray) Xor(other BitArray) error {
	if ba.size != other.size {
		return errors.New("Sizes don't match")
	}

	for i := uint32(0); i < uint32(len(ba.bits)); i++ {
		ba.bits[i] ^= other.bits[i]
	}

	return nil
}

func (ba *BitArray) ToBytes(bitOffset uint32, array []uint8, offset, numBytes uint32) {
	for i := uint32(0); i < numBytes; i++ {
		theByte := 0
		for j := uint32(0); j < 8; j++ {
			if ba.Get(bitOffset) {
				theByte |= 1 << (7 - j)
			}
			bitOffset++
		}
		array[offset+i] = uint8(theByte)
	}
}

func (ba *BitArray) GetBitArray() []uint32 {
	return ba.bits
}

func (ba *BitArray) Reverse() {
	newBits := make([]uint32, len(ba.bits))

	len := ((ba.size - 1) / 32)
	oldBitsLen := len + 1

	for i := uint32(0); i < oldBitsLen; i++ {
		x := ba.bits[i]

		x = (((x >> 1) & 0x55555555) | ((x & 0x55555555) << 1))
		x = (((x >> 2) & 0x33333333) | ((x & 0x33333333) << 2))
		x = (((x >> 4) & 0x0F0F0F0F) | ((x & 0x0F0F0F0F) << 4))
		x = (((x >> 8) & 0x00FF00FF) | ((x & 0x00FF00FF) << 8))
		x = (((x >> 16) & 0x0000FFFF) | ((x & 0x0000FFFF) << 16))
		newBits[len-i] = x
	}

	if ba.size != (oldBitsLen * 32) {
		leftOffset := ((oldBitsLen * 32) - ba.size)
		currentInt := newBits[0] >> leftOffset
		for i := uint32(1); i < oldBitsLen; i++ {
			nextInt := newBits[i]
			currentInt |= (nextInt << (32 - leftOffset))
			newBits[i-1] = currentInt
		}
		newBits[oldBitsLen-1] = currentInt
	}
	ba.bits = newBits
}

func (ba *BitArray) Equals(other BitArray) bool {
	if ba.size != other.size {
		return false
	}

	for i := 0; i < len(ba.bits); i++ {
		if ba.bits[i] != other.bits[i] {
			return false
		}
	}

	return true
}

func (ba *BitArray) String() string {
	var result bytes.Buffer

	for i := uint32(0); i < ba.size; i++ {
		if (i & 0x07) == 0 {
			result.WriteString(" ")
		}

		if ba.Get(i) {
			result.WriteString("X")
		} else {
			result.WriteString(".")
		}
	}

	return result.String()
}

func (ba *BitArray) Clone() *BitArray {
	return &BitArray{internal.SliceCloneU32(ba.bits), ba.size}
}
