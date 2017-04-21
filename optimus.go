package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	username := "OptimusPrime"
	conn := establishConnection()
	screen := make([]byte, 803) // using small tmo buffer for demonstrating
	secondLine := make([]byte, 39)
	bumper := make([]byte, 39)
	for {
		n, err := conn.Read(screen)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println("got", n, "bytes.")
		fmt.Println(string(screen))
		if screen[8] == 87 {
			//Name entry code
			_, err = conn.Write([]byte(username + "\n"))
		} else {
			//actual movement code
			copy(secondLine[:], screen[41:79])
			copy(bumper[:], screen[481:520])
			xPosition := findX(secondLine)
			if xPosition != -1 {

				fmt.Printf("Found x at: %v",xPosition)
				fmt.Println()
				carPosition := updateCarPosition(bumper)
				fmt.Printf("Car from %v to %v",carPosition,carPosition+14)
				fmt.Println()
				fmt.Println("----SECOND LINE----")
				fmt.Println(string(secondLine))
				fmt.Println("----SECOND LINE----")
				if xPosition>=carPosition && xPosition<= carPosition+14{
					vector:=createVector(xPosition,carPosition)
					fmt.Printf("Vector %v", vector)
					steer(vector, conn)
					time.Sleep(4 * time.Millisecond)
				}
			}
		}

	}
}

func establishConnection() net.Conn {
	conn, err := net.Dial("tcp", "protury.info:4242")
	if err != nil {
		fmt.Println("Connection not established")
		os.Exit(1)
	} else {
		fmt.Println("Connection established to protury.info")
	}
	return conn
}

func updateCarPosition(screen []byte) int {
	for i, b := range screen {
		if b == 47 { // this equals the forward slash that is at the top left corner of the car
			return i-3
		}
	}
	return -1
}

func findX(firstline []byte) int{
	for i, b := range firstline {
		if b == 88 { // this equals an "X"
			return i
		}
	}
	return -1
}

func createVector(xPosition int, carPosition int) int{
	if xPosition<=carPosition+6 {
		fmt.Println("Steer right")
		vector:=xPosition-carPosition
		if carPosition+13+vector < 40{
			return vector
		}

	} else {
		fmt.Println("Steer left")
		vector:=(xPosition)-(carPosition+15)
		if carPosition+vector > 1 {
			return vector
		}
	}
	if xPosition>=carPosition+6 {
		fmt.Println("Steer right")
		vector:=xPosition-carPosition

			return vector


	} else {
		fmt.Println("Steer left")
		vector:=(xPosition)-(carPosition+15)

			return vector

	}
	return 0
}

func steer(vector int, conn net.Conn){
	fmt.Printf("Steering with vector %v", vector)
	fmt.Println()
	if vector > 0 {
		for i:=vector; i > 0; i--{
			conn.Write([]byte{67})
		}
	} else {
		for i:=vector*-1; i > 0; i--{
			conn.Write([]byte{68})
		}
	}
}