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

//go:build example
// +build example

package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var mplusNormalFont font.Face

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	inputField *inpututil.InputField

	w, h int
}

func NewGame() *Game {
	g := &Game{
		inputField: inpututil.NewInputField(mplusNormalFont),
	}
	return g
}

func (g *Game) Update() error {
	err := g.inputField.Update()
	if err != nil {
		return fmt.Errorf("failed to update input field: %s", err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.inputField.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == g.w && outsideHeight == g.h {
		return outsideWidth, outsideHeight
	}

	padding := 10

	w, h := outsideWidth-padding, 100
	if h > outsideHeight {
		h = outsideHeight
	}

	x, y := outsideWidth/2-w/2, outsideHeight/2-h/2

	g.inputField.SetRect(image.Rect(x, y, x+w, y+h))

	g.w, g.h = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("InputField (Ebiten Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
