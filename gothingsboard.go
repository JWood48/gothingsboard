package main

import (
    //    "encoding/json"
    "os"
    "fmt"
    "time"
    "log"
    "io"
    "bytes"
    "encoding/json"
    "net/http"
    "github.com/creamdog/gonfig"

    //"golang.org/x/exp/io/i2c"
    //"github.com/biturbo/bme280"

    "github.com/davecheney/i2c"
    "github.com/quinte17/bme280"
//    "github.com/d2r2/go-i2c"
)


type Message struct {
    Temperature float64 `json:"temperature"`
    Humidity float64 `json:"humidity"`
    Pressure float64 `json:"pressure"`
    Active bool `json:"active"`
}


func main() {


    confFile := os.Args[1]

    fmt.Printf("CFG_FILE: "+confFile+"\n")

    f, err := os.Open(confFile)
    if err != nil {
	log.Print(err)
        // TODO: error handling
    }
    defer f.Close();
    config, err := gonfig.FromJson(f)
    if err != nil {
	log.Print(err)
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

    ticker := time.NewTicker(2 * time.Minute)
    quit := make(chan struct{})
    go func() {
        for {
           select {
           case <- ticker.C:
                // do stuff
	
    data := getData()
log.Print(data)
	url := "http://"+host+":"+port+"/api/v1/"+token+"/telemetry"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	log.Print("url: "+url)	
	//log.Print("data: "+string(data))	
	//log.Print("data2:"+json.Marshal(data).value)	
	
	jd,err := json.Marshal(data)
	log.Print(jd)

	res, err := http.Post(url, "application/json; charset=utf-8", b)
	
    if err != nil {
	log.Print(err)
    }
	log.Print("Request sent...")
	log.Print(res)
	io.Copy(os.Stdout, res.Body)


           case <- quit:
           ticker.Stop()
           return
        }
    }
 }()

select {}

}

func getData() Message {

	dev, err := i2c.New(0x76, 1)
	if err != nil {
		log.Print(err)
	}
	bme, err := bme280.NewI2CDriver(dev)
	if err != nil {
		log.Print(err)
	}
	
	readData, err := bme.Readenv()
	if err != nil {
		log.Print(err)
	}
	

	log.Print("DATA READ1: ")
	log.Print(readData)

	data := Message{readData.Temp, readData.Hum, readData.Press, true}

	log.Print("DATA READ2: ")
	log.Print(data)

	return data

}
