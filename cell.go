package main

import "net/http"
import "io/ioutil"
import "image/color"
import "math/rand"
import "strconv"
import "strings"

type cell struct {
	client       *http.Client
	url          string
	index        int
	baseColor    color.RGBA
	currentColor uint32
	fadeCycles   int
	lastError    error
	currentCycle int
	fetching     bool
}

var errColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}

func fadeCycles() int {
	return rand.Int() % 100
}

type cycler interface {
	cycle()
	fetch(url string)
}

func newCell(client *http.Client, url string, i int) cell {
	black := color.RGBA{R: 0, G: 0, B: 0, A: 0xFF}
	return cell{url: url, index: i, baseColor: black, currentColor: black, fadeCycles: 0, currentCycle: 0, fetching: false, lastError: nil, client: client}
}

func (c *cell) cycle() {
	if c.currentCycle == 0 && !c.fetching {
		go func() {
			c.fetch()
		}()
	} else if c.fetching {
		// nothing
	} else {
		alpha := uint8(0xFF * (float64(c.currentCycle) / float64(c.fadeCycles)))
		nrgba := color.NRGBA{R: c.baseColor.R, G: c.baseColor.G, B: c.baseColor.B, A: alpha}
		fr, fg, fb, fa := nrgba.RGBA()
		c.currentColor = (fr << 16) + (fg << 8) + fb
		c.currentCycle--
	}
}

func (c *cell) fetch() {
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		c.lastError = err
		c.baseColor = errColor
		c.fadeCycles = fadeCycles()
		c.fetching = false
		return
	}
	req.Header.Set("Connection", "close")
	req.Close = true
	resp, err := c.client.Get(c.url)
	if err != nil {
		c.lastError = err
		c.baseColor = errColor
		c.fadeCycles = fadeCycles()
		c.fetching = false
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.lastError = err
		c.baseColor = errColor
		c.fadeCycles = fadeCycles()
		c.fetching = false
		return
	}
	err = resp.Body.Close()
	if err != nil {
		c.lastError = err
		c.baseColor = errColor
		c.fadeCycles = fadeCycles()
		c.fetching = false
		return
	}
	baseColor, err := strconv.ParseInt(strings.Trim(string(body[:]), "\n"), 16, 32)
	if err != nil {
		c.lastError = err
		c.baseColor = errColor
		c.fadeCycles = fadeCycles()
		c.fetching = false
		return
	}
	r := uint8(baseColor & 0xFF0000 >> 16)
	g := uint8(baseColor & 0xFF00 >> 8)
	b := uint8(baseColor & 0xFF)
	c.baseColor = color.RGBA{R: r, G: g, B: b, A: 0xFF}
	c.fadeCycles = fadeCycles()
	c.currentCycle = c.fadeCycles
	c.lastError = nil
	c.fetching = false
}
