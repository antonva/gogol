package plugins

import (
    "os"
    "fmt"
)

func bots() string{
    fmt.Println(os.Args[0])
    return os.Args[0]
}
