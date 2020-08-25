package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"io"
	"log"
	"strconv"
	"time"
	"flag"
	"sort"
)

func printMaps(tmax, tmin int, m1, m2 map[int][]int) {
	fmt.Println("Year", "Temp Above", tmax)
	for k, v := range m1 {
		t := 0
		for _, d := range v {
			t += d
		}
		sort.Ints(v)
		x := v[len(v)-1]
		fmt.Printf("%d, %d, %d %v\n", k, t, x, v)
	}
	fmt.Println("Year", "Temp Below", tmin)
	for k, v := range m2 {
		t := 0
		for _, d := range v {
			t += d
		}
		sort.Ints(v)
		x := v[len(v)-1]
		fmt.Printf("%d, %d, %d %v\n", k, t, x, v)
	}
}

func printMap(m1 map[int][]int) {
	for k, v := range m1 {
		t := 0
		for _, d := range v {
			t += d
		}
		sort.Ints(v)
		x := v[len(v)-1]
		fmt.Printf("%d, %d, %d %v\n", k, t, x, v)
	}
}

func main() {

	var heatSpell, coldSpell bool

	tmaxPtr := flag.Int("maxF", 90, "max Temp above")
	tminPtr := flag.Int("minF", 35, "min Temp below")

	flag.Parse()

	fArgs := flag.Args()

	if len(fArgs) == 0 {
		fmt.Println("Missing file")
		return
	}

	f, err := os.Open(fArgs[0])
	if err != nil {
		fmt.Println("No such file")
		return
	}

	defer f.Close()

	yrsHtSpells := make(map[int][]int)
	yrsCldSpells := make(map[int][]int)
	Tmax := int(*tmaxPtr)
	Tmin := int(*tminPtr)
	heatSpell = false
	coldSpell = false
	heatSpellLen := 0
	coldSpellLen := 0

	r := csv.NewReader(f)
	r.LazyQuotes = true
	tyear := 0
	for {
		line, err := r.Read()
		if err == io.EOF {
			if coldSpell == true {
				yrsCldSpells[tyear] = append(yrsCldSpells[tyear],
									coldSpellLen)
			}
			if heatSpell == true {
				yrsHtSpells[tyear] = append(yrsHtSpells[tyear],
									heatSpellLen)
			}
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		//fmt.Println(line[2], " ", line[4], " ", line[5])
		date := line[2]
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue;
		}
		tmax, err := strconv.ParseInt(line[4], 10, 32)
		if err != nil {
			tmax = 0
		}
		tmin, err := strconv.ParseInt(line[5], 10, 32)
		if err != nil {
			tmin = 0
		}
		
		tyear = t.Year()

		if tmax > 0 && int(tmax) >= Tmax {
			heatSpell = true
			heatSpellLen += 1
		} else {
			if heatSpell == true {
				if heatSpellLen > 3 {
					if len(fArgs) > 1 && fArgs[1] == "max" {
					//	fmt.Printf("%d Days above %dF %s\n",
					//		heatSpellLen, Tmax, date)
					}
				}
				heatSpell = false
				yrsHtSpells[t.Year()] = append(yrsHtSpells[t.Year()],
									heatSpellLen)
				heatSpellLen = 0
			}
		}
		if tmin > 0 && int(tmin) <= Tmin {
			coldSpell = true
			coldSpellLen += 1
		} else {
			if coldSpell == true {
				if coldSpellLen > 3 {
					if len(fArgs) > 1 && fArgs[1] == "min" {
					//	fmt.Printf("%d Days below %dF %s\n",
					//		coldSpellLen, Tmin, date)
					}
				}
				coldSpell = false
				yrsCldSpells[t.Year()] = append(yrsCldSpells[t.Year()],
								coldSpellLen)
				coldSpellLen = 0
			}
		}
	}
	if len(fArgs) > 1 {
		if fArgs[1] == "min" {
			printMap(yrsCldSpells)
		} else if fArgs[1] == "max" {
			printMap(yrsHtSpells)
		} else {
			printMaps(int(*tmaxPtr), int(*tminPtr), yrsHtSpells, yrsCldSpells)
		}
	} else {
		printMaps(int(*tmaxPtr), int(*tminPtr), yrsHtSpells, yrsCldSpells)
	}
}
