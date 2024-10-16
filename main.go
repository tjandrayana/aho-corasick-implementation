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
	NameIdentifier    string = "name"
	AddressIdentifier string = "address"
)

type Module struct {
	AddressDictionary   map[string][]string
	NameDictionary      map[string][]string
	AddressMapAhorasick map[string]ahocorasick.AhoCorasick
	NameMapAhorasick    map[string]ahocorasick.AhoCorasick
}

func Init() Module {

	nameDic := CreateDictionary("name.csv")
	addDic := CreateDictionary("add.csv")

	addBuilder := ahocorasick.NewAhoCorasickBuilder(ahocorasick.Opts{
		AsciiCaseInsensitive: true,
		MatchOnlyWholeWords:  true,
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

	for k, v := range dictionary {
		fmt.Printf("L3_geo_id: %s, Word: %v\n", k, v)
	}

	fmt.Printf("\n")

	return dictionary
}

func BuildAhorasickMap(builder ahocorasick.AhoCorasickBuilder, dictionary map[string][]string) map[string]ahocorasick.AhoCorasick {
	start := time.Now()
	mapAhorasick := make(map[string]ahocorasick.AhoCorasick)

	for k, v := range dictionary {
		mapAhorasick[k] = builder.Build(v)
	}

	fmt.Printf("BuildAhorasickMap  taken %v micro second\n ", time.Since(start).Microseconds())

	return mapAhorasick
}

type TestCase struct {
	L3_geo_id   string
	AddressName string
	Name        string
	Identifier  string
}

func main() {

	st := time.Now()
	m := Init()

	tc := []TestCase{
		{
			L3_geo_id:   "14_184_2389",
			AddressName: "Perum Negeri Jaya Pedasong Pekajangan Kendaldoyong--###--###--###--lamparCilik,rumah BILLY wea 081246944944",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbahnari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "14_188_2466",
			AddressName: "NANGGUNGAN RT 5 RW 3 rumah mbah nari tukang--###--###--###--",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "15_222_3043",
			AddressName: "jln sahrun dsn nanggungan rt 02 rw 08--###--###--###--( tembok kramik merah)",
			Identifier:  AddressIdentifier,
		},
		{
			L3_geo_id:   "15_222_3043",
			AddressName: "jln Kasman. dsn nanggungan RT/RW 01/09 ds Jatirejo--###--###--###--rumah belakang masjid nanggungan",
			Identifier:  AddressIdentifier,
		},

		// Name Testcase
		{
			L3_geo_id:   "15_228_3194",
			AddressName: "SaPNCd Evander",
			Identifier:  NameIdentifier,
		},
		{
			L3_geo_id:   "11_146_1636",
			AddressName: "Bang Wizi Satria",
			Identifier:  NameIdentifier,
		},
		{
			L3_geo_id:   "11_146_1636",
			AddressName: "BangWizi Satria",
			Identifier:  NameIdentifier,
		},
	}

	m.RunTestCase(tc)

	dur := time.Since(st)
	fmt.Printf("Main Package total duration with %d Testcases took %v micro second (%vms)\n", len(tc), dur.Microseconds(), dur.Milliseconds())
}

func (m Module) RunTestCase(tc []TestCase) {

	for i, t := range tc {
		start := time.Now()

		text := t.AddressName
		l3 := t.L3_geo_id
		fmt.Printf("Test case %d: %s \n", i+1, text)

		switch t.Identifier {
		case NameIdentifier:
			fmt.Printf("dictionary[%s] : %v \n", l3, m.NameDictionary[l3])
			if len(m.NameDictionary[l3]) > 0 {
				matches := m.NameMapAhorasick[l3].FindAll(text)

				fmt.Printf("code words found : ")

				fmt.Printf("[ ")
				for _, match := range matches {
					fmt.Printf("%s ", text[match.Start():match.End()])
				}
				fmt.Printf(" ]\n")
			}
		case AddressIdentifier:
			fmt.Printf("dictionary[%s] : %v \n", l3, m.AddressDictionary[l3])
			if len(m.AddressDictionary[l3]) > 0 {
				matches := m.AddressMapAhorasick[l3].FindAll(text)

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
