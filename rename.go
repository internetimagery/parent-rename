// Rename files to that of their containing folder.
// Also ignore anything in square brackets, for
// compatability with TagSpaces.
// Version 1
package main

import (
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
  "sync"
  "regexp"
  "strconv"
  "math"
)

// Directory prepped for renaming
type Directory struct {
  Name string // Name of Folder
  Files map[string]string // Map old names to new names for rename
}

// Directory Factory
func NewDirectory(path string) (*Directory) {
  dir := new(Directory)
  dir.Name = filepath.Base(path)
  dir.Files = make(map[string]string)
  for _, f := range getfiles(path) {
    dir.Files[f] = ""
  }
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
func getfiles(path string) ([]string) {
  if paths, err := ioutil.ReadDir(path); err == nil {
    result := make([]string, 0)
    for _, p := range paths {
      if !p.IsDir() {
        result = append(result, p.Name())
      }
    }
    return result
  } else {
    panic(err)
  }
}

// Validate files in a folder match our naming scheme
// For files that do not match, add their corresponding scheme name
func validate(dir *Directory) {
  // Quote our base name to use in regexp
  base := regexp.QuoteMeta(dir.Name)
  reg1 := regexp.MustCompile(base + "_(\\d+)")
  // Run through our files and pick out only files that are valid
  count := 0.0
  for k, _ := range dir.Files {
    match := reg1.FindStringSubmatch(k)
    if len(match) > 0 {
      num, err := strconv.ParseFloat(match[1], 64)
      if err != nil { panic(err) }
      count = math.Max(num, count)
    } else {
      // Add a value to our file to mark it as a rename candidate
      // This is a temporary name
      // We're using the actual name just incase something goes wrong
      // and we try to rename it, we'll just rename it to the same name. *shrug*
      dir.Files[k] = k
    }
  }
  countInt := int(count)
  // Now we have a maximum file number, build new names that incriment it
  reg2 := regexp.MustCompile("\\[.*\\]") // Keep tags for Tagspaces
  for k, v := range dir.Files {
    if v != "" {
      countInt += 1
      tags := reg2.FindString(k)
      name := dir.Name + "_" + strconv.Itoa(countInt) + tags + filepath.Ext(k)
      dir.Files[k] = name
    }
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
    // Kick off our validation and renaming
    if len(dirs) > 0 {
      var hold sync.WaitGroup
      for k, v := range dirs {
        hold.Add(1)
        go func(path string, dir *Directory){
          defer hold.Done()
          // Validate and build our rename candidates
          validate(dir)
          for old, _new := range dir.Files {
            if _new != "" {
              // Build our absolute paths, and rename files!
              fmt.Printf("Rename: %s => %s ", old, _new)
              old_abs := filepath.Join(path, old)
              new_abs := filepath.Join(path, _new)
              err := os.Rename(old_abs, new_abs)
              if err == nil {
                fmt.Printf("OK\n")
              } else {
                fmt.Printf("FAILED\n")
              }
            }
          }
        }(k, v)
      }
      hold.Wait()
      fmt.Println("Done!")
      return
    }
  }
  // No args or something, print help menu
  fmt.Println("Usage: rename <directory> ...")
}
