package main


import ("gfx2"
	"fmt"
	"time")
	
func main () {
	fmt.Println("start")
	
	gfx2.Fenster(400,400)
	
	channel := make(chan int)
	go rk(channel)
	
	for {
		select{
		case data := <- channel:
			fmt.Println(data)
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("--")
		}
	}}
	
func rk (channel chan int) {
	for {
		button, status, _, _ := gfx2.MausLesen1()
		if button == 1 && status == 1 {
			channel <- 1
		}
	}
}
