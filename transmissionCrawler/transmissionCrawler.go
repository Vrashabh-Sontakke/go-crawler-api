package transmissionCrawler

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
	router.HandleFunc("/{vin}", transmissionCrawler)

	fmt.Println("Running... localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}

func transmissionCrawler(w http.ResponseWriter, r *http.Request) {

	//var transmission []Transmission
	vars := mux.Vars(r)
	vin := vars["vin"]

	// INITIATE MEGA CRAWLERS

	// Crawler that gets link to actual transmission page
	link, err := getTransmissionLink(vin)
	if err != nil {
		log.Println("Error getting Transmission Page Link")
	}

	// Crawler for actual transmission Page
	transmissions, err := getTransmissions(*link)
	if err != nil {
		log.Println("Error getting transmissions.")
	}
	// Crawler to get individual transmission Data
	transmissionData, err := getIndividualTransmissionData(transmissions)
	if err != nil {
		log.Println("Error getting transmissions.")
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(transmissionData)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(j)

}

func getTransmissionLink(vin string) (*string, error) {
	var links []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColTwo]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Transmission" && h.ChildAttr("a", "href") != "" {
				e := Output{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, e)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
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

	transmissionLink := links[1]

	return &transmissionLink.URL, nil

}

func getTransmissions(url_suffix string) ([]Output, error) {
	var transmissions []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			e := Output{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			transmissions = append(transmissions, e)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return transmissions, nil

}

func getIndividualTransmissionData(links []Output) ([]Output, error) {

	var transmissions []Output

	for _, transmission := range links {
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

			transmissions = append(transmissions, e)

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
		})

		c.Visit("https://www.hollanderparts.com/" + transmission.URL)

		c.Wait()

	}

	fmt.Println(transmissions)

	return transmissions, nil

}
