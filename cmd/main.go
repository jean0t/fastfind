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


func main() {
    cli_flags()
    flag.Parse()

    var wg sync.WaitGroup
    var fileChan chan string = make(chan string)

    wg.Add(1)
    go func() {
        defer wg.Done()
        err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error accessing path %q: %v", path, err)
                return nil
            }

            if !includeDirs && d.IsDir() {
                return nil
            }

            if pattern != "" && !strings.Contains(d.Name(), pattern) {
                return nil
            }

            if extFilter != "" && !strings.HasSuffix(d.Name(), extFilter) {
                return nil
            }

            fileChan <- path
            return nil
        })
        if err != nil {
            fmt.Fprintf(os.Stderr, "Walk error: %v\n", err)
        }
    }()

    go func() {
        wg.Wait()
        close(fileChan)
    }()

    for f := range fileChan {
        fmt.Println(f)
    }

}
