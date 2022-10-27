package unpack

import (
    "os"
    "time"
    "image/color"
    "encoding/json"

    "example/CRT/neopixel"
)

const (
    WIDTH, HEIGHT = 15, 11
)

type Anim struct {
	Frames []Frame
    Fps int
    Loop bool
    LoopDelay time.Duration
}

type Frame [WIDTH][HEIGHT]color.RGBA

func (f *Frame) Draw(mat *neopixel.Matrix) {
    for x, col := range f {
        for y, rgb := range col {
            mat.Set(x, mat.Height-(y+1), rgb)
        }
    }
}

func UnpackFile(path string) (*Anim, error) {
    var unpacked Anim
    file, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    json.Unmarshal([]byte(file), &unpacked)

	var newframes []Frame
	for _, f := range unpacked.Frames {
		newframes = append([]Frame{f}, newframes...)
	}

    return &unpacked, nil
}


