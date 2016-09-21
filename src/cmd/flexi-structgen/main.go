package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
)

type evidences struct {
	XMLName     xml.Name   `xml:"evidences"`
	CompanyId   string     `xml:"companyId"`
	CompanyName string     `xml:"companyName"`
	Evidence    []evidence `xml:"evidence"`
}

type evidence struct {
	Type         string `xml:"evidenceType"`
	Name         string `xml:"evidenceName"`
	Path         string `xml:"evidencePath"`
	ImportStatus string `xml:"importStatus"`
	ClassName    string `xml:"className"`
	FormCode     string `xml:"formCode"`
}

type properties struct {
	XMLName      xml.Name   `xml:"properties"`
	EvidenceName string     `xml:"evidenceName"`
	TagName      string     `xml:"tagName"`
	Property     []property `xml:"property"`
}

type property struct {
	Name string `xml:"name"`
}

// Pro trideni evidenci
type byPath []evidence

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func readUrl(url string) ([]byte, error) {
	var data []byte
	var err error

	response, err := http.Get(url)
	if err == nil {
		defer response.Body.Close()
		data, err = ioutil.ReadAll(response.Body)
	}

	return data, err
}

func main() {
	var xmldata []byte
	var data evidences

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s BASEURL (https://demo.flexibee.eu/c/demo/)\n", os.Args[0])
		os.Exit(1)
	}

	baseURL := os.Args[1]

	xmldata, err := readUrl(baseURL + "evidence-list.xml")
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(xmldata, &data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Company: %s has %d evidences\n", data.CompanyName, len(data.Evidence))
	sort.Sort(byPath(data.Evidence))

	log.Print("Fetching data")
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, e := range data.Evidence {
		propdata := properties{}
		xmldata, err := readUrl(baseURL + e.Path + "/properties.xml")
		if err == nil {
			err = xml.Unmarshal(xmldata, &propdata)
		}
		fmt.Fprintf(w, "%s\t%s (%d atributu)\n", e.Path, e.Name, len(propdata.Property))
		if err != nil {
			log.Printf("[%s] %s", e.Path, err)
		}
	}
	w.Flush()
}
