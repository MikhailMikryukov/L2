package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	return recursiveMerge(channels...)[0]
}

// Рекурсивно объединяем каналы пока всё не объединится в один
func recursiveMerge(channels ...<-chan interface{}) []<-chan interface{} {
	if len(channels) == 1 {
		return channels
	}

	var mergedChannels []<-chan interface{}

	if len(channels)%2 == 0 {
		mergedChannels = make([]<-chan interface{}, 0, len(channels)/2)
	} else {
		// Для нечетного количества каналов.
		// Сразу добавляем последний элемент, так как далее объединение попарно
		mergedChannels = make([]<-chan interface{}, 0, len(channels)/2+1)
		mergedChannels = append(mergedChannels, channels[len(channels)-1])
	}
	for i := 0; i < len(channels)-1; i += 2 {
		c := merge(channels[i], channels[i+1])
		mergedChannels = append(mergedChannels, c)
	}

	return recursiveMerge(mergedChannels...)
}

// Функция объединения двух каналов из L2.7
func merge(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil || b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

// L2.14 Функция or (объединение done-каналов)
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
