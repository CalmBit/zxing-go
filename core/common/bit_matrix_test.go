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
	"strconv"
	"testing"

	"github.com/discesoft/zxing-go/core/common"
	"github.com/discesoft/zxing-go/core/internal"
)

func TestGetSet(t *testing.T) {
	matrix, err := common.NewBitMatrixFromDimension(33)
	internal.AssertSuccess(t, err)
	internal.AssertEquals(t, uint32(33), matrix.GetHeight(), "matrix initialized with height "+strconv.Itoa(int(matrix.GetHeight()))+" instead of 33")
	for y := uint32(0); y < 33; y++ {
		for x := uint32(0); x < 33; x++ {
			if ((y * x) % 3) == 0 {
				matrix.Set(x, y)
			}
		}
	}
	for y := uint32(0); y < 33; y++ {
		for x := uint32(0); x < 33; x++ {
			internal.AssertEquals(
				t,
				(((y * x) % 3) == 0),
				matrix.Get(x, y),
				"value at "+strconv.Itoa(int(x))+","+strconv.Itoa(int(y))+" was incorrect")
		}
	}
}

func TestSetRegion(t *testing.T) {
	matrix, err := common.NewBitMatrixFromDimension(5)
	internal.AssertSuccess(t, err)
	err = matrix.SetRegion(1, 1, 3, 3)
	internal.AssertSuccess(t, err)
	for y := uint32(0); y < 5; y++ {
		for x := uint32(0); x < 5; x++ {
			internal.AssertEquals(
				t,
				(y >= 1 && y <= 3 && x >= 1 && x <= 3),
				matrix.Get(x, y),
				"value at "+strconv.Itoa(int(x))+","+strconv.Itoa(int(y))+" was incorrect")
		}
	}
}

func TestEnclosing(t *testing.T) {
	matrix, err := common.NewBitMatrixFromDimension(5)
	internal.AssertSuccess(t, err)
	b := matrix.GetEnclosingRectangle()
	internal.AssertNil(t, b, "enclosing rectangle was not nil on init")
	internal.AssertSuccess(t, matrix.SetRegion(1, 1, 1, 1))
	internal.AssertSlicesEqualU32(t, []uint32{1, 1, 1, 1}, matrix.GetEnclosingRectangle(), "enclosing rectangle was not equal to {1,1,1,1}")
	internal.AssertSuccess(t, matrix.SetRegion(1, 1, 3, 2))
	internal.AssertSlicesEqualU32(t, []uint32{1, 1, 3, 2}, matrix.GetEnclosingRectangle(), "enclosing rectangle was not equal to {1,1,3,2}")
	internal.AssertSuccess(t, matrix.SetRegion(0, 0, 5, 5))
	internal.AssertSlicesEqualU32(t, []uint32{0, 0, 5, 5}, matrix.GetEnclosingRectangle(), "enclosing rectangle was not equal to {0,0,5,5}")
}

func TestOnBit(t *testing.T) {
	matrix, err := common.NewBitMatrixFromDimension(5)
	internal.AssertSuccess(t, err)
	internal.AssertNil(t, matrix.GetTopLeftOnBit(), "top left on bit not nil on init")
	internal.AssertNil(t, matrix.GetBottomRightOnBit(), "bottom right on bit not nil on init")
	matrix.SetRegion(1, 1, 1, 1)
	internal.AssertSlicesEqualU32(t, []uint32{1, 1}, matrix.GetTopLeftOnBit(), "1 | top left on bit not {1,1}")
	internal.AssertSlicesEqualU32(t, []uint32{1, 1}, matrix.GetBottomRightOnBit(), "1 | bottom right on bit not {1,1}")
	matrix.SetRegion(1, 1, 3, 2)
	internal.AssertSlicesEqualU32(t, []uint32{1, 1}, matrix.GetTopLeftOnBit(), "2 | top left on bit not {1,1}")
	internal.AssertSlicesEqualU32(t, []uint32{3, 2}, matrix.GetBottomRightOnBit(), "2 | bottom right on bit not {3,2}")
	matrix.SetRegion(0, 0, 5, 5)
	internal.AssertSlicesEqualU32(t, []uint32{0, 0}, matrix.GetTopLeftOnBit(), "3 | top left on bit not {0,0}")
	internal.AssertSlicesEqualU32(t, []uint32{4, 4}, matrix.GetBottomRightOnBit(), "3 | bottom right on bit not {4,4}")

}

