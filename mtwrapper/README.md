### How to use
```go

package main

import (
    "fmt"
    "project/mtwrapper" // Import the mtwrapper package
)

func main() {
    // Generate a random number
    randomNum := mtwrapper.GenerateRandom(1, 100)
    fmt.Printf("Random Number: %d\n", randomNum)

    // Simulate jackpot drop
    if mtwrapper.DropJackpot(10000) {
        fmt.Println("🎉 Jackpot Dropped!")
    } else {
        fmt.Println("❌ No Jackpot This Time.")
    }
}

```
