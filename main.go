package	main

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type CsvLine struct {
	domains string
}

func readCsvFile() [][]string {
	f, err := os.Open("domain.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}


func gethttpstatuscode() (*http.Response, int64) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	starttime := time.Now().UTC().UnixNano()
	resp, err := http.Head("https://192.168.1.100:444")
	latency := time.Now().UTC().UnixNano() - starttime
	if err != nil {
		log.Fatal(err)
	}
	return resp, latency
}

func gettcpstatus() net.Conn {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("192.168.1.100", "23"), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort("192.168.1.100", "23"))
	}
		return conn
}

func main()  {
	resp, latency := gethttpstatuscode()
	csvLines := readCsvFile()
	for _, line := range csvLines{
		data := CsvLine{
			domains: line[0],
		}
		println(data.domains)
	}
	tcpstatus := gettcpstatus()
	fmt.Println(tcpstatus)
	fmt.Println(resp.StatusCode, latency, csvLines[1])
}