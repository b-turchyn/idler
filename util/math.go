package util

import (
  "math"
)

const coefficient = 1.15

func Max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func Min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func Within(min, val, max int) int {
  return Max(min, Min(val, max))
}

func Cost(initialCost uint64, count uint64) uint64 {
  return uint64(float64(initialCost) * math.Pow(coefficient, float64(count)))
}

