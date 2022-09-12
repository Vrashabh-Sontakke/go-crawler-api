package frontBumperCrawler

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
	router.HandleFunc("/{vin}", frontBumperCrawler)

	fmt.Println("Running... localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}

func frontBumperCrawler(w http.ResponseWriter, r *http.Request) {

	//var frontBumpers []FrontBumper
	vars := mux.Vars(r)
	vin := vars["vin"]

	// INITIATE MEGA CRAWLERS

	// Crawler that gets link to actual frontBumper page
	link, err := getFrontBumperLink(vin)
	if err != nil {
		log.Println("Error getting FrontBumper Page Link")
	}

	// Crawler for actual FrontBumper Page
	frontBumpers, err := getFrontBumpers(*link)
	if err != nil {
		log.Println("Error getting frontBumpers.")
	}
	// Crawler to get individual frontBumper Data
	frontBumperData, err := getIndividualFrontBumperData(frontBumpers)
	if err != nil {
		log.Println("Issue getting frontBumpers.")
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.Marshal(frontBumperData)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(j)

}

func getFrontBumperLink(vin string) (*string, error) {
	var links []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			if h.ChildText("a") == "Front Body" && h.ChildAttr("a", "href") != "" {
				e := Output{
					Name: h.ChildText("a"),
					URL:  h.ChildAttr("a", "href"),
				}
				links = append(links, e)
				c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
			} else {
				c.OnHTML("div[class=searchColTwo]", func(h *colly.HTMLElement) {
					h.ForEach("div", func(i int, h *colly.HTMLElement) {
						if h.ChildText("a") == "Front Bumper" && h.ChildAttr("a", "href") != "" {
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

	frontBumperLink := links[1]

	return &frontBumperLink.URL, nil

}

func getFrontBumpers(url_suffix string) ([]Output, error) {
	var frontBumpers []Output
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("div[class=searchColOne]", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {
			e := Output{
				Name: h.ChildText("a"),
				URL:  h.ChildAttr("a", "href"),
			}
			frontBumpers = append(frontBumpers, e)
			c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
		})

	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit("https://www.hollanderparts.com/" + url_suffix)
	c.Wait()

	return frontBumpers, nil

}

func getIndividualFrontBumperData(links []Output) ([]Output, error) {

	var frontBumpers []Output

	for _, frontBumper := range links {
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

			frontBumpers = append(frontBumpers, e)

		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("visiting", r.URL.String())
		})

		c.Visit("https://www.hollanderparts.com/" + frontBumper.URL)

		c.Wait()

	}

	fmt.Println(frontBumpers)

	return frontBumpers, nil

}
