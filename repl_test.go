package main

import (
  "testing"
)
func TestCleanInput(t *testing.T)  {
  cases := []struct {
    input string
    expected []string
  }{
    {
      input: " hello world ",
      expected: []string{"hello","world"},
    },
    {
      input: "Charmander Bulbasaur PIKACHU",
      expected: []string{"charmander","bulbasaur", "pikachu"},
    },

  }
  //The excecution loop
  for _, c := range cases{
    actual := CleanInput(c.input)
    
    if len(actual) != len(c.expected) {
      t.Errorf("For input %q: expected slice of length %d, got %d", c.input, len(c.expected), len(actual))
      continue
    }

    //Check each word in the slice
    for i := range actual{
      word := actual[i]
      expectedWord := c.expected[i]
      if word != expectedWord{
        t.Errorf("For input %q: at index %d, expected %q but got %q", c.input, i, expectedWord,word)
      }

    }
  }
}
