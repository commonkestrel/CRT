package unpack

import (
	"encoding/json"
	"errors"
	"image/color"
	"os"
	"strings"
	"time"

	"example/CRT/neopixel"
)

var errNotFound = errors.New("command not found for key.")
var keys = map[string]string{"power": "KEY_POWER", "volumeup": "KEY_VOLUMEUP", "stop": "KEY_STOP", "previous": "KEY_PREVIOUS", "playpause": "KEY_PLAYPAUSE", "next": "KEY_NEXT", "down": "KEY_DOWN", "volumedown": "KEY_VOLUMEDOWN", "up": "KEY_UP", "equal": "KEY_EQUAL", "start": "BTN_START", "1": "KEY_1", "2": "KEY_2", "3": "KEY_3", "4": "KEY_4", "5": "KEY_5", "6": "KEY_6", "7": "KEY_7", "8": "KEY_8", "9": "KEY_9", "0": "KEY_0"}

const (
    WIDTH, HEIGHT = 15, 11
)

type Anim struct {
    Frames []Frame
    Fps int
    Loop bool
    LoopDelay time.Duration
}

type frames struct {
    Frames []Frame
    Fps int
}

type in struct {
    Name string
    Filename string
    Key string
    Loop bool
    LoopDelay int //In milliseconds
}

type Frame [][]color.RGBA

func (f *Frame) Draw(mat *neopixel.Matrix) {
    for x, col := range *f {
        for y, rgb := range col {
            mat.Set(x, mat.Height-(y+1), rgb)
        }
    }
}

func readindex() ([]in, error) {
    file, err := os.ReadFile("./store/index.json")
    if err != nil {
        return nil, err
    }

    var unmarshaled []in
    err = json.Unmarshal(file, &unmarshaled)
    if err != nil {
        return nil, err
    }

    return unmarshaled, nil
}

func Search(key string) (*Anim, error) {
    index, err := readindex()
    if err != nil {
        return nil, err
    }

    var found *in = nil
    for _, c := range index {
        if keys[strings.ToLower(c.Key)] == key {
            found = &c
        }
    }
    if found == nil {
        return nil, errNotFound
    }

    path := found.Filename

    var unmarshaled frames
    file, err := os.ReadFile("./store/" + path)
    if err != nil {
        return nil, err
    }
    json.Unmarshal([]byte(file), &unmarshaled)

    unpacked := &Anim{unmarshaled.Frames, unmarshaled.Fps, found.Loop, time.Duration(found.LoopDelay)*time.Millisecond}

    return unpacked, nil
}

func IsNotFound(err error) bool {
    return err == errNotFound
}
