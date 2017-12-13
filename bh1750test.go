package main

import "fmt"
import "github.com/explicite/i2c/bh1750"
//import "testing"
import "time"


func main() {
	fmt.Printf("hello, world\n")
	
	//var dev bh1750.BH1750

	dev := new(bh1750.BH1750)

	dev.Init(0x23, 0x01)

	start := time.Now()
	
	//resolution := bh1750.ConHRes1lx 
	//resolution := bh1750.ConHRes05lx
	resolution := bh1750.ConLRes4lx

	n := 100

	var sum float32 

	for i := 0; i < n; i++ {
	//for {
		
		
		res,err := dev.Illuminance(byte(resolution))
		fmt.Printf("TEST: %v - err=%v\n", res, err)
		if err != nil {
			fmt.Printf("TEST: err=%v\n", err)
		}
		sum += res
	}

	elapsed := time.Since(start)
	elapsedPerCall := elapsed / time.Duration(n)
	mean := sum / float32(n)

	fmt.Printf("elapsed=%v, prRequest=%v\n", elapsed, elapsedPerCall)

	fmt.Printf("Sum=%v, mean=%v\n", sum, mean)


	//var dev = bh1750.BH1750()
	//var dev = BH1750.Init(0x23, 0x01)



}


//func ReadIllumination(dev *bh1750.BH1750) {

	
	
//	dev.Illuminance(bh1750.ConHRes1lx)
	
//}



