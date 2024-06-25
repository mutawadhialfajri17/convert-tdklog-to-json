package utils

import (
	"fmt"
	"log"
)

func PrintLog(in interface{}) {
	if e, ok := in.(error); ok {
		log.Fatalln("ERROR :", e)
	} else {
		log.Println("INFO :", in)
	}
}

func init() {
	fmt.Println("======= Convert Request From Metadata (TDK Log) =======")
}
