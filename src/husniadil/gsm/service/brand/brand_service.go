package brand

import (
	"errors"
	"strconv"
	"strings"

	. "husniadil/gsm/model"
	web_scraper "husniadil/gsm/util/scraper"

	"github.com/PuerkitoBio/goquery"
)

const Endpoint = "http://www.gsmarena.com/makers.php3"

func getSlug(s *goquery.Selection) (slug string, err error) {
	tmpSlug, exists := s.Attr("href")

	if exists {
		slug = strings.TrimSuffix(tmpSlug, ".php")
	} else {
		err = errors.New("Cannot resolve slug")
	}

	return
}

func getID(s *goquery.Selection) (id int, err error) {
	slug, err := getSlug(s)

	if err != nil {
		return
	}

	tmpID := strings.Split(slug, "-")
	tmpID2 := tmpID[len(tmpID)-1]

	id, err = strconv.Atoi(tmpID2)

	return
}

func getName(s *goquery.Selection) (name string, err error) {
	tmpName, err := s.Html()

	if err != nil {
		return
	}

	tmpName2 := strings.Split(tmpName, "<br/>")
	if len(tmpName2) > 0 {
		name = tmpName2[0]
	}

	return
}

func getNumberOfDevices(s *goquery.Selection) (numberOfDevices int, err error) {
	numberOfDevices, err = strconv.Atoi(strings.Split(s.Find("span").Text(), " ")[0])

	return
}

func newBrand(s *goquery.Selection) (brand Brand, err error) {
	id, err := getID(s)

	if err != nil {
		return
	}

	slug, err := getSlug(s)

	if err != nil {
		return
	}

	name, err := getName(s)

	if err != nil {
		return
	}

	numberOfDevices, err := getNumberOfDevices(s)

	if err != nil {
		return
	}

	brand = Brand{
		ID:              id,
		Name:            name,
		Slug:            slug,
		NumberOfDevices: numberOfDevices,
	}

	return
}

func GetAllBrands() (brands []Brand, err error) {
	scraper := web_scraper.NewScraper()
	document, err := scraper.Scrap(Endpoint)

	if err != nil {
		return
	}

	links := document.Find(".st-text").Find("a")
	for i := 0; i < links.Length(); i++ {
		brand, err := newBrand(links.Eq(i))

		if err != nil {
			return nil, err
		}

		brands = append(brands, brand)
	}

	return
}
