package imx

import (
  "fmt"
  "os"
  "bytes"
	"encoding/binary"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

/**
 * Checks for the validity of a given IMX magic number.
 */
func magicNumCheck(number []byte) bool {
  valid := false
  valid = ((number[0] | number[1]) == 0)
  return valid
}

var typeMap map[byte]string = map[byte]string{
  0x08: "ubyte",
  0x09: "byte",
  0x0B: "short",
  0x0C: "int",
  0x0D: "float",
  0x0E: "double",
}

func getDataType(dt byte) string {
  return typeMap[dt]
}

func LoadFile(path *string) (int, [][]byte) {
  f, err := os.Open(*path) // Open file for reading
  check(err)

  /************************************************************************
   * IMX File Format.
   *
   * Read the magic number, defined as four bytes at start of file.
   * First two bytes are *always* zero
   *
   * Third byte represents data type stored in data structure.
   * - 0x08: unsigned byte
   * - 0x09: signed byte
   * - 0x0B: short (2 bytes)
   * - 0x0C: int (4 bytes)
   * - 0x0D: float (4 bytes)
   * - 0x0E: double (8 bytes)
   *
   * Fourth byte stores the dimensionality of the structure
   * - 1: Vector
   * - 2: Matrices
   * If there are > 2 dimensions they are accumulated in matrix columns.
   *
   * Following this are dimension sizes.
   * Stored as four byte integers in BigEndian format.
   *
   * Finally data is stored as a series of unsigned bytes.
   *
   ************************************************************************/

  // Read in magic number and associated data.
  magicNumber := make([]byte, 4)
  magicNumBytesCheck, err := f.Read(magicNumber)
  check(err)
  if magicNumBytesCheck != 4 || !magicNumCheck(magicNumber) {
    panic("File Format Invalid")
  }
  // TODO: Handle data of different types
  // dataType := getDataType(magicNumber[2])
  noDimensions := magicNumber[3] * 4

  // Read in dimension sizes.
  dimensionality := make([]byte, noDimensions)
  dCount, err := f.Read(dimensionality)
  check(err)
  if dCount != int(noDimensions) {
    panic("Unexpected dimension count")
  }

  // Convert dimension bytes into integer representation.
  intDimensions := make([]uint32, noDimensions/4)
  r := bytes.NewReader(dimensionality)
  if err := binary.Read(r, binary.BigEndian, &intDimensions); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

  // Get number of rows
  rows := intDimensions[0]

  // Get column sizes so we can pack into a 1d format
  var columnsPacked int = 1
  for i := 1; i < len(intDimensions); i++ {
    columnsPacked *= int(intDimensions[i])
  }

  // Define appropriate storage for data.
  a := make([][]byte, rows)
  for i := range a {
    a[i] = make([]byte, columnsPacked)
  }

  // Now read the data one row at a time.
  for i := 0; i < int(rows); i++ {
    f.Read(a[i])
  }

  // TODO: Support verbose output (look into Go logging)
  // fmt.Printf("Magic Number Valid: %d\nData Type: %s\nDimension Count: %d\nDimensions:%d\n", magicNumber, dataType, len(intDimensions), intDimensions)

  return columnsPacked, a
}
