package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

func ProductCardParse(bow *browser.Browser) {
	bow.Dom().Find("div.product-card__content").Each(func(_ int, s *goquery.Selection) {
		new_price := s.Find("div.product-unit-prices__actual-wrapper").Find("span.product-price__sum-rubles").Text()
		old_price := s.Find("div.product-unit-prices__old-wrapper").Find("span.product-price__sum-rubles").Text()
		disc := s.Find("div.product-discount.nowrap.catalog-2-level-product-card__icon-discount.style--catalog-2-level-product-card").Text()
		name := s.Find("span.product-card-name__text").Text()
		fmt.Println(name, "Старая цена:", old_price, "Новая цена:", new_price, "Скидка:", disc)
	})
}

func AllCategoricalProductParse(mapa map[string]string) {
	for k, v := range mapa {
		fmt.Printf("%s\n%s\n%s\n", strings.Repeat("=", 50), k, strings.Repeat("=", 50))
		CategoricalProductParse(v)
	}
}

func UrlParse(url string) (*browser.Browser, error) {
	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		return nil, err
	}

	return bow, nil
}

func CategoricalProductParse(url string) {
	var num int

	bow, err := UrlParse(url)
	if err != nil {
		panic(err)
	}

	pagination := bow.Dom().Find("a.v-pagination__item.catalog-paginate__item").Each(func(_ int, s *goquery.Selection) {
		num, _ = strconv.Atoi(s.Text())
	})

	if pagination.Length() > 0 {
		for i := 1; i <= num; i++ {
			ProductCardParse(bow)

			if i < num {
				bow, _ = UrlParse(fmt.Sprintf("%s?page=%d", url, i+1))
			}
		}
	} else {
		ProductCardParse(bow)
	}
}

func CategoryParse(url string) map[string]string {
	urlMap := make(map[string]string, 0)
	bow, err := UrlParse(url)
	if err != nil {
		panic(err)
	}

	bow.Dom().Find("a.catalog-filters-categories__item.level-first.is-bold").Each(func(_ int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		spanText := s.Find("span").Text()

		if spanText == "" {
			spanText = s.Text()
		}

		urlMap[spanText] = "https://online.metro-cc.ru" + link
	})

	return urlMap
}
