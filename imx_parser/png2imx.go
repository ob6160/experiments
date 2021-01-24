package main

import (
  "fmt"
  _"github.com/ob6160/imxutils/imx"
  _"flag"
  _"strings"
  "image"
  "image/color"
  _"image/png"
  "os"
  "strings"
)

func main() {
    infile, err := os.Open(os.Args[1])
    if err != nil {
        // replace this with real error handling
        panic(err)
    }
    defer infile.Close()

    // Decode will figure out what type of image is in the file on its own.
    // We just have to be sure all the image packages we want are imported.
    src, _, err := image.Decode(infile)
    if err != nil {
        // replace this with real error handling
        panic(err)
    }

    // Richard's Neural Net Header Data------------------
    fmt.Println(fmt.Sprintf("%d %d", 784, 10), "%.1f %.0f %.3f")

    dataCellCount := 784
    labelCellCount := 10
    // Input labels
    for i := 0; i < (dataCellCount + labelCellCount); i++ {
      fmt.Printf("P%d ",i)
    }
    fmt.Print("\n")

    // Scale min
    for i := 0; i < (dataCellCount + labelCellCount); i++ {
      fmt.Printf("0 ")
    }
    fmt.Print("\n")

    // Scale max
    for i := 0; i < (dataCellCount); i++ {
      fmt.Printf("255 ")
    }

    for i := 0; i < labelCellCount; i++ {
      fmt.Printf("1 ")
    }
    fmt.Print("\n")

    // Create a new grayscale image
    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y
    for y := 0; y < h; y++ {
      var acc string = ""
      for x := 0; x < w; x++ {
        acc += fmt.Sprintf("%s\t",strings.Trim(fmt.Sprintf("%d",color.GrayModel.Convert(src.At(x, y))), "{}"))
      }
      fmt.Printf("%s", acc)
    }


}
