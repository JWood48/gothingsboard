package main

import (
    //    "encoding/json"
"os"
"fmt"
"time"
"log"
"log/syslog"
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


    l3, err := syslog.New(syslog.LOG_ERR, "GoThingsBoard")
    //l3, err := syslog.Dial("udp", "localhost", syslog.LOG_ERR, "GoExample") // connection to a log daemon
    defer l3.Close()
    if err != nil {
        log.Fatal("error")
    }

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

   l3.Info("cfg: \n")
   l3.Info("host =>" +host+"<\n")
   l3.Info("port =>" +port+"<\n")
   l3.Info("token =>" +token+"<\n")
   l3.Info("\n")

   ticker := time.NewTicker(2 * time.Minute)
   quit := make(chan struct{})
   go func() {
    for {
     select {
     case <- ticker.C:
                // do stuff

        data := getData()
        
        l3.Print("data: "+string(data)

        url := "http://"+host+":"+port+"/api/v1/"+token+"/telemetry"
        b := new(bytes.Buffer)
        json.NewEncoder(b).Encode(data)
        l3.Info("url: "+url)	
	//log.Print("data: "+string(data))	
	//log.Print("data2:"+json.Marshal(data).value)	

        jd,err := json.Marshal(data)
        l3.Info("JSON: "string(jd))

        res, err := http.Post(url, "application/json; charset=utf-8", b)

        if err != nil {
           log.Print(err)
       }
       l3.Info("Request sent...")
       l3.Info(res)
       //io.Copy(os.Stdout, res.Body)


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
