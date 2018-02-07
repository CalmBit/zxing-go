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

package core

type InvertedLuminanceSource struct {
	height, width int
	delegate      LuminanceSource
}

func NewInvertedLuminanceSource(delegate LuminanceSource) *InvertedLuminanceSource {
	return &InvertedLuminanceSource{delegate.GetHeight(), delegate.GetWidth(), delegate}
}
func (this *InvertedLuminanceSource) GetRow(y int, row []uint8) []uint8 {
	row = this.delegate.GetRow(y, row)
	width := this.GetWidth()
	for i := 0; i < width; i++ {
		row[i] = (255 - row[i])
	}
	return row
}

func (this *InvertedLuminanceSource) GetMatrix() [][]uint8 {
	matrix := this.delegate.GetMatrix()
	Inverted := make([][]uint8, this.GetHeight())
	for i := range Inverted {
		Inverted[i] = make([]uint8, this.GetWidth())
		for j := range Inverted[i] {
			Inverted[i][j] = (255 - matrix[i][j])
		}
	}
	return Inverted
}

func (this *InvertedLuminanceSource) GetWidth() int {
	return this.width
}

func (this *InvertedLuminanceSource) GetHeight() int {
	return this.height
}

func (this *InvertedLuminanceSource) IsCropSupported() bool {
	return this.delegate.IsCropSupported()
}

func (this *InvertedLuminanceSource) Crop(left, top, width, height int) LuminanceSource {
	return NewInvertedLuminanceSource(this.delegate.Crop(left, top, width, height))
}

func (this *InvertedLuminanceSource) IsRotateSupported() bool {
	return this.delegate.IsRotateSupported()
}

func (this *InvertedLuminanceSource) Invert() LuminanceSource {
	return this.delegate
}

func (this *InvertedLuminanceSource) RotateCounterClockwise() LuminanceSource {
	return NewInvertedLuminanceSource(this.delegate.RotateCounterClockwise())
}

func (this *InvertedLuminanceSource) RotateCounterClockwise45() LuminanceSource {
	return NewInvertedLuminanceSource(this.delegate.RotateCounterClockwise45())
}
