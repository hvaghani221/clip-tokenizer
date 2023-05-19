package main

import (
	"flag"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pandodao/tokenizer-go"
)

var (
	SignLen int
	Freq    time.Duration
)

func initFunc() {
	flag.IntVar(&SignLen, "signlen", 60, "Length of the signature")
	flag.DurationVar(&Freq, "freq", time.Millisecond*17, "Clipboard fetch frequency")
	flag.Parse()
}

type Result[T any] struct {
	Value T
	Error error
}

func clipStream(d time.Duration) <-chan Result[string] {
	res := make(chan Result[string], 5)
	go func() {
		prev, _ := clipboard.ReadAll()
		timer := time.NewTicker(d)
		defer timer.Stop()
		for range timer.C {

			current, err := clipboard.ReadAll()
			if err != nil {
				res <- Result[string]{
					Error: err,
				}
			}
			if current == prev {
				continue
			}
			res <- Result[string]{
				Value: current,
			}
			prev = current
		}
	}()
	return res
}

type TokenResult struct {
	Tokens int
	Words  int
	Chars  int
	Sign   string
}

func tokeniseStream(clips <-chan Result[string]) <-chan Result[TokenResult] {
	res := make(chan Result[TokenResult], 5)
	go func() {
		cache := NewLRU[string, Result[TokenResult]](16)
		for clip := range clips {
			if clip.Error != nil {
				res <- Result[TokenResult]{Error: clip.Error}
				continue
			}
			if cached, found := cache.Get(clip.Value); found {
				res <- cached
				continue
			}
			tokens, err := tokenizer.CalToken(clip.Value)
			if err != nil {
				res <- Result[TokenResult]{Error: err}
				continue
			}
			sign := clip.Value
			if len(sign) > SignLen {
				sign = sign[:SignLen*3/4] + "..." + sign[len(sign)-SignLen/4-3:]
			}
			strings.Fields(clip.Value)
			value := Result[TokenResult]{Value: TokenResult{
				Tokens: tokens,
				Words:  len(strings.Fields(clip.Value)),
				Chars:  len(clip.Value),
				Sign:   sign,
			}}
			res <- value
			cache.Put(clip.Value, value)
		}
	}()
	return res
}
