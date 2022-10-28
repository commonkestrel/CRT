// A utility package to make rpi-ws821x-go simpler and more intuitive
package neopixel

import (
    "errors"
    "image/color"

    "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Strip struct {
    Strip *ws2811.WS2811
    Length int
}

func NewStrip(pin, length int) (*Strip, error) {
    pins := [4]int{10, 12, 18, 21}
    var valid bool
    for _, v := range pins {
        if pin == v {
            valid = true
        }
    }
    if !valid {
        return nil, errors.New("invalid pin. valid pins are 10, 12, 18, and 21")
    }

    opt := ws2811.DefaultOptions
    opt.Channels[0].LedCount = length
    opt.Channels[0].GpioPin = pin
    strip, err := ws2811.MakeWS2811(&opt)
    if err != nil {
        return nil, err        
    }
    strip.Init()
    return &Strip{strip, length}, nil
}

func colorToBytes(rgb color.Color) uint32 {
    r, g, b, _ := rgb.RGBA()
    return ((r>>8)&0xff)<<16 + ((g>>8)&0xff)<<8 + ((b >> 8) & 0xff)
}

func (s *Strip) Set(index int, rgb color.Color) error {
    if index >= s.Length {
        return errors.New("index out of bounds")
    }

    s.Strip.Leds(0)[index] = colorToBytes(rgb)
    return nil
}

func (s *Strip) Clear(rgb color.RGBA) {
    for i := 0; i < s.Length; i++ {
        s.Set(i, rgb)
    }
}

func (s *Strip) Render() error {
    return s.Strip.Render()
}

func (s *Strip) Close() {
    s.Strip.Fini()
}

type Matrix struct {
    strip *Strip
    Width, Height int
}

func NewMatrix(width, height, pin int) (*Matrix, error) {
    strip, err := NewStrip(pin, width*height)
    if err != nil {
        return nil, err
    }
    return &Matrix{strip, width, height}, nil
}

func (m *Matrix) Set(x, y int, rgb color.Color) error {
    var index int
    if x >= m.Width || y >= m.Height {
        return errors.New("index out of bounds")
    }

    if y%2 == 0 {
        index = y*m.Width + x
    } else {
        index = (y+1)*m.Width - (x+1)
    }
    m.strip.Set(index, rgb)
    return nil
}

func (m *Matrix) Clear(rgb color.RGBA) {
    m.strip.Clear(rgb)
}

func (m *Matrix) Render() error {
    return m.strip.Render()
}

func (m *Matrix) Close() {
    m.strip.Close()
}
