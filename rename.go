// Rename files to that of their containing folder.
// Also ignore anything in square brackets, for
// compatability with TagSpaces.
package main

import (
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
)

// Directory prepped for renaming
type Directory struct {
  Name string
  Files []os.FileInfo
}

// Directory Factory
func NewDirectory(path string) (*Directory) {
  dir := new(Directory)
  dir.Name = filepath.Base(path)
  dir.Files = getfiles(path)
  return dir
}

// Check if path is a directory
func isdir(path string) bool {
  if st, err := os.Stat(path); err == nil && st.IsDir() {
    return true
  }
  return false
}

// Get all files from a directory
func getfiles(path string) ([]os.FileInfo) {
  if paths, err := ioutil.ReadDir(path); err == nil {
    result := make([]os.FileInfo, 0)
    for _, p := range paths {
      if !p.IsDir() {
        result = append(result, p)
      }
    }
    return result
  } else {
    panic(err)
  }
}

func main()  {
  // Grab our cli arguments
  args := os.Args[1:]
  if len(args) > 0 {
    // Get our current dir
    cwd, _ := os.Getwd()
    // Store valid directories. Use map to eliminate doubleups
    dirs := make(map[string]*Directory)
    for _, p := range args {
      if abs := filepath.Join(cwd, p); isdir(abs) {
        dirs[abs] = NewDirectory(abs)
      }
    }
    fmt.Println("Directories", dirs)
  } else {
    // No args, print help menu
    fmt.Println("Usage: rename <directory> ...")
  }
}
