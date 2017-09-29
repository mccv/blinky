package main

import "net/http"
import "io/ioutil"
import "math/rand"
import "strconv"
import "strings"

type cell struct {
	client       *http.Client
	url          string
	index        int
	baseRed      uint8
	baseGreen    uint8
	baseBlue     uint8
	fadeCycles   int
	lastError    error
	currentColor uint32
	currentCycle int
	fetching     bool
}

func fadeCycles() int {
	return 50 + rand.Int() % 100
}

type cycler interface {
	cycle()
	fetch(url string)
}

func (c *cell) setError(err error) {
	c.lastError = err
	c.baseRed = 255
	c.baseGreen = 0
	c.baseBlue = 0
	c.setCurrentColor(c.baseRed, c.baseGreen, c.baseBlue)
	c.fadeCycles = fadeCycles()
	c.fetching = false
}

func newCell(client *http.Client, url string, i int) cell {
	return cell{url: url, index: i, baseRed: 0, baseGreen: 0, baseBlue: 0, fadeCycles: 0, currentCycle: 0, fetching: false, lastError: nil, client: client}
}

func (c *cell) setCurrentColor(red uint8, green uint8, blue uint8) {
	c.currentColor = (uint32(green) << 16) + (uint32(red) << 8) + uint32(blue)
}

func (c *cell) cycle() {
	if c.currentCycle == 0 && !c.fetching {
		go func() {
			c.fetch()
		}()
	} else if c.fetching {
		// nothing
	} else {
		alpha := float64(c.currentCycle) / float64(c.fadeCycles)
		fr := uint8(alpha * float64(c.baseRed))
		fg := uint8(alpha * float64(c.baseGreen))
		fb := uint8(alpha * float64(c.baseBlue))
		// LED colors are GRB
		c.setCurrentColor(fr, fg, fb)
		c.currentCycle--
	}
}

func (c *cell) fetch() {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		c.setError(err)
		return
	}
	req.Header.Set("Connection", "close")
	req.Close = true
	resp, err := c.client.Get(c.url)
	if err != nil {
		c.setError(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.setError(err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		c.setError(err)
		return
	}
	baseColor, err := strconv.ParseInt(strings.Trim(string(body[:]), "\n"), 16, 32)
	if err != nil {
		c.setError(err)
		return
	}
	c.baseRed = uint8(baseColor & 0xFF0000 >> 16)
	c.baseGreen = uint8(baseColor & 0xFF00 >> 8)
	c.baseBlue = uint8(baseColor & 0xFF)
	c.setCurrentColor(c.baseRed, c.baseGreen, c.baseBlue)
	c.fadeCycles = fadeCycles()
	c.currentCycle = c.fadeCycles
	c.lastError = nil
	c.fetching = false
}
