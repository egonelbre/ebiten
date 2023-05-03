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

type Vector2 struct {
	X, Y float32
}

func (a Vector2) Add(b Vector2) Vector2 {
	return Vector2{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a Vector2) Sub(b Vector2) Vector2 {
	return Vector2{X: a.X - b.X, Y: a.Y - b.Y}
}

func (a Vector2) Scale(s float32) Vector2 {
	return Vector2{X: a.X * s, Y: a.Y * s}
}
