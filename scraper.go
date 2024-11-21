package main

import (
	"fmt"
	"log"

	// importing Colly
	"encoding/csv"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

func main() {

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	var products []Product
	var visitedUrls sync.Map

	c.OnHTML("li.product", func(e *colly.HTMLElement) {

		product := Product{}

		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		products = append(products, product)

	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {

		nextPage := e.Attr("href")

		if _, found := visitedUrls.Load(nextPage); !found {
			fmt.Println("Scraping:", nextPage)

			visitedUrls.Store(nextPage, struct{}{})

			e.Request.Visit(nextPage)
		}
	})

	c.OnScraped(func(r *colly.Response) {

		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Cannot create file", err)
		}

		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}

		writer.Write(headers)

		for _, product := range products {
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			writer.Write(record)
		}

		defer writer.Flush()

		fmt.Println("Scraping finished, check file products.csv")

	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")

}
