package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	ahocorasick "github.com/petar-dambovaliev/aho-corasick"
)

const (
	NameIdentifier    string = "NameIdentifier"
	AddressIdentifier string = "AddressIdentifier"
)

type Module struct {
	AddressDictionary   map[string][]string
	NameDictionary      map[string][]string
	AddressMapAhorasick map[string]ahocorasick.AhoCorasick
	NameMapAhorasick    map[string]ahocorasick.AhoCorasick
}

func Init() Module {

	nameDic := CreateDictionary("files/name.csv")
	addDic := CreateDictionary("files/add.csv")

	addBuilder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  false,
		MatchKind:            ahocorasick.LeftMostFirstMatch,
		DFA:                  true,
	})
	nameBuilder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  false,
		MatchKind:            ahocorasick.LeftMostFirstMatch,
		DFA:                  true,
	})

	mAddAho := BuildAhorasickMap(addBuilder, addDic)
	mNameAho := BuildAhorasickMap(nameBuilder, nameDic)

	return Module{
		NameDictionary:      nameDic,
		NameMapAhorasick:    mNameAho,
		AddressDictionary:   addDic,
		AddressMapAhorasick: mAddAho,
	}
}

// Function to generate dictionary
func CreateDictionary(filename string) map[string][]string {
	dictionary := make(map[string][]string)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return dictionary
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return dictionary
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if i == 0 {
			continue
		}
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			continue
		}

		l3GeoID := parts[0]
		raw := strings.TrimSpace(parts[1])

		raw = strings.Trim(raw, "[]")

		var word []string
		err := json.Unmarshal([]byte(raw), &word)
		if err != nil {
			word = strings.Split(raw, ",")
			for i := range word {
				word[i] = strings.Trim(word[i], `"`)
			}
		}

		dictionary[l3GeoID] = word
	}

	fmt.Printf("Total Dictionary Loaded : %v \n", len(dictionary))

	return dictionary
}

func BuildAhorasickMap(builder ahocorasick.AhoCorasickBuilder, dictionary map[string][]string) map[string]ahocorasick.AhoCorasick {
	start := time.Now()
	mapAhorasick := make(map[string]ahocorasick.AhoCorasick)

	for k, v := range dictionary {
		mapAhorasick[k] = builder.Build(v)
	}

	fmt.Printf("BuildAhorasickMap  taken %v micro second\n", time.Since(start).Microseconds())

	return mapAhorasick
}

type TestCase struct {
	KeyID      string `json:"key_id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

func main() {

	st := time.Now()
	m := Init()

	// Unmarshal the JSON file
	tc, err := UnmarshalJSONFile("files/testcases.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	m.RunTestCase(tc)

	dur := time.Since(st)
	fmt.Printf("Main Package total duration with %d Testcases took %v micro second (%vms)\n", len(tc), dur.Microseconds(), dur.Milliseconds())
}

func (m Module) RunTestCase(tc []TestCase) {

	for i, t := range tc {
		start := time.Now()

		text := t.Name
		key := t.KeyID
		fmt.Printf("Test case %d: %s \n", i+1, text)

		switch t.Identifier {
		case NameIdentifier:
			fmt.Printf("dictionary[%s] : %v \n", key, m.NameDictionary[key])
			if len(m.NameDictionary[key]) > 0 {
				matches := m.NameMapAhorasick[key].FindAll(text)

				fmt.Printf("code words found : ")

				fmt.Printf("[ ")
				for _, match := range matches {
					fmt.Printf("%s ", text[match.Start():match.End()])
				}
				fmt.Printf(" ]\n")
			}
		case AddressIdentifier:
			fmt.Printf("dictionary[%s] : %v \n", key, m.AddressDictionary[key])
			if len(m.AddressDictionary[key]) > 0 {
				matches := m.AddressMapAhorasick[key].FindAll(text)

				fmt.Printf("code words found : ")

				fmt.Printf("[ ")
				for _, match := range matches {
					fmt.Printf("%s ", text[match.Start():match.End()])
				}
				fmt.Printf(" ]\n")
			}

		}

		fmt.Printf("Process FindAll  taken %v micro second\n ", time.Since(start).Microseconds())
		fmt.Printf("============= \n\n")
	}

}

// Function to unmarshal JSON file to struct
func UnmarshalJSONFile(filename string) ([]TestCase, error) {
	// Read the file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the test cases
	var testCases []TestCase

	// Unmarshal the JSON data into the slice of structs
	err = json.Unmarshal(file, &testCases)
	if err != nil {
		return nil, err
	}

	return testCases, nil
}