func TestRectangularMatrix(t *testing.T) {
	matrix, err := common.NewBitMatrix(75, 20)
	internal.AssertSuccess(t, err)
	internal.AssertEquals(t, uint32(75), matrix.GetWidth(), "matrix width not 75")
	internal.AssertEquals(t, uint32(20), matrix.GetHeight(), "matrix height not 20")
	matrix.Set(10, 0)
	matrix.Set(11, 1)
	matrix.Set(50, 2)
	matrix.Set(51, 3)
	matrix.Flip(74, 4)
	matrix.Flip(0, 5)

	internal.AssertEquals(t, true, matrix.Get(10, 0), "value at {10,0} not true")
	internal.AssertEquals(t, true, matrix.Get(11, 1), "value at {11,1} not true")
	internal.AssertEquals(t, true, matrix.Get(50, 2), "value at {50,2} not true")
	internal.AssertEquals(t, true, matrix.Get(51, 3), "value at {51,3} not true")
	internal.AssertEquals(t, true, matrix.Get(74, 4), "value at {74,4} not true")
	internal.AssertEquals(t, true, matrix.Get(0, 5), "value at {0,5} not true")

	matrix.Flip(50, 2)
	matrix.Flip(51, 3)
	internal.AssertEquals(t, false, matrix.Get(50, 2), "value at {50,2} didn't flip off")
	internal.AssertEquals(t, false, matrix.Get(51, 3), "value at {51,3} didn't flip off")
}

func TestRectangularSetRegion(t *testing.T) {
	matrix, err := common.NewBitMatrix(320, 240)
	internal.AssertSuccess(t, err)
	internal.AssertEquals(t, uint32(320), matrix.GetWidth(), "matrix width not 320")
	internal.AssertEquals(t, uint32(240), matrix.GetHeight(), "matrix height not 240")
	matrix.SetRegion(105, 22, 80, 12)

	for y := uint32(0); y < 240; y++ {
		for x := uint32(0); x < 320; x++ {
			internal.AssertEquals(t, y >= 22 && y < 34 && x >= 105 && x < 185, matrix.Get(x, y), "value at {"+strconv.Itoa(int(x))+","+strconv.Itoa(int(y))+"} was incorrect")
		}
	}
}

func TestGetRow(t *testing.T) {
	matrix, err := common.NewBitMatrix(102, 5)
	internal.AssertSuccess(t, err)
	for x := uint32(0); x < 102; x++ {
		if (x & 0x03) == 0 {
			matrix.Set(x, 2)
		}
	}

	array := matrix.GetRow(2, nil)
	internal.AssertEquals(t, uint32(102), array.GetSize(), "newly allocated array was not size 102")

	array2 := common.NewBitArray(60)
	array2 = matrix.GetRow(2, array2)
	internal.AssertEquals(t, uint32(102), array2.GetSize(), "reallocated array was not size 102")

	array3 := common.NewBitArray(200)
	array3 = matrix.GetRow(2, array3)
	internal.AssertEquals(t, uint32(200), array3.GetSize(), "array must have been reallocated, size not 200")

	for x := uint32(0); x < 102; x++ {
		on := ((x & 0x03) == 0)
		internal.AssertEquals(t, on, array.Get(x), "array had wrong value at {"+strconv.Itoa(int(x))+",2}")
		internal.AssertEquals(t, on, array2.Get(x), "array2 had wrong value at {"+strconv.Itoa(int(x))+",2}")
		internal.AssertEquals(t, on, array3.Get(x), "array3 had wrong value at {"+strconv.Itoa(int(x))+",2}")
	}
}

func TestRotate180Simple(t *testing.T) {
	matrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	matrix.Set(0, 0)
	matrix.Set(0, 1)
	matrix.Set(1, 2)
	matrix.Set(2, 1)

	internal.AssertSuccess(t, matrix.Rotate180())

	internal.AssertEquals(t, true, matrix.Get(2, 2), "flip: {2,2} was false")
	internal.AssertEquals(t, true, matrix.Get(2, 1), "flip: {2,1} was false")
	internal.AssertEquals(t, true, matrix.Get(1, 0), "flip: {1,0} was false")
	internal.AssertEquals(t, true, matrix.Get(0, 1), "flip: {0,1} was false")
}

func TestRotate180(t *testing.T) {
	testRotate180(t, 7, 4)
	testRotate180(t, 7, 5)
	testRotate180(t, 8, 4)
	testRotate180(t, 8, 5)
}

