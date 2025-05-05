package main


import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

var (
    rootDir string = "."
    pattern string = ""
    includeDirs bool = false
    extFilter string = ""
)


func cli_flags() {
    flag.StringVar(&rootDir, "path", ".", "Root directory to start searching")
    flag.StringVar(&pattern, "name", "", "Filename contains this pattern")
    flag.StringVar(&extFilter, "ext", "", "Filter by extension (e.g.: .go, .txt, .py)")
    flag.BoolVar(&includeDirs, "dirs", false, "Include directories in search results")
}

func walkDir(path string, wg *sync.WaitGroup, fileChan chan<- string) {
    defer wg.Done()

    entries, err := os.ReadDir(path)
    if err != nil {
        return
    }

    for _, entry := range entries {
        fullPath := filepath.Join(path, entry.Name())

        if entry.IsDir() {
            wg.Add(1)
            go walkDir(fullPath, wg, fileChan)
        } else if pattern != "" && strings.Contains(entry.Name(), pattern) {
            fileChan <- fullPath
        } else if extFilter != "" && strings.HasSuffix(entry.Name(), extFilter) {
            fileChan <- fullPath
        }

    }
}


func main() {
    cli_flags()
    flag.Parse()

    var wg sync.WaitGroup
    var fileChan chan string = make(chan string, 100)

    go func() {
        for file := range fileChan {
            fmt.Println(file)
        }
    }()

    wg.Add(1)
    go walkDir(rootDir, &wg, fileChan)

    wg.Wait()
    close(fileChan)
}
