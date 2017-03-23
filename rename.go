// Rename files to that of their containing folder.
// Also ignore anything in square brackets, for
// compatability with TagSpaces.
package main

import (
  "fmt"
  "os"
  "path/filepath"
)

// Check if path is a directory
func isdir(path string) bool {
  if st, err := os.Stat(path); err == nil && st.IsDir() {
    return true
  }
  return false
}

func main()  {
  // Grab our cli arguments
  args := os.Args[1:]
  if len(args) > 0 {
    // Get our current dir
    cwd, _ := os.Getwd()
    for _, p := range args {
      abs := filepath.Join(cwd, p)
      if isdir(abs) {
        fmt.Println("OK", abs)
      }
    }
  } else {
    // No args, print help menu
    fmt.Println("Usage: rename <directory> ...")
  }
}
