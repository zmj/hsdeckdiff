package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	formHTML, err := loadForm()
	if err != nil {
		fmt.Printf("Failed to load deck submission form: %v\n", err.Error())
		return
	}
	p, err := newParser()
	if err != nil {
		fmt.Printf("Failed to create parser: %v\n", err.Error())
		return
	}

	http.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			wr.Write(formHTML)
			return
		} else if req.Method != http.MethodPost {
			err = fmt.Errorf("Unsupported HTTP method %v ", req.Method)
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		err := req.ParseForm()
		if err != nil {
			err = fmt.Errorf("Failed to parse submitted form: %v", err.Error())
			http.Error(wr, err.Error(), http.StatusBadRequest)
			return
		}

		decksToParse := []deckToParse{
			{"maindeck", "main deck"},
			{"altdeck1", "alternate deck 1"},
			{"altdeck2", "alternate deck 2"},
		}

		formField := func(name, desc string) (string, error) {
			values, ok := req.PostForm[name]
			if !ok || len(values) == 0 {
				err := fmt.Errorf("Missing %v content", desc)
				http.Error(wr, err.Error(), http.StatusBadRequest)
				return "", err
			}
			return values[0], nil
		}

		parseDeck := func(raw, desc string) (*deck, error) {
			d, err := p.parse(strings.NewReader(raw))
			if err != nil {
				err := fmt.Errorf("Failed to parse deck %v: %v", desc, err.Error())
				http.Error(wr, err.Error(), http.StatusInternalServerError)
				return nil, err
			}
			return d, nil
		}

		var parsedDecks []*deck
		for _, d := range decksToParse {
			raw, err := formField(d.Field, d.Desc)
			if err != nil {
				return
			}
			deck, err := parseDeck(raw, d.Desc)
			if err != nil {
				return
			}
			parsedDecks = append(parsedDecks, deck)
		}
		maindeck, altdeck1, altdeck2 := parsedDecks[0], parsedDecks[1], parsedDecks[2]

		diff1 := diff(maindeck, altdeck1)
		diff2 := diff(maindeck, altdeck2)
		wr.Write([]byte(diff1.String()))
		wr.Write([]byte("\n"))
		wr.Write([]byte("\n"))
		wr.Write([]byte(diff2.String()))
		return
	})
	err = http.ListenAndServe(":6177", nil)
	if err != nil {
		fmt.Printf("Failed to start http server: %v", err.Error())
		return
	}

	fmt.Println("Exiting! bye")
}

type deckToParse struct {
	Field string
	Desc  string
}

func loadForm() ([]byte, error) {
	file, err := os.Open("form.html")
	if err != nil {
		return nil, fmt.Errorf("Failed to open form.html: %v", err.Error())
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to read form.html: %v", err.Error())
	}
	return content, nil
}
