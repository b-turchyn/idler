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

func Cost(initialCost uint64, count uint64) uint64 {
  return uint64(float64(initialCost) * math.Pow(coefficient, float64(count)))
}

