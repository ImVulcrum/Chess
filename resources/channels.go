package main

import (
	"fmt"
	"gfxw"
	"time"
)

func main() {
	fmt.Println("start")

	gfxw.Fenster(400, 400)

	channel := make(chan int)
	go rk(channel)

	for {
		select {
		case data := <-channel:
			fmt.Println(data)

		default:
			fmt.Println("----")

		}
		time.Sleep(50 * time.Millisecond)
		fmt.Println("hellooooo")
	}
}

func rk(channel chan int) {
	for {
		button, status, _, _ := gfxw.MausLesen1()
		if button == 1 && status == 1 {
			channel <- 1
		}
	}
}
