// Copyright 2023 The Ebiten Authors
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

package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600

	elasticity       = 0.85
	startBunnies     = 10000
	bunniesPerSecond = 10000
)

type Mode struct {
	Colored bool

	BunniesPerSecond int
}

type Game struct {
	Mode Mode

	Size    Vector2
	Gravity Vector2

	Bunnies []Bunny

	Generating          bool
	CurrentBunnyVariant int
	BunnyVariants       []*ebiten.Image
}

func (g *Game) Update() error {
	dt := 1 / float32(ebiten.TPS())
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Generating = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Generating = false
		g.CurrentBunnyVariant = (g.CurrentBunnyVariant + 1) % len(g.BunnyVariants)
	}
	if g.Generating {
		g.GenerateBunnies(int(float32(g.Mode.BunniesPerSecond) * dt))
	}

	for i := range g.Bunnies {
		b := &g.Bunnies[i]
		b.Update(g, dt)
	}
	return nil
}

func (g *Game) GenerateBunnies(n int) {
	start := len(g.Bunnies)
	g.Bunnies = append(g.Bunnies, make([]Bunny, n)...)
	for i := start; i < len(g.Bunnies); i++ {
		g.Bunnies[i].init(g)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	switch {
	case !g.Mode.Colored:
		for i := range g.Bunnies {
			b := &g.Bunnies[i]
			b.Draw(screen)
		}
	case g.Mode.Colored:
		for i := range g.Bunnies {
			b := &g.Bunnies[i]
			b.DrawColored(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.Size.X), int(g.Size.Y)
}

type Bunny struct {
	Sprite     *ebiten.Image
	Color      ebiten.ColorScale
	Position   Vector2
	Size       Vector2
	Velocity   Vector2
	Elasticity float32
}

func (bunny *Bunny) init(game *Game) {
	bunny.Sprite = game.BunnyVariants[game.CurrentBunnyVariant]
	bunny.Velocity.X = rand.Float32() * screenWidth * 0.2
	bunny.Velocity.Y = rand.Float32() * screenHeight * 0.05
	bunny.Size.X = float32(bunny.Sprite.Bounds().Dx())
	bunny.Size.Y = float32(bunny.Sprite.Bounds().Dy())
	bunny.Elasticity = elasticity + (1-elasticity)*rand.Float32()
}

func (bunny *Bunny) Update(game *Game, dt float32) {
	bunny.Position = bunny.Position.Add(bunny.Velocity.Scale(dt))
	bunny.Velocity = bunny.Velocity.Add(game.Gravity.Scale(dt))

	if bunny.Position.X+bunny.Size.X > game.Size.X && bunny.Velocity.X > 0 {
		bunny.Velocity.X *= -1
	}
	if bunny.Position.X < 0 && bunny.Velocity.X < 0 {
		bunny.Velocity.X *= -1
	}
	if bunny.Position.Y+bunny.Size.Y > game.Size.Y && bunny.Velocity.Y > 0 {
		bunny.Velocity.Y *= -bunny.Elasticity

		//TODO: add energy back into the system
	}
}

func (bunny *Bunny) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(bunny.Position.X), float64(bunny.Position.Y))
	screen.DrawImage(bunny.Sprite, op)
}

func (bunny *Bunny) DrawColored(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(bunny.Position.X), float64(bunny.Position.Y))
	op.ColorScale = bunny.Color
	screen.DrawImage(bunny.Sprite, op)
}

func main() {
	const (
		frameOX     = 0
		frameOY     = 32
		frameWidth  = 32
		frameHeight = 32
		frameCount  = 8
	)

	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	tileset := ebiten.NewImageFromImage(img)

	game := &Game{
		Mode: Mode{
			BunniesPerSecond: bunniesPerSecond,
		},
		Size: Vector2{
			X: screenWidth,
			Y: screenHeight,
		},
		Gravity: Vector2{
			X: 0,
			Y: 300,
		},
	}

	for i := 0; i < frameCount; i++ {
		sx, sy := frameOX+i*frameWidth, frameOY
		frame := tileset.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
		game.BunnyVariants = append(game.BunnyVariants, frame)
	}

	game.GenerateBunnies(startBunnies)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Bunnymark (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
