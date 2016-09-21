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

type Evidences struct {
	XMLName     xml.Name   `xml:"evidences"`
	CompanyId   string     `xml:"companyId"`
	CompanyName string     `xml:"companyName"`
	Evidence    []Evidence `xml:"evidence"`
}

type Evidence struct {
	Type         string `xml:"evidenceType"`
	Name         string `xml:"evidenceName"`
	Path         string `xml:"evidencePath"`
	ImportStatus string `xml:"importStatus"`
	ClassName    string `xml:"className"`
	FormCode     string `xml:"formCode"`
}

type byPath []Evidence

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

func main() {
	var xmldata []byte
	var data Evidences

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}

	response, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		xmldata, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = xml.Unmarshal(xmldata, &data)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Company: %s has %d evidences\n", data.CompanyName, len(data.Evidence))
	sort.Sort(byPath(data.Evidence))

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, e := range data.Evidence {
		fmt.Fprintf(w, "%s\t%s\n", e.Path, e.Name)
	}
	w.Flush()
}
