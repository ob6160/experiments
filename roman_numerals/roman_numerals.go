package main

import (
  "fmt"
  "math"
)

func digitAtPos(number, pos int) int {
  return number / int(math.Pow(10, float64(pos))) % 10
}

func toNumeral(number, pow int) string {
  if number == 0 {
    return ""
  }

  base := 'I'
  mid := 'V'
  max := 'X'
  switch pow {
  case 1:
    base = 'X'
    mid = 'L'
    max = 'C'
  case 2:
    base = 'C'
    mid = 'D'
    max = 'M'
  case 3:
    base='M'
  }

  sizes := [10]int{1, 2, 3, 2, 1, 2, 3, 4, 2, 1}
  size := sizes[number-1]
  b := make([]rune, size)

  if number <= 3 {
    for i := 0; i < size; i++ {
        b[i] = base
    }
  } else if number == 4 {
    b[0] = base
    b[1] = mid
  } else if number <= 8 {
    b[0] = mid
    if number > 5 {
      for i := 1; i < size; i++ {
        b[i] = base
     }
   }
  } else if number == 9 {
    b[0] = base
    b[1] = max
  }

  return string(b)
}

func main() {
  toConvert := 678
  fmt.Println("Converting: ", toConvert)
  for i := 0; i < 4; i++ {
    currentDigit := digitAtPos(toConvert, i)
    numeral := toNumeral(currentDigit, i)
    fmt.Println(currentDigit, i, numeral)
  }

}
