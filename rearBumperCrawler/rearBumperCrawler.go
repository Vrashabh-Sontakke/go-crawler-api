package rearBumperCrawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type Output struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Price    string `json:"price,optional"`
	Shipping string `json:"shipping,optional"`
	Img      string `json:"img,optional"`
	Grade    string `json:"grade,optional"`
}

func Router() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{vin}", rearBumperCrawler)

	fmt.Println("Running... localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}

func rearBumperCrawler(w http.ResponseWriter, r *http.Request) {

	//var rearBumpers []RearBumper
	vars := mux.Vars(r)
	vin := vars["vin"]

	// INITIATE MEGA CRAWLERS

	// Crawler that gets link to actual rearBumper page
	link, err := getRearBumperLink(vin)
	if err != nil {
		log.Println("Error getting RearBumper Page Link")
	}

	// Crawler for actual RearBumper Page
	rearBumpers, err := getRearBumpers(*link)
	if err != nil {
		log.Println("Error getting rearBumpers.")
	}
	// Crawler to get individual rearBumper Data
	rearBumperData, err := getIndividualRearBumperData(rearBumpers)
	if err != nil {
		log.Println("Issue getting rearBumpers.")
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(rearBumperData)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(j)

}

func getRearBumperLink(vin string) (*string, error) {
	var links []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColTwo]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Rear Body" && h.ChildAttr("a", "href") != "" {
				e := Output{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, e)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
			} else {
				c.OnHTML("div[class=searchColThree]", func(h *colly.HTMLElement) {
					h.ForEach("div", func(i int, h *colly.HTMLElement) {
						if h.ChildText("a") == "Rear Bumper" && h.ChildAttr("a", "href") != "" {
							e := Output{
								Name: h.ChildText("a"),
								URL:  h.ChildAttr("a", "href"),
							}
							links = append(links, e)
							c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
						}
					})

				})
			}
		})

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.PostMultipart("https://www.hollanderparts.com/Home", map[string][]byte{
		"hdnVIN": []byte(vin),
	})
	c.Wait()

	rearBumperLink := links[1]

	return &rearBumperLink.URL, nil

}

func getRearBumpers(url_suffix string) ([]Output, error) {
	var rearBumpers []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			e := Output{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			rearBumpers = append(rearBumpers, e)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return rearBumpers, nil

}

func getIndividualRearBumperData(links []Output) ([]Output, error) {

	var rearBumpers []Output

	for _, rearBumper := range links {
		c := colly.NewCollector(
			colly.AllowURLRevisit(),
		)

		c.OnHTML("div[class=individualPartHolder]", func(h *colly.HTMLElement) {

			name := strings.Split(h.Response.Request.URL.String(), "/")
			price := h.ChildText("div[class=partPrice]")
			shipping := h.ChildText("div[class=partShipping]")
			img := h.ChildAttr("img", "src")
			grade := h.ChildText("div[class=gradeText]")

			e := Output{
				Name:     name[len(name)-1],
				URL:      h.Response.Request.URL.String(),
				Grade:    grade,
				Img:      img,
				Price:    price,
				Shipping: shipping,
			}

			rearBumpers = append(rearBumpers, e)

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
		})

		c.Visit("https://www.hollanderparts.com/" + rearBumper.URL)

		c.Wait()

	}

	fmt.Println(rearBumpers)

	return rearBumpers, nil

}
