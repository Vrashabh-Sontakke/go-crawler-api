package main

import (
	"fmt"

	"go-crawler-api/engineCrawler"
	"go-crawler-api/frontBumperCrawler"
	"go-crawler-api/rearBumperCrawler"
	"go-crawler-api/transmissionCrawler"
)

func main() {

	fmt.Println("Choose any one option:")
	fmt.Println("1. Engine")
	fmt.Println("2. Transmission")
	fmt.Println("3. Front Bumper")
	fmt.Println("4. Rear Bumper")

	fmt.Print("Enter the part number:")

	var partNum int
	fmt.Scan(&partNum)

	switch partNum {
	case 1:
		engineCrawler.Router()
	case 2:
		transmissionCrawler.Router()
	case 3:
		frontBumperCrawler.Router()
	case 4:
		rearBumperCrawler.Router()
	default:
		fmt.Println("Please enter correct part option number such as 1, 2, 3, 4, etc...")
	}
}
