package main

import (
  "fmt"
  "github.com/ob6160/imxutils/imx"
  "github.com/fogleman/gg"
  "flag"
  "strings"
)

func main() {
  // Init graphics context
	const W = 28
	const H = 28
  dc := gg.NewContext(W, H)
	dc.SetRGB(100, 100, 100)
	dc.Clear()

  // Get path to imx file
  pathDataFlag  := flag.String("data", "", "Path to data file")
  pathLabelsFlag := flag.String("labels", "", "Path to data labels")
  dataRowsOut := flag.Int("rows", -1, "Rows to consider")
  dataRowShift := flag.Int("shift", 0, "Offset to start row output")

  // TODO: Add support for verbose logging output
  // verboseFlag := flag.String("v", "", "Verbose output")
  pngOut := flag.Bool("pngOut", false, "Output png files?")
  labelRange := flag.Int("labelRange", 10, "Size of label to convert to one hot")

  flag.Parse() // Parse flags
  dataCellCount, data := imx.LoadFile(pathDataFlag)
  labelCellCount, labels := imx.LoadFile(pathLabelsFlag)

  // Get row count
  dataRows  := len(data)
  labelRows := len(labels)

  if dataRows != labelRows {
    panic("Data/Label row mismatch")
  }

  if *dataRowsOut < 0 {
   dataRowsOut = &dataRows
  }


  // Richard's Neural Net Header Data------------------
  fmt.Println(fmt.Sprintf("%d %d", dataCellCount, *labelRange), "%.1f %.0f %.3f")

  // Input labels
  for i := 0; i < (dataCellCount + labelCellCount); i++ {
    fmt.Printf("P%d ",i)
  }
  fmt.Print("\n")

  // Scale min
  for i := 0; i < (dataCellCount + *labelRange); i++ {
    fmt.Printf("0 ")
  }
  fmt.Print("\n")

  // Scale max
  for i := 0; i < (dataCellCount); i++ {
    fmt.Printf("255 ")
  }

  for i := 0; i < *labelRange; i++ {
    fmt.Printf("1 ")
  }
  fmt.Print("\n")
  // End Richard's Neural Net Header Data---------------

  for i := *dataRowShift; i < (*dataRowsOut+*dataRowShift); i++ {
    logData := strings.Trim(fmt.Sprintf("%v", data[i]), "[]")
    logData = strings.Replace(logData, " ", "	", -1)
    oneHotLabel := make([]int, *labelRange)
    for j := range oneHotLabel {
      oneHotLabel[j] = 0
      if int(labels[i][0]) == j {
        oneHotLabel[j] = 1
      }
    }
    logLabels := strings.Trim(fmt.Sprintf("%v", oneHotLabel), "[]")
    logLabels = strings.Replace(logLabels, " ", " ", -1)
    fmt.Printf("%s\t%s\n", logData, logLabels)

    if *pngOut {
      dc.Clear()
      for i, point := range data[i] {
        x := i % 28
        y := i / 28
        dc.SetRGB255(int(point), int(point), int(point))
        dc.SetPixel(x, y)
      }
      dc.SavePNG(fmt.Sprintf("images/%d.png",i))
    }
  }
}
