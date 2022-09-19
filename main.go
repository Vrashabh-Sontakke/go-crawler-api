package main

import (
	"fmt"
	"go-crawler-api/engineCrawler"
	"go-crawler-api/fenderCrawler"
	"go-crawler-api/frontBumperCrawler"
	"go-crawler-api/frontSuspensionCrawler"
	"go-crawler-api/hoodCrawler"
	"go-crawler-api/radiatorCrawler"
	"go-crawler-api/rearBumperCrawler"
	"go-crawler-api/rearQuarterPanelCrawler"
	"go-crawler-api/rearSuspensionCrawler"
	"go-crawler-api/transmissionCrawler"
)

func main() {

	fmt.Println("Choose any one option:")
	fmt.Println("1. Engine")
	fmt.Println("2. Transmission")
	fmt.Println("3. Front Bumper")
	fmt.Println("4. Rear Bumper")
	fmt.Println("5. Fender")
	fmt.Println("6. Rear Quarter Panels")
	fmt.Println("7. Front Suspension")
	fmt.Println("8. Rear Suspension")
	fmt.Println("9. Hood")
	fmt.Println("10. Radiator")

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
	case 5:
		fenderCrawler.Router()
	case 6:
		rearQuarterPanelCrawler.Router()
	case 7:
		frontSuspensionCrawler.Router()
	case 8:
		rearSuspensionCrawler.Router()
	case 9:
		hoodCrawler.Router()
	case 10:
		radiatorCrawler.Router()

	default:
		fmt.Println("Please enter correct part option number such as 1, 2, 3, 4, etc...")
	}
}
