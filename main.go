package main

import (
	"encoding/csv"
	"fmt"
	"strings"

	mcsv "github.com/Nidal-Bakir/first_go/pkg/csv"
)

func main() {
	data :=
		`name,age,has_pet
Jon,"100",true
"Fred ""The Hammer"" Smith",42,false
Martha,37,"true"`

	csvReader := csv.NewReader(strings.NewReader(data))
	csvArrData, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(csvArrData)

	type Person struct {
		HasPet bool   `csv:"has_pet"`
		Name   string `csv:"name"`
		Age    int    `csv:"age"`
	}

	personArr := new([]Person)
	err = mcsv.Unmarshal(csvArrData, personArr)
	if err != nil {
		panic(err)
	}

	fmt.Println(*personArr)
	fmt.Println((*personArr)[1].Name)
	fmt.Println("=====================================")
	fmt.Println("=====================================")

	marsheldData, err := mcsv.Marshal(personArr)
	if err != nil {
		panic(err)
	}
	fmt.Println(marsheldData)

}
