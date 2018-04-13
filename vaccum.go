package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
	"strconv"
)

// Owner is the repository owner
type Owner struct {
	Login string
}

// Item is the single repository data structure
type Item struct {
	ID              int
	Name            string
	FullName        string `json:"full_name"`
	Owner           Owner
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

// JSONData contains the GitHub API response
type JSONData struct {
	Count int `json:"total_count"`
	Items []Item
}


func main() {
	st:= 1000000
	for i:=1;i<1000;i++ {
		time.Sleep(5 * time.Second)
		res, err := http.Get("https://api.github.com/search/repositories?q=stars:<"+strconv.Itoa(st)+"&sort=stars&order=desc&page=1&per_page=100")
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			log.Fatal("Unexpected status code", res.StatusCode)
		}
		data := JSONData{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}
		st = saveData(data)
	}	
}

func saveData(data JSONData)(int) {
	lastNbStars := 0
	log.Printf("Repositories found: %d", data.Count)
	const format = "%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Name", "Stars", "Created at", "Description")
	fmt.Fprintf(tw, format, "----------", "-----", "----------", "----------")
	//I will append the name of the project in this file :
	f, err := os.OpenFile("out.gi", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// And the description here
	f1, err := os.OpenFile("out.en", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range data.Items {
		desc := i.Description
		if len(desc) > 50 {
			desc = string(desc[:50]) + "..."
		}
		t, err := time.Parse(time.RFC3339, i.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(tw, format, i.Name, i.StargazersCount, t.Year(), desc)
		_, err = f.WriteString(i.Name+".\n")
		_, err = f1.WriteString(i.Description+".\n")
		if err != nil {
			log.Fatal(err)
		}
		lastNbStars = i.StargazersCount
	}
	defer f.Close()
	defer f1.Close()
	tw.Flush()
	return lastNbStars
}
