
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
    change = make(chan bool)
    matrix   *neopixel.Matrix

    BLACK = color.RGBA{0, 0, 0, 0}
    WHITE = color.RGBA{255, 255, 255, 255}

    fn int
    anim *unpack.Anim
)

const (
    WIDTH, HEIGHT = 15, 11
)

func event(event lirc.Event) {
    if event.Repeat > 0 {
        return
    }

    newanim, err := unpack.Search(event.Button)
    if err != nil {
        if unpack.IsNotFound(err) {
            log.Printf("Key %v not registered\n", event.Button)
            return
        }
        log.Println(err)
        return
    }

    anim = newanim
    change <- true
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
    ir.Handle("", "", event)
    matrix.Set(0, 0, color.White)
    matrix.Set(WIDTH-1, 0, color.White)
    matrix.Set(0, HEIGHT-1, color.White)
    matrix.Set(WIDTH-1, HEIGHT-1, color.White)
    matrix.Render()

    go ir.Run()
    log.Println("running")
    run()
}

func run() {
    <-change
    out:
    for {
        select {
        case <-change:
            fn = 0
        default:
        }

        marcher:
        for fn < len(anim.Frames) {
            fps := time.NewTimer(time.Second/time.Duration(anim.Fps))
            select {
            case <-change:
                continue out
            default:
            }

            frame := anim.Frames[fn]

            if len(frame) == 0 {
                fps.Stop()
                timer := time.NewTimer(anim.LoopDelay)

                select {
                case <-timer.C:
                    fn++
                    timer.Stop()
                    continue marcher
                case <-change:
                    timer.Stop()
                    continue out
                }
            }

            frame.Draw(matrix)
            matrix.Render()
            fn++

            select {
            case <-fps.C:
                fps.Stop()
                continue marcher
            case <-change:
                fps.Stop()
                continue out
            }
        }
        if anim.Loop {
            fn = 0
            if anim.LoopDelay > 0 {
                timer := time.NewTicker(time.Duration(anim.LoopDelay))
                select {
                case <-timer.C:
                    timer.Stop()
                    continue out
                case <-change:
                    timer.Stop()
                    continue out
                }
            }
        }
    }
}
