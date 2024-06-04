package main

import (
	"encoding/json"
	"fmt"

	handler "github.com/sb-nour/providers-endpoints/api"
)

func main() {
	regions := handler.GetRegions()
	regionsJson, err := json.Marshal(regions)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(regionsJson))
}
