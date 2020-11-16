package base62

import (
  "errors"
  "math"
  "strings"
)

// Base
const (
  alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  length   = uint64(len(alphabet))
)

// For encoding 
func Encode(number uint64) string { // Return string 
  var encodedBuilder strings.Builder // Create a string builder variable 
  encodedBuilder.Grow(11) // Grow function 

  for ; number > 0; number = number / length {
     encodedBuilder.WriteByte(alphabet[(number % length)])
  }

  return encodedBuilder.String() // We return the encodedBuilder as String 
}

// For decoding 
func Decode(encoded string) (uint64, error) { // Return uint64 and error
  var number uint64 // Create a number variable as uint64 

  for i, symbol := range encoded {
     alphabeticPosition := strings.IndexRune(alphabet, symbol)

     if alphabeticPosition == -1 {
        return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
     }
     number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
  }

  return number, nil // return the uint64 and error 
}

/* ('Sami Ghasemi) */
