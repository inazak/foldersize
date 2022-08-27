package main

import (
  "flag"
  "fmt"
  "os"
  "path/filepath"
  "sort"
  "sync"
  "sync/atomic"

  "github.com/mattn/go-runewidth"
  "github.com/saracen/walker"
)

var usage =`
Usage:
  foldersize [OPTION] [PATTERN]

  Print total size of folders.
  Sort from largest to smallest size.

  OPTION:
    -m      ... display in MB
    -k      ... display in KB
`

var optionsMB    bool
var optionsKB    bool

func main() {

  flag.BoolVar(&optionsMB, "m", false, "display in MB")
  flag.BoolVar(&optionsKB, "k", false, "display in KB")
  flag.Parse()

  if len(flag.Args()) > 1 {
    fmt.Printf("%s", usage)
    os.Exit(1)
  }

  var pattern string

  if len(flag.Args()) == 1 {
    pattern = flag.Args()[0]
  } else {
    // if no argument is given, current directory is scanned
    pwd, _ := os.Getwd()
    pattern = filepath.Join(pwd, "*")
  }

  pathlist, err := GetFolderpathList(pattern)
  if err != nil {
    errorExit(err.Error())
  }

  sizemap := GetSizeMap(pathlist)

  sorted := []string{}
  maxlen_path := 0
  maxlen_size := 0

  for path, size := range sizemap {
    sorted = append(sorted, path)

    lp := runewidth.StringWidth(path)
    if lp > maxlen_path {
      maxlen_path = lp
    }

    ls := len(fmt.Sprintf("%d", size))
    if ls > maxlen_size {
      maxlen_size = ls
    }
  }

  sort.Slice(sorted, func(i, j int) bool {
    return sizemap[sorted[i]] > sizemap[sorted[j]]
  })

  for _, path := range sorted {
    size := sizemap[path] //path is space-filled next line
    path := runewidth.FillRight(path, maxlen_path)

    if optionsMB { size = (size / 1024 + 1 ) / 1024 + 1 }
    if optionsKB { size =  size / 1024 + 1 }

    // Output
    format := fmt.Sprintf("%%s  %%%dd\n", maxlen_size)
    fmt.Printf(format, path, size)
  }

  os.Exit(0)
}


func GetFolderpathList(filter string) (list []string, err error) {
  entries, err := filepath.Glob(filter)
  if err != nil {
    return nil, err
  }

  for _, entry := range entries {
    isdir, err := isDir(entry)
    if err != nil {
      return nil, err
    }
    if isdir {
      list = append(list, entry)
    }
  }

  return list, nil
}


func isDir(p string) (bool, error) {
  f, err := os.Stat(p)
  if err != nil {
    return false, err
  }
  return f.Mode().IsDir(), nil
}

func GetSizeMap(pathlist []string) map[string]int64 {

  sizemap := make(map[string]int64)

  var wg sync.WaitGroup
  var mu sync.Mutex

  for _, path := range pathlist {
    wg.Add(1)

    go func(p string, s map[string]int64) {
      defer wg.Done()

      size, err := GetFolderSize(p)
      if err != nil {
        size = -1 //FIXME
      }

      mu.Lock()
      s[p] = size
      mu.Unlock()
    }(path, sizemap) //goroutine

  } //for

  wg.Wait()

  return sizemap
}

func GetFolderSize(folderpath string) (int64, error) {
  var size int64

  err := walker.Walk(folderpath, func(_ string, info os.FileInfo) error {
    if ! info.IsDir() {
      atomic.AddInt64(&size, info.Size())
    }
    return nil
  })

  return size, err
}

func errorExit(f string, v ...interface{}) {
  fmt.Fprintf(os.Stderr, f, v...)
  os.Exit(1)
}

