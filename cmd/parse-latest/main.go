// parses most recent test-data file and generates yaml file

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	nameday "github.com/skrawler/crawler-nameday"
)

func mostRecentDataFile() string {
	all, err := filepath.Glob("test-data/nameday-swe-dagensnamn-*")
	if err != nil {
		log.Fatal(err)
	}
	outName := all[0]
	highest, err := dateFromFilenameSuffix(outName)
	if err != nil {
		log.Fatal(err)
	}
	// find most recent date
	for _, file := range all {
		t, err := dateFromFilenameSuffix(file)
		if err == nil {
			if t.After(highest) {
				highest = t
				outName = file
			}
		}
	}
	return outName
}

func dateFromFilenameSuffix(file string) (time.Time, error) {
	name := filepath.Base(file)
	if len(name) < 24 {
		return time.Now(), fmt.Errorf("too short name")
	}
	name = name[23:] // YYYY-MM-DD
	return time.Parse("2006-01-02", name)
}

func main() {
	inFile := mostRecentDataFile()
	log.Println("Parsing", inFile)
	s, err := readTextFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	info, err := nameday.ExtractNamedaysDagensnamn(s)
	if err != nil {
		log.Fatal(err)
	}

	outFile := "data/nameday_swe.yml"
	err = saveAsYaml(outFile, info)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated parsed data file", outFile)
}

func readTextFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveAsYaml(filename string, list []nameday.Nameday) error {
	y, err := yaml.Marshal(list)
	if err != nil {
		spew.Dump(list)
		return fmt.Errorf("marshal error: %v", err)
	}
	s := string(y)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(s)
	return err
}
