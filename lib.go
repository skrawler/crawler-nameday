package nameday

import (
	"fmt"
	"log"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kennygrant/sanitize"
	natural "github.com/skrawler/go-natural"
)

// ...
var (
	NamedaySwe = []Nameday{}
)

// Nameday ...
type Nameday struct {
	Name     string
	Date     natural.MonthDay
	Official bool
}

func init() {
	log.SetFlags(log.Lshortfile)

	data, err := DataFile("data/nameday_swe.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &NamedaySwe)
	if err != nil {
		panic(err)
	}
}

// DataFile returns data from a go-bindata embedded yaml file from the "data/" folder
func DataFile(filename string) ([]byte, error) {
	return Asset(filename)
}

// SweNamesOnDate returns name days matching given query
func SweNamesOnDate(month, day int) []Nameday {
	matches := []Nameday{}
	query := fmt.Sprintf("%02d-%02d", month, day)
	for _, nd := range NamedaySwe {
		if nd.Date.String() == query {
			matches = append(matches, nd)
		}
	}
	return matches
}

// SweNamedayFor returns the Nameday info for given name
func SweNamedayFor(name string) (Nameday, error) {
	name = ucfirst(name)
	for _, nd := range NamedaySwe {
		if nd.Name == name {
			return nd, nil
		}
	}
	return Nameday{}, fmt.Errorf("no match")
}

// ucfirst returns s with first letter in upper case
func ucfirst(s string) string {
	first := string([]rune(s)[0:1])
	rest := ""
	if len(s) > 1 {
		rest = string([]rune(s)[1:])
	}
	return strings.ToUpper(first) + strings.ToLower(rest)
}

// StrBetween returns the string between starter and ender markers
func StrBetween(s, starter, ender string) (string, error) {
	startI := strings.Index(s, starter)
	if startI == -1 {
		return "", fmt.Errorf("starter not found")
	}
	startI += len(starter)
	endI := strings.Index(s, ender)
	if endI == -1 {
		return "", fmt.Errorf("ender not found")
	}
	return s[startI:endI], nil
}

// ExtractNamedaysDagensnamn ...
func ExtractNamedaysDagensnamn(s string) ([]Nameday, error) {
	matches := []Nameday{}

	s = strings.Replace(s, "\n", "___", -1)
	s = strings.Replace(s, `<br />`, "\n", -1)

	s, err := StrBetween(s,
		`<div class="span3 name-list no-margin">`,
		`<div class="clear"></div>`)
	if err != nil {
		return nil, err
	}

	rows := strings.Split(s, "\n")

	for _, row := range rows {
		current := Nameday{}

		row = sanitize.HTML(row)
		row = strings.TrimSpace(row)
		if row == "" {
			continue
		}

		tmpParts := strings.Split(row, "___")
		if len(tmpParts) < 7 {
			fmt.Println("Skipping", tmpParts)
			continue
		}

		parts := []string{}
		for idx := range tmpParts {
			part := strings.TrimSpace(tmpParts[idx])
			if part == "" {
				continue
			}
			parts = append(parts, part)
		}

		name := parts[0]
		if len(name) <= 0 {
			fmt.Println("Skipping", parts)
			for idx := range parts {
				fmt.Println("part", idx, ":", parts[idx])
			}

			continue
		}

		if name[len(name)-1] == '*' {
			name = name[0 : len(name)-1]
			name = strings.TrimSpace(name)
		}

		current.Name = name
		date, err := parseDateToMMDD(parts[1])
		if err != nil {
			fmt.Println("date parse error:", err, "on line '", row, "'")
			continue
		}
		current.Date = date
		current.Official = true
		if len(parts) > 2 && parts[2] == "(inofficiell)" {
			current.Official = false
		}
		//fmt.Println("Found", current)
		matches = append(matches, current)
	}
	return matches, nil
}

// returns a "MM-DD" string from input such as "15 december"
func parseDateToMMDD(s string) (natural.MonthDay, error) {
	return natural.ParseDateIntoMonthDay(s)
}
