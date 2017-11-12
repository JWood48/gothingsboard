package main

import (
    //    "encoding/json"
    "os"
    "fmt"
    "github.com/creamdog/gonfig"
    //    "github.com/d2r2/go-i2c"
)

type Configuration struct {
    host string
    port string
    token string
}


func main() {

    fmt.Printf("HOME: "+os.Getenv("HOME")+"\n")

    userHome := os.Getenv("HOME")
    confFile := userHome+"/gotb.cfg"

    fmt.Printf("CFG_FILE: "+confFile+"\n")

    f, err := os.Open(confFile)
    if err != nil {
        // TODO: error handling
    }
    defer f.Close();
    config, err := gonfig.FromJson(f)
    if err != nil {
        // TODO: error handling
    }


    host, _ := config.GetString("host", "n/a")
    port, _ := config.GetString("port", "n/a")
    token, _ := config.GetString("token", "n/a")

    fmt.Printf("cfg: \n")
    fmt.Printf("host =>" +host+"<\n")
    fmt.Printf("port =>" +port+"<\n")
    fmt.Printf("token =>" +token+"<\n")
    fmt.Printf("\n")


    fmt.Printf("hello, world\n")
}

func getData() {

    // Create new connection to I2C bus on 2 line with address 0x27
    i2c, err := i2c.NewI2C(0x27, 2)
    if err != nil { log.Fatal(err) }
    // Free I2C connection on exit
    defer i2c.Close()
    
    
    // Here goes code specific for sending and reading data
    // to and from device connected via I2C bus, like:
    _, err := i2c.Write([]byte{0x1, 0xF3})
    if err != nil { log.Fatal(err) }

}
