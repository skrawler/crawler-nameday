// downloads and writes to disk list of name days in test-data folder

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://www.dagensnamnsdag.nu/namnsdagar/"
	data, err := httpGetDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	fileName := "test-data/nameday-swe-dagensnamn-" + time.Now().Format("2006-01-02")

	log.Println("Updating cache for", url, "to", fileName)
	err = writeBinaryFile(fileName, data)
	if err != nil {
		log.Fatal(err)
	}
}

func writeBinaryFile(fileName string, data []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func httpGetDocument(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
