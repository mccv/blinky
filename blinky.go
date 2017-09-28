package main

import "fmt"
import "time"
import "net/http"

func main() {
	var cells [25]cell
	var ledBitmap [25]uint32
	client := &http.Client{Timeout: 2 * time.Second}
	for i := 0; i < len(cells); i++ {
		cells[i] = newCell(client, "http://104.196.242.214/api/logs", i)
	}
	sleepTime := 30 * time.Millisecond
	for i := 0; i < 100; i++ {
		for j := 0; j < len(cells); j++ {
			cells[j].cycle()
			if cells[j].lastError != nil {
				fmt.Printf("cell %d error: %v", j, cells[j].lastError)
				cells[j].lastError = nil
				fmt.Print("E")
			} else {
				fmt.Print(".")
			}
			ledBitmap[j] = cells[j].currentColor
		}
		fmt.Println("")
		time.Sleep(sleepTime)
	}
}
