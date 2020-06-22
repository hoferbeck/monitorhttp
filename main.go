package	main

import (
	"crypto/tls"
	"encoding/csv"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)
type CsvLine struct {
	domains string
	protocol string
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


func gethttpreq(url string) (int, int64, error) {
	starttime := time.Now().UTC().UnixNano()
	resp, err := http.Get(url)
	if err != nil {
		//log.Print(err)
		return 0, 0, err
	}
	latency := time.Now().UTC().UnixNano() - starttime
	return resp.StatusCode, latency, nil
}

func gettcpreq (host string, port string) (net.Conn, int64, error) {
	timeout := time.Second
	starttime := time.Now().UTC().UnixNano()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	latency := time.Now().UTC().UnixNano() - starttime
	return conn, latency, err
}

func main()  {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	csvLines := readCsvFile()
	for{
		for _, line := range csvLines{
			data := CsvLine{
				domains: line[0],
				protocol: line[1],
			}
			println(data.domains)
			spliteddata := strings.Split(data.domains, ":")

			if data.protocol == "tcp"{
				tcpreq, tcplatency, err := gettcpreq(spliteddata[0], spliteddata[1])
				println(tcpreq,tcplatency,err)

			} else if data.protocol == "http" {
				httpreq, httplatency, err := gethttpreq(data.domains)
				println(httpreq, httplatency, err)

			}

		}
		time.Sleep(1 * time.Second)
	}
}