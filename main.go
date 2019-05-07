package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

const (
	mapWidth  = 128
	mapHeight = 128
	mapSize   = mapHeight * (mapWidth + 1)
)

var mapCharacters = "ox"

func main() {
	http.HandleFunc("/random/map", GetMap)

	log.Printf("Listening at :8000, using %d out of %d CPUs\n", runtime.GOMAXPROCS(-1), runtime.NumCPU())

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func GetMap(w http.ResponseWriter, r *http.Request) {
	data := NewMap()
	fmt.Fprint(w, data)
}

func NewMap() string {
	var (
		mapData = make([]byte, 0, mapSize)
		src     = rand.NewSource(time.Now().UnixNano())
	)

	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			mapData = append(mapData, mapCharacters[src.Int63()%2])
		}
		mapData = append(mapData, '\n')
	}

	return string(mapData)
}