func TestParse(t *testing.T) {
	var fullMatrix, centerMatrix, emptyMatrix24 *common.BitMatrix

	emptyMatrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)

	fullMatrix, err = common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	fullMatrix.SetRegion(0, 0, 3, 3)

	centerMatrix, err = common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	centerMatrix.SetRegion(1, 1, 1, 1)

	emptyMatrix24, err = common.NewBitMatrix(2, 4)
	internal.AssertSuccess(t, err)

	result, err := common.ParseStringToBitMatrix("   \n   \n   \n", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, emptyMatrix, result)

	result, err = common.ParseStringToBitMatrix("   \n   \r\r\n   \n\r", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, emptyMatrix, result)

	result, err = common.ParseStringToBitMatrix("   \n   \n   ", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, emptyMatrix, result)

	result, err = common.ParseStringToBitMatrix("xxx\nxxx\nxxx\n", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, fullMatrix, result)

	result, err = common.ParseStringToBitMatrix("   \n x \n   \n", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, centerMatrix, result)

	result, err = common.ParseStringToBitMatrix("      \n  x   \n      \n", "x ", "  ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, centerMatrix, result)

	result, err = common.ParseStringToBitMatrix("   \n xy\n   \n", "x", " ")
	internal.AssertFailure(t, err, "parse should have failed")

	result, err = common.ParseStringToBitMatrix("  \n  \n  \n  \n", "x", " ")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, emptyMatrix24, result)

	result, err = common.ParseStringToBitMatrix(centerMatrix.ToString("x", "."), "x", ".")
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, centerMatrix, result)
}

func TestUnset(t *testing.T) {
	emptyMatrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	matrix := emptyMatrix.Clone()
	matrix.Set(1, 1)
	assertMatrixInequality(t, emptyMatrix, matrix)
	matrix.Unset(1, 1)
	assertMatrixEquality(t, emptyMatrix, matrix)
	matrix.Unset(1, 1)
	assertMatrixEquality(t, emptyMatrix, matrix)
}

func TestXOR(t *testing.T) {
	emptyMatrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)

	fullMatrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	fullMatrix.SetRegion(0, 0, 3, 3)

	centerMatrix, err := common.NewBitMatrix(3, 3)
	internal.AssertSuccess(t, err)
	centerMatrix.SetRegion(1, 1, 1, 1)

	invertedCenterMatrix := fullMatrix.Clone()
	invertedCenterMatrix.Unset(1, 1)

	badMatrix, err := common.NewBitMatrix(4, 4)
	internal.AssertSuccess(t, err)

	testXor(t, emptyMatrix, emptyMatrix, emptyMatrix)
	testXor(t, emptyMatrix, centerMatrix, centerMatrix)
	testXor(t, emptyMatrix, fullMatrix, fullMatrix)

	testXor(t, centerMatrix, emptyMatrix, centerMatrix)
	testXor(t, centerMatrix, centerMatrix, emptyMatrix)
	testXor(t, centerMatrix, fullMatrix, invertedCenterMatrix)

	testXor(t, invertedCenterMatrix, emptyMatrix, invertedCenterMatrix)
	testXor(t, invertedCenterMatrix, centerMatrix, fullMatrix)
	testXor(t, invertedCenterMatrix, fullMatrix, centerMatrix)

	testXor(t, fullMatrix, emptyMatrix, fullMatrix)
	testXor(t, fullMatrix, centerMatrix, invertedCenterMatrix)
	testXor(t, fullMatrix, fullMatrix, emptyMatrix)

	err = emptyMatrix.Clone().Xor(badMatrix)
	internal.AssertFailure(t, err, "empty matrix can't XOR against badMatrix :(")

	err = badMatrix.Clone().Xor(emptyMatrix)
	internal.AssertFailure(t, err, "bad matrix can't XOR against emptyMatrix :(")
}

func testXor(t *testing.T, dataMatrix, flipMatrix, expectedMatrix *common.BitMatrix) {
	matrix := dataMatrix.Clone()
	err := matrix.Xor(flipMatrix)
	internal.AssertSuccess(t, err)
	assertMatrixEquality(t, expectedMatrix, matrix)
}

func assertMatrixInequality(t *testing.T, a, b *common.BitMatrix) {
	internal.AssertEquals(t, false, a.Equals(b), "matrix equality detected")
}

func assertMatrixEquality(t *testing.T, a, b *common.BitMatrix) {
	internal.AssertEquals(t, true, a.Equals(b), "matrix inequality detected")
}

func testRotate180(t *testing.T, width, height uint32) {
	input := getInput(width, height)
	input.Rotate180()
	expected := getExpected(width, height)
	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			internal.AssertEquals(t, expected.Get(x, y), input.Get(x, y), "flip180: non-matching value at {"+strconv.Itoa(int(x))+","+strconv.Itoa(int(y))+"}")
		}
	}
}

func getExpected(width, height uint32) *common.BitMatrix {
	result, _ := common.NewBitMatrix(width, height)
	points := []uint32{1, 2, 2, 0, 3, 1}

	for i := 0; i < len(points); i += 2 {
		result.Set(width-1-points[i], height-1-points[i+1])
	}
	return result
}

func getInput(width, height uint32) *common.BitMatrix {
	// ignore errors for now
	result, _ := common.NewBitMatrix(width, height)
	points := []uint32{1, 2, 2, 0, 3, 1}

	for i := 0; i < len(points); i += 2 {
		result.Set(points[i], points[i+1])
	}
	return result
}
