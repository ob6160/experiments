/**
 * Heapsort Implementation in Golang for CS1FC16.
 * @author Oliver Barnwell <qb011080@reading.ac.uk>
 */

package main

import (
  "fmt"
  "time"
  "math"
  "math/rand"
)

/**
 * Global Variables.
 */

/**
 * Helper Functions.
 */

/**
 * Returns the array index of the specified heap parent.
 */
func iParent(i int) int {
  return int(math.Floor(float64((i - 1) / 2)))
}

/**
 * Returns the array index of the specified left heap child.
 */
func iLeftChild(i int) int {
  return 2 * i + 1
}

func swap(a []float64, i, j int) {
  a[i], a[j] = a[j], a[i]
}

/**
 * Completes a heapsort on the provided slice.
 *
 * a []float64
 *   The slice to be sorted.
 */
func heapsort(a []float64) {
  count := len(a)

  // Construct the initial heap.
  heapify(a, count)

  // Initialise the end as the last element in the heap (smallest).
  end := count - 1

  // Until the end escapes the bounds of the slice iterate.
  for(end > 0) {
    // Swap the last element with the first.
    swap(a, end, 0)

    // We now consider the new last element to be fixed, so we decrement end so it is not considered.
    end--

    // "Sift" the new first element into the correct position in the heap.
    siftDownIterative(a, 0, end)
  }
}

/**
 * Create the initial heap structure.
 *
 * a []float64
 *   Slice to be heapified.
 * count int
 *   Max size of the heap
 */
func heapify(a []float64, count int) {
  start := iParent(count - 1)
  for(start >= 0) {
    siftDownIterative(a, start, count - 1)
    start--
  }
}

/**
 * Repairs the provided heap where the root element is at index start.
 *   This method is applied iteratively.
 *
 * a []float64
 *   The heap to be repaired (heapified).
 * start, end int
 *   The start and end indexes.
 *
 */
func siftDownIterative(a []float64, start, end int) {
  // Set the root index to the specified start.
  root := start

  // While the index of the root's left child is less than the specified end index iterate. Instead this could be implemented as a recursive method.
  for(iLeftChild(root) <= end) {
    // Calculate the children's indexes and initialise a variable to keep track of the largest child.
    left := iLeftChild(root)
    right := left + 1
    largestChild := root

    // Compare root with its left child.
    // If left child is greater choose it.
    if sortCompareLessThan(a[largestChild], a[left]) {
      largestChild = left
    }

    // Compare current largest with right child.
    // If right child is larger choose it.
    if right <= end && sortCompareLessThan(a[largestChild], a[right]) {
      largestChild = right
    }

    // Check if the root is the largest of it and its children.
    // If so we have achieved our goal of ordering the tree.
    if largestChild == root {
      return
    }

    // Otherwise, swap the largest child with the original root.
    swap(a, root, largestChild)

    /// Assign a new root for the next iteration.
    root = largestChild
  }
}

/**
 * Custom comparison function which returns the desired results when comparing against NaN and negative zero.
 *
 */
func sortCompareLessThan(a, b float64) bool {
  var result bool = false
  aNaN := math.IsNaN(a)
  aSign := math.Signbit(a)
  bSign := math.Signbit(b)

  // Base less than check.
  if a < b {
    result = true
  }

  // Check for NaN.
  if aNaN {
    result = true
  }

  // Check for negative zero.
  if (a == 0 && aSign) && (b == 0 && !bSign) {
    result = true
  }

  return result
}

/**
 * Repairs the provided heap where the root element is at index start.
 *   This method is applied recursively.
 *
 * a []float64
 *   The heap to be repaired (heapified).
 * start, end int
 *   The start and end indexes.
 *
 */
func siftDownRecursive(a []float64, start, end int) {
    // Calculate the children's indexes and initialise a variable to keep track of the largest child.
    left := iLeftChild(start)
    right := left + 1
    largestChild := start

    // Compare root with its left child.
    // If left child is greater choose it.
    if left <= end && sortCompareLessThan(a[largestChild], a[left]) {
      largestChild = left
    }

    // Compare current largest with right child.
    // If right child is larger choose it.
    if right <= end && sortCompareLessThan(a[largestChild], a[right]) {
      largestChild = right
    }

    // Check if the root is the largest of it and its children.
    // If so we have achieved our goal of ordering the tree.
    if largestChild != start {
      // Otherwise, swap the largest child with the original root.
      swap(a, start, largestChild)
      siftDownRecursive(a, largestChild, end)
    }
}

func testRunner(arr []float64, message string) {
  fmt.Printf("%s\n", message)
  heapsort(arr)
  fmt.Printf("%v\n", arr);
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  testArray := []float64{}
  testRunner(testArray, "Sorting an empty list.")

  testArray = []float64{1.0}
  testRunner(testArray, "Sorting a list with one element.")

  testArray = []float64{2.0, 1.0}
  testRunner(testArray, "Sorting a list with two elements.")

  testArray = []float64{2.0, 4.0, 1.0, 5.0}
  testRunner(testArray, "Sorting a list with a general number of elements")

  testArray = []float64{6.0, 1.0, 2.0, 1.0, 3.0, 100.0, 1.0}
  testRunner(testArray, "Sorting a list where some elements are the same")

  // math.Inf(sign int) -> positive infinty if sign >= 0 negative infinty if < 0
  testArray = []float64{1.0, 100.0, math.Inf(1.0), 2.0, math.Inf(-1.0)}
  testRunner(testArray, "Sorting a list where some of the elements are infinity and minus infinity")

  zero := float64(0)
  neg_zero := -zero
  testArray = []float64{-100.0,  neg_zero, -50.0, -1.0, 3.0, 0.0, 2.0, neg_zero, 1.0, neg_zero, 5.0, 0.0, 0.0, 1.0}
  testRunner(testArray, "Sorting a list where some of the elements are minus zero")

  testArray = []float64{math.NaN(), 6.0, 1.0, math.NaN()}
  testRunner(testArray, "Sorting a list where some elements are NaN")

  var example_nan float64 = math.Sqrt(-1)
  var example_nan_two float64 = math.Sin(math.Inf(1.0))
  testArray = []float64{-1.0, 6.0, 2.0, math.NaN(), 2.0, 1000.0, example_nan, example_nan_two}
  testRunner(testArray, "Sorting a list where some elements are NaN with different mantissas.")

  testArray = make([]float64, 100)
  for i := 0; i < 100; i++ {
    testArray[i] = rand.Float64() * 100
  }
  testRunner(testArray, "Sorting a list of 100 pseudorandom numbers")
}
