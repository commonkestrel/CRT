package main

import (
    "image/color"
    "log"
    "time"

    "example/CRT/neopixel"
    "example/CRT/unpack"

    "github.com/chbmuc/lirc"
)

var (
    rgb      = color.RGBA{255, 255, 255, 255}
    stop bool
    running  bool
    matrix   *neopixel.Matrix

    BLACK = color.RGBA{0, 0, 0, 0}
    WHITE = color.RGBA{255, 255, 255, 255}
)

const (
    WIDTH, HEIGHT = 15, 11
)

func change(event lirc.Event) {
    log.Println(event)
    if event.Repeat == 0 {
        if running {
            stop = true
            running = false
        }
        
        log.Println(stop)
        for stop {}

        switch event.Button {
        case "KEY_1":
            go Anim("eye.json")
        case "KEY_2":
            go Anim("clock.json")
        case "KEY_3":
            go Anim("spiral.json")
        case "KEY_4":
            go Anim("tv.json")
        case "KEY_5":
            go Anim("heart.json")
        case "KEY_6":
            matrix.Clear(WHITE)
        case "KEY_0":
            
        case "KEY_POWER":
            matrix.Clear(BLACK)
            matrix.Render()
        }
    }
}

func Anim(path string) {
    build, err := unpack.UnpackFile(path)
    if err != nil {
        log.Println(err)
        matrix.Clear(BLACK)
        matrix.Render()
        return
    }
    
    running = true
    if build.LoopDelay > 0 && build.Loop {
        build.Frames[len(build.Frames)-1].Draw(matrix)
        matrix.Render()
    }
    for {
        if build.LoopDelay > 0 && build.Loop {
            delay := time.NewTimer(build.LoopDelay)
            <-delay.C
        }

        for _, frame := range build.Frames {
            fps := time.NewTicker(time.Second / time.Duration(build.Fps))
            defer fps.Stop()

            if stop {
                stop = false
                return
            }

            frame.Draw(matrix)
            matrix.Render()

            <-fps.C
        }
        if !build.Loop {
            break
        }
    }
    running = false
}

func main() {
    var err error
    matrix, err = neopixel.NewMatrix(WIDTH, HEIGHT, 21)
    if err != nil {
        log.Fatal(err)
    }
    defer matrix.Close()

    ir, err := lirc.Init("/var/run/lirc/lircd")
    if err != nil {
        panic(err)
    }
    ir.Handle("", "", change)
    matrix.Set(5, 5, color.White)
    matrix.Render()

    log.Println("running")
    ir.Run()
}
