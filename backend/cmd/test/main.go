package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("test")

	// var gracefulStop = make(chan os.Signal)
	// signal.Notify(gracefulStop, syscall.SIGTERM)
	// signal.Notify(gracefulStop, syscall.SIGINT)
	// go func() {
	// 	sig := <-gracefulStop
	// 	fmt.Printf("caught sig: %+v", sig)
	// 	fmt.Println("Wait for 5 second to finish processing")
	// 	time.Sleep(5 * time.Second)
	// 	os.Exit(0)
	// }()

	go test()

	time.Sleep(80 * time.Second)
}

func test() {
	fmt.Println("test 2")
	time.Sleep(60 * time.Second)
	fmt.Println("test 3")
}
