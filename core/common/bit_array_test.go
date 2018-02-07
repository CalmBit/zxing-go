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

package common_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/discesoft/zxing-go/core/internal"

	"github.com/discesoft/zxing-go/core/common"
)

func TestBitArray_GetSet(t *testing.T) {
	array := common.NewBitArray(33)
	for i := uint32(0); i < 33; i++ {
		internal.AssertFalse(t, array.Get(i), "uninitialized point {"+strconv.Itoa(int(i))+"} should be false")
		array.Set(i)
		internal.AssertTrue(t, array.Get(i), "initialized point {"+strconv.Itoa(int(i))+"} should be true")
	}
}

func TestBitArray_GetNextSet1(t *testing.T) {
	array := common.NewBitArray(32)
	for i := uint32(0); i < array.GetSize(); i++ {
		internal.AssertEquals(t, uint32(32), array.GetNextSet(i), strconv.Itoa(int(i)))
	}

	array = common.NewBitArray(33)
	for i := uint32(0); i < array.GetSize(); i++ {
		internal.AssertEquals(t, uint32(33), array.GetNextSet(i), strconv.Itoa(int(i)))
	}
}

func TestBitArray_GetNextSet2(t *testing.T) {
	array := common.NewBitArray(33)
	array.Set(31)
	for i := uint32(0); i < array.GetSize(); i++ {
		if i <= 31 {
			internal.AssertEquals(t, uint32(31), array.GetNextSet(i), strconv.Itoa(int(i)))
		} else {
			internal.AssertEquals(t, uint32(33), array.GetNextSet(i), strconv.Itoa(int(i)))
		}
	}
	array = common.NewBitArray(33)
	array.Set(32)
	for i := uint32(0); i < array.GetSize(); i++ {
		internal.AssertEquals(t, uint32(32), array.GetNextSet(i), strconv.Itoa(int(i)))
	}
}

func TestBitArray_GetNextSet3(t *testing.T) {
	array := common.NewBitArray(63)
	array.Set(31)
	array.Set(32)
	for i := uint32(0); i < array.GetSize(); i++ {
		var expected uint32
		if i <= 31 {
			expected = 31
		} else if i == 32 {
			expected = 32
		} else {
			expected = 63
		}
		internal.AssertEquals(t, expected, array.GetNextSet(i), strconv.Itoa(int(i)))
	}
}

func TestBitArray_GetNextSet4(t *testing.T) {
	array := common.NewBitArray(63)
	array.Set(33)
	array.Set(40)
	for i := uint32(0); i < array.GetSize(); i++ {
		var expected uint32
		if i <= 33 {
			expected = 33
		} else if i <= 40 {
			expected = 40
		} else {
			expected = 63
		}
		internal.AssertEquals(t, expected, array.GetNextSet(i), strconv.Itoa(int(i)))
	}
}

func TestBitArray_GetNextSet5(t *testing.T) {
	r := rand.New(rand.NewSource(0xDEADBEEF))
	for i := 0; i < 10; i++ {
		array := common.NewBitArray(1 + uint32(r.Intn(100)))
		numSet := uint32(r.Intn(20))
		for j := uint32(0); j < numSet; j++ {
			array.Set(uint32(r.Intn(int(array.GetSize()))))
		}
		numQueries := uint32(r.Intn(20))
		for j := uint32(0); j < numQueries; j++ {
			query := uint32(r.Intn(int(array.GetSize())))
			expected := query
			for (expected < array.GetSize()) && !(array.Get(expected)) {
				expected++
			}
			actual := array.GetNextSet(query)
			internal.AssertEquals(t, expected, actual, "GetNextSet5 failed")
		}
	}
}

func TestBitArray_SetBulk(t *testing.T) {
	array := common.NewBitArray(64)
	array.SetBulk(32, 0xFFFF0000)
	for i := uint32(0); i < 48; i++ {
		internal.AssertFalse(t, array.Get(i), "false expected for "+strconv.Itoa(int(i)))
	}
	for i := uint32(48); i < 64; i++ {
		internal.AssertTrue(t, array.Get(i), "true expected for "+strconv.Itoa(int(i)))
	}
}

func TestBitArray_SetRange(t *testing.T) {
	array := common.NewBitArray(64)
	array.SetRange(28, 36)
	internal.AssertFalse(t, array.Get(27), "false expected for 27")
	for i := uint32(28); i < 36; i++ {
		internal.AssertTrue(t, array.Get(i), "true expected for "+strconv.Itoa(int(i)))
	}
	internal.AssertFalse(t, array.Get(36), "false expected for 36")
}

