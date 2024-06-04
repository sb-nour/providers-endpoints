package main

import (
	"encoding/json"
	"fmt"

	"github.com/sb-nour/providers-endpoints/lib"
)

func main() {
	regions := lib.GetRegions()
	regionsJson, err := json.Marshal(regions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(regionsJson))
}
