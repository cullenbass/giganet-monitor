package main

import (
	"io"
	"net/http"
	"net"
	"encoding/csv"
	"os"
	"strconv"
	"github.com/tcnksm/go-httpstat"
	"time"
	"io/ioutil"
	"fmt"
	"flag"
)

type Data struct {
	pingTime int
	validData bool
}

func write(ch chan Data) {
	f, err := os.OpenFile("data.csv", os.O_APPEND, 0777)
	created := false
	if err != nil {
		f, err = os.Create("data.csv")
		created = true
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if created {
		err := w.Write([]string{"datestamp","ping", "valid"})
		if err != nil {
			panic(err)
		}
	}
	for data := range ch {
		err:= w.Write([]string{time.Now().Format(time.RFC3339), strconv.Itoa(data.pingTime), strconv.FormatBool(data.validData)})
		if err!= nil {
			fmt.Println(err)
			continue
		}
		w.Flush()
	}
}

func verify(url string, ch chan Data) {
	req, _ := http.NewRequest("GET", url, nil)
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 60 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
	}
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)
	client := &http.Client{Timeout:60 * time.Second, Transport: netTransport}
	res, err := client.Do(req)
	if err != nil {
		ch <- Data{-1, false}
		return
	}
	_, err = io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil{
		ch <- Data{-1, false}
	} else {
		var valid bool
		if res.StatusCode >= 300 {
			valid = false
		} else {
			valid = true
		}
		ping := int(result.Connect / time.Millisecond)
		ch <- Data{ping, valid}
	}
	netTransport.CloseIdleConnections()
}

func main() {
	var comChan = make(chan Data)
	url := "http://www.google.com"
	interval := 1
	flag.Parse()
	if flag.NArg() > 2 || flag.NArg() == 1 {
		fmt.Println("Invalid arguments. Example: gigabet <full http url to test> <number of minutes between tests>")
		os.Exit(1)
	} else if flag.NArg()) == 2 {
		url = flag.Arg(0)
		i,err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fmt.Println("Invalid arguments. Example: gigabet <full http url to test> <number of minutes between tests>")
			os.Exit(1)
		} else {
			interval = i
		}

	}
	defer close(comChan)
	go write(comChan)
	for {
		go verify(url, comChan)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

