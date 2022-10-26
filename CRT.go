package main

import (
    "image/color"
    "log"
    "time"
    "os"

    "example/CRT/neopixel"

    "github.com/chbmuc/lirc"
)

var rgb = color.RGBA{255, 255, 255, 255}

func change(event lirc.Event) {
    switch event.Button {
    case "KEY_1":
        rgb = color.RGBA{255, 0, 0, 255}
    case "KEY_2":
        rgb = color.RGBA{255, 127, 80, 255}
    case "KEY_3":
        rgb = color.RGBA{255, 255, 0, 255}
    case "KEY_4":
        rgb = color.RGBA{0, 255, 0, 255}
    case "KEY_5":
        rgb = color.RGBA{0, 0, 255, 255}
    case "KEY_6":
        rgb = color.RGBA{255, 0, 255, 255}
    case "KEY_0":
        rgb = color.RGBA{255, 255, 255, 255}
    case "KEY_POWER":
        os.Exit(0)
    }
}

func main() {
    matrix, err := neopixel.NewMatrix(15, 11, 21)
    if err != nil {
        log.Fatal(err)
    }
    defer matrix.Close()

    ir, err := lirc.Init("/var/run/lirc/lircd")
    if err != nil {
        panic(err)
    }
    ir.Handle("", "", change)

    go ir.Run()
    for {
        for x := 0; x < matrix.Width; x++ {
            matrix.Clear()
            for y := 0; y < matrix.Width; y++ {
                matrix.Set(x, y, rgb)
            }
            matrix.Render()
            time.Sleep(time.Second / 4)
        }
    }
}
