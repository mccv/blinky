package main

import "fmt"
import "time"
import "net/http"
import "github.com/jgarff/ws2811"

func main() {
	numCells := 25
	cells := make([]cell, numCells)
	ledBitmap := make([]uint32, numCells)
	client := &http.Client{Timeout: 2 * time.Second}
	ws2811.Init(18, numCells, 255)
	for i := 0; i < len(cells); i++ {
		cells[i] = newCell(client, "http://104.196.242.214/api/logs", i)
	}
	sleepTime := 30 * time.Millisecond
	cycles := 0
	errs := 0
	for {
		cycles++
		if cycles%100 == 0 {
			fmt.Printf("%d cycles, %d errs\n", cycles, errs)
		}
		for j := 0; j < len(cells); j++ {
			cells[j].cycle()
			if cells[j].lastError != nil {
				errs++
				cells[j].lastError = nil
			}
			ledBitmap[j] = cells[j].currentColor
		}
		ws2811.SetBitmap(ledBitmap)
		ws2811.Render()
		time.Sleep(sleepTime)
	}
}
