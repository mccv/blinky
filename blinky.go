package main

import "fmt"
import "time"
import "net/http"
import "github.com/jgarff/ws2811"

func main() {
	numCells := 15
	cells := make([]cell, numCells)
	ledBitmap := make([]uint32, numCells)
	client := &http.Client{Timeout: 2 * time.Second}
	ws2811.Init(18, 15, 255)
	for i := 0; i < len(cells); i++ {
		cells[i] = newCell(client, "http://104.196.242.214/api/logs", i)
	}
	sleepTime := 30 * time.Millisecond
	for i := 0; i < 10000; i++ {
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
		ws2811.SetBitmap(ledBitmap)
		ws2811.Render()
		fmt.Println("")
		time.Sleep(sleepTime)
	}
}
