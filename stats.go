package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func initStats() map[string]int {
	statsListMapBuffer := make(map[string]int)
	//statsListMap["target hit"] = 0
	go statsLog()

	return statsListMapBuffer
}

func hitTarget(ipAddress string) {
	value, exist := statsListMap[ipAddress]
	if exist == false {
		statsListMap[ipAddress] = 1
		return
	}
	statsListMap[ipAddress] = value + 1
	return
}

func statsLog() {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	//f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\n")
	for one := range statsListMap {
		f.WriteString(one + "\t " + strconv.Itoa(statsListMap[one]) + "\n")
	}

	fmt.Println("Write current stats: ", statsListMap)

	time.Sleep(time.Minute * 1)
	statsLog()
}
