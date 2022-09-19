package rearSuspensionCrawler

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
	router.HandleFunc("/{vin}", partCrawler)

	fmt.Println("Running... localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}

func partCrawler(w http.ResponseWriter, r *http.Request) {

	//var parts []Part
	vars := mux.Vars(r)
	vin := vars["vin"]

	// INITIATE MEGA CRAWLERS

	// Crawler that gets link to actual part page
	link, err := getPartLink(vin)
	if err != nil {
		log.Println("Error getting Part Page Link")
	}

	// Crawler for actual Part Page
	parts, err := getParts(*link)
	if err != nil {
		log.Println("Error getting parts.")
	}
	// Crawler to get individual part Data
	partData, err := getIndividualPartData(parts)
	if err != nil {
		log.Println("Issue getting parts.")
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(partData)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(j)

}

func getPartLink(vin string) (*string, error) {
	var links []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColTwo]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Suspension-Steering" && h.ChildAttr("a", "href") != "" {
				e := Output{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, e)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
			} else {
				c.OnHTML("div[class=searchColThree]", func(h *colly.HTMLElement) {
					h.ForEach("div", func(i int, h *colly.HTMLElement) {
						if h.ChildText("a") == "Rear Suspension" && h.ChildAttr("a", "href") != "" {
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

	partLink := links[1]

	return &partLink.URL, nil

}

func getParts(url_suffix string) ([]Output, error) {
	var parts []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			e := Output{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			parts = append(parts, e)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return parts, nil

}

func getIndividualPartData(links []Output) ([]Output, error) {

	var parts []Output

	for _, part := range links {
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

			parts = append(parts, e)

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
		})

		c.Visit("https://www.hollanderparts.com/" + part.URL)

		c.Wait()

	}

	fmt.Println(parts)

	return parts, nil

}
