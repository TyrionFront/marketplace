package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"common"
	"consequent"
	"parallel"
)

func errCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fStart := time.Now()
	content, err := os.Open("../Archive/data.json")
	errCheck(err)
	defer content.Close()

	var ds common.PointsSet
	err2 := json.NewDecoder(content).Decode(&ds)
	errCheck(err2)

	log.Printf("Reading + decoding took: %v", time.Since(fStart))
	fmt.Printf("ds size: %v\n\n", len(ds))

	startC := time.Now()
	consequentRes := consequent.Consequent(ds)

	log.Printf("Consequent calculations took: %v\n\n", time.Since(startC))

	consequentOut, err3 := json.Marshal(consequentRes)
	errCheck(err3)
	os.WriteFile("./output/consequentOut.json", consequentOut, 0644)

	startP := time.Now()
	parallelRes := parallel.Parallel(ds)
	log.Printf("Parallel calculations took: %v\n\n", time.Since(startP))

	parallelOut, err3 := json.Marshal(parallelRes)
	errCheck(err3)
	os.WriteFile("./output/parallelOut.json", parallelOut, 0644)
}
