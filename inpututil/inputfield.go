package inpututil

// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"image"
	"image/color"
	"strings"
	"sync"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const initialPadding = 2

// InputField is a text input field. Call Update and Draw when your Game's
// Update and Draw methods are called. The field is hidden by default. Set the
// position and size of the field by calling SetRect to enable display.
type InputField struct {
	// r is the position of the field.
	r image.Rectangle

	// img is the cached image of the field.
	img *ebiten.Image

	// buffer is the actual content of the field.
	buffer string

	// bufferWrapped is the content of the field as it appears on the screen.
	bufferWrapped []string

	// face is the font face of the text within the field.
	face font.Face

	// textHeight is the height of a single line of text.
	textHeight int

	// textColor is the color of the text within the field.
	textColor color.Color

	// padding is the amount of padding around the text within the field.
	padding int

	sync.Mutex
}

func NewInputField(face font.Face) *InputField {
	f := &InputField{
		face:      face,
		textColor: color.RGBA{0, 0, 0, 255},
		padding:   initialPadding,
	}

	bounds := text.BoundString(f.face, "ATZgpq.")
	f.textHeight = bounds.Dy()

	return f
}

func (f *InputField) GetRect() image.Rectangle {
	f.Lock()
	defer f.Unlock()

	return f.r
}

func (f *InputField) SetRect(r image.Rectangle) error {
	f.Lock()
	defer f.Unlock()

	f.r = r
	if rectIsZero(r) {
		f.img = nil
		return nil
	}

	f.img = ebiten.NewImage(f.r.Dx(), f.r.Dy())
	return f.drawFieldImage()
}

func (f *InputField) wrapContent() {
	f.bufferWrapped = nil
	w := f.r.Dx()
	for _, line := range strings.Split(f.buffer, "\n") {
		l := len(line)
		var start int
		var end int
		for start < l {
			for end = l; end > start; end-- {
				bounds := text.BoundString(f.face, line[start:end])
				if bounds.Dx() < w-(f.padding*2) {
					// Break on whitespace.
					if end < l && !unicode.IsSpace(rune(line[end])) {
						for endOffset := 0; endOffset < end-start; endOffset++ {
							if unicode.IsSpace(rune(line[end-endOffset])) {
								end = end - endOffset
								break
							}
						}
					}
					f.bufferWrapped = append(f.bufferWrapped, line[start:end])
					break
				}
			}
			start = end
		}
	}
}

func (f *InputField) Update() error {
	f.Lock()
	defer f.Unlock()

	var redraw bool
	if IsKeyJustPressed(ebiten.KeyBackspace) && len(f.buffer) > 0 {
		f.buffer = f.buffer[:len(f.buffer)-1]
		redraw = true
	}

	if IsKeyJustPressed(ebiten.KeyEnter) {
		f.buffer += "\n"
	}

	var b []rune
	b = ebiten.AppendInputChars(b[:0])
	if len(b) > 0 {
		f.buffer += string(b)
		redraw = true
	}

	if redraw {
		return f.drawFieldImage()
	}
	return nil
}

func (f *InputField) drawFieldImage() error {
	if f.img == nil {
		return nil
	}

	f.wrapContent()

	f.img.Fill(color.RGBA{255, 255, 255, 255})

	for i, line := range f.bufferWrapped {
		text.Draw(f.img, line, f.face, 0, f.textHeight*(i+1), f.textColor)
	}
	return nil
}

func (f *InputField) Draw(screen *ebiten.Image) {
	f.Lock()
	defer f.Unlock()

	if f.img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.r.Min.X), float64(f.r.Min.Y))
	screen.DrawImage(f.img, op)
}

func rectIsZero(r image.Rectangle) bool {
	return r == image.Rectangle{}
}