func TestBitArray_Clear(t *testing.T) {
	array := common.NewBitArray(32)
	for i := uint32(0); i < 32; i++ {
		array.Set(i)
	}
	array.Clear()
	for i := uint32(0); i < 32; i++ {
		internal.AssertFalse(t, array.Get(i), "false expected at "+strconv.Itoa(int(i)))
	}
}

func TestBitArray_Flip(t *testing.T) {
	array := common.NewBitArray(32)
	internal.AssertFalse(t, array.Get(5), "false expected for 5")
}

func TestBitArray_GetArray(t *testing.T) {
	array := common.NewBitArray(64)
	array.Set(0)
	array.Set(63)
	ints := array.GetBitArray()
	internal.AssertEquals(t, uint32(1), ints[0], "expected {0} = 1")
	internal.AssertEquals(t, uint32(0x80000000), ints[1], "expected {1} = 0x80000000")
}

func TestBitArray_IsRange(t *testing.T) {
	array := common.NewBitArray(64)

	result, err := array.IsRange(0, 64, false)
	internal.AssertSuccess(t, err)
	internal.AssertTrue(t, result, "expected range [0,64] to be false")

	result, err = array.IsRange(0, 64, true)
	internal.AssertSuccess(t, err)
	internal.AssertFalse(t, result, "expected range [0,64] to not be true")

	array.Set(32)
	result, err = array.IsRange(32, 33, true)
	internal.AssertSuccess(t, err)
	internal.AssertTrue(t, result, "expected range [32,33] to be true")

	array.Set(31)
	result, err = array.IsRange(31, 33, true)
	internal.AssertSuccess(t, err)
	internal.AssertTrue(t, result, "expected range [31,33] to be true")

	array.Set(35)
	result, err = array.IsRange(31, 35, true)
	internal.AssertSuccess(t, err)
	internal.AssertFalse(t, result, "expected range [31,33] to not be true")

	for i := uint32(0); i < 31; i++ {
		array.Set(i)
	}
	result, err = array.IsRange(0, 33, true)
	internal.AssertSuccess(t, err)
	internal.AssertTrue(t, result, "expected range [0,33] to be true")

	for i := uint32(33); i < 64; i++ {
		array.Set(i)
	}

	result, err = array.IsRange(0, 64, true)
	internal.AssertSuccess(t, err)
	internal.AssertTrue(t, result, "expected range [0,64] to be true")

	result, err = array.IsRange(0, 64, false)
	internal.AssertSuccess(t, err)
	internal.AssertFalse(t, result, "expected range [0,64] to not be false")

}

func TestBitArray_ReverseAlgorithmTest(t *testing.T) {
	oldBits := []uint32{128, 256, 512, 6453324, 50934953}
	for size := uint32(1); size < 160; size++ {
		newBitsOriginal := reverseOriginal(internal.SliceCloneU32(oldBits), size)
		newBitArray := common.NewBitArrayTesting(internal.SliceCloneU32(oldBits), size)
		newBitArray.Reverse()
		newBitsNew := newBitArray.GetBitArray()
		internal.AssertSlicesEqualU32(t, newBitsOriginal, newBitsNew, strconv.Itoa(int(size)))
	}
}

func TestBitArray_TestClone(t *testing.T) {
	array := common.NewBitArray(32)
	array.Clone().Set(0)
	internal.AssertFalse(t, array.Get(0), "expected {0} after clone to be false")
}

func TestBitArray_TestEquals(t *testing.T) {
	a := common.NewBitArray(32)
	b := common.NewBitArray(32)
	internal.AssertTrue(t, a.Equals(b), "a != b!!!")
	internal.AssertFalse(t, a.Equals(common.NewBitArray(31)), "new array == a!!!")
	a.Set(16)
	internal.AssertFalse(t, a.Equals(common.NewBitArray(31)), "new array == a!!!")
	internal.AssertFalse(t, a.Equals(b), "a = b!!!")
	b.Set(16)
	internal.AssertTrue(t, a.Equals(b), "a != b!!!")

}

func reverseOriginal(oldBits []uint32, size uint32) []uint32 {
	newBits := make([]uint32, len(oldBits))
	for i := uint32(0); i < size; i++ {
		if bitSet(oldBits, size-i-1) {
			newBits[i/32] |= (1 << (i & 0x1F))
		}
	}
	return newBits
}

func bitSet(bits []uint32, i uint32) bool {
	return (bits[i/32] & (1 << (i & 0x1F))) != 0
}
