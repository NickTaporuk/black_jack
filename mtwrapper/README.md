### How to use
```go

package main

import (
    "fmt"
    "project/jackpot" // Import the jackpot package
)

func main() {
    // Generate a random number
    randomNum := jackpot.GenerateRandom(1, 100)
    fmt.Printf("Random Number: %d\n", randomNum)

    // Simulate jackpot drop
    if jackpot.DropJackpot(10000) {
        fmt.Println("🎉 Jackpot Dropped!")
    } else {
        fmt.Println("❌ No Jackpot This Time.")
    }
}

```