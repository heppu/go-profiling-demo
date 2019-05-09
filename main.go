package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

const (
	mapWidth      = 128
	mapHeight     = 128
	mapSize       = mapHeight * (mapWidth + 1)
	mapCharacters = "ox"
)

var srcPool = &sync.Pool{
	New: func() interface{} {
		return rand.NewSource(time.Now().UnixNano())
	},
}

func main() {
	http.HandleFunc("/random/map", GetMap)

	log.Printf("Listening at :8000, using %d out of %d CPUs\n", runtime.GOMAXPROCS(-1), runtime.NumCPU())

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func GetMap(w http.ResponseWriter, r *http.Request) {
	data := NewMap()
	w.Write(data)
}

func NewMap() []byte {
	var (
		mapData   = make([]byte, 0, mapSize)
		src       = srcPool.Get().(rand.Source)
		r         = src.Int63()
		randIndex = 0
	)

	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			if randIndex == 63 {
				r = src.Int63()
				randIndex = 0
			}
			mapData = append(mapData, mapCharacters[r&1])
			randIndex++
			r >>= 1
		}
		mapData = append(mapData, '\n')
	}
	srcPool.Put(src)

	return mapData
}
