package main

import (
    //    "encoding/json"
"os"
"fmt"
"time"
"log"
"log/syslog"
//"io"
"bytes"
"encoding/json"
"net/http"
"github.com/creamdog/gonfig"

    //"golang.org/x/exp/io/i2c"
    //"github.com/biturbo/bme280"

	"golang.org/x/exp/io/i2c"

	"github.com/quhar/bme280"


//"github.com/davecheney/i2c"
//"github.com/quinte17/bme280"


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
   i2cdev, _ := config.GetString("i2c", "n/a")
   interval, _ := config.GetInt("interval", 120)

   l3.Info("cfg: \n")
   l3.Info("host =>" +host+"<\n")
   l3.Info("port =>" +port+"<\n")
   l3.Info("token =>" +token+"<\n")
   l3.Info("i2c_dev =>" +i2cdev+"<\n")
   l3.Info(""+fmt.Sprintf("interval =>%v<\n",interval))
   l3.Info("\n")

   tickerInterval := time.Duration(interval) * time.Second

   l3.Info("tickerInterval: "+fmt.Sprintf("%v", tickerInterval)+"\n")
   fmt.Printf("tickerInterval: %v", tickerInterval)


   ticker := time.NewTicker(tickerInterval)
   quit := make(chan struct{})
   go func() {
    for {
     select {
     case <- ticker.C:
                // do stuff

        data := getData(i2cdev)
        //data := getData2()
        
        l3.Info("data: "+fmt.Sprintf("%v", data))

        url := "http://"+host+":"+port+"/api/v1/"+token+"/telemetry"
        b := new(bytes.Buffer)
        json.NewEncoder(b).Encode(data)
        l3.Info("url: "+url)	
	//log.Print("data: "+string(data))	
	//log.Print("data2:"+json.Marshal(data).value)	

        //jd,err := json.Marshal(data)
        //l3.Info("JSON: "+fmt.Sprintf("%v", jd))

        res, err := http.Post(url, "application/json; charset=utf-8", b)

        if err != nil {
           log.Print(err)
       }
       l3.Info("Request sent...")
       l3.Info(fmt.Sprintf("Response: %v", res))
       //io.Copy(os.Stdout, res.Body)


   case <- quit:
     ticker.Stop()
     return
 }
}
}()

select {}

}


func getData(i2cdev string) Message {

	//d, err := i2c.Open(&i2c.Devfs{Dev: i2cdev}, bme280.I2CAddr)
	d, err := i2c.Open(&i2c.Devfs{Dev: i2cdev}, 0x76)
	if err != nil {
		fmt.Printf("Could not open device "+i2cdev)
		panic(err)
	}

	b := bme280.New(d)
	err = b.Init()

	if err != nil {
		fmt.Printf("Could not init bme280 from "+i2cdev)
		panic(err)
	}
	
	t, p, h, err := b.EnvData()

	if err != nil {
		fmt.Printf("Could not read bme280 data from "+i2cdev)
		panic(err)
	}

	fmt.Printf("Temp: %fC, Press: %fhPa, Hum: %f%%\n", t, p, h)

	data := Message{t, h, p, true}

	return data
}

