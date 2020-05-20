package crawler

// The palce that you can scrap everything!

import (
	"fmt"
	"strings"
	"top100-scrapy/pkg/model"
	"top100-scrapy/pkg/preference"

	"github.com/PuerkitoBio/goquery"
)

const UnavailbaleProduct = "This item is no longer available"

func ScrapeProductNames(opts *preference.Options) (names []string) {
	opts.Doc.Find("ol#zg-ordered-list li.zg-item-immersion").Each(func(i int, s *goquery.Selection) {
		var name string
		nameNode := s.Find("span.zg-text-center-align").Next()
		if len(nameNode.Nodes) == 1 {
			name = nameNode.Text()
		} else {
			name = UnavailbaleProduct
		}
		names = append(names, strings.TrimSpace(name))
	})
	return names
}

func ScrapeProducts(row *model.CategoryRow, opts *preference.Options) (set []*model.ProductRow, err error) {
	names := ScrapeProductNames(opts)
	if len(names) == 0 {
		err := fmt.Errorf("The names scraped from the url `%s` are empty, the category id stored into the DB is %d", row.Url, row.Id)
		return set, &EmptyError{row, err}
	}
	for i, name := range names {
		productRow := &model.ProductRow{
			Name:       name,
			Rank:       model.BuildRank(i, opts.Page),
			Page:       opts.Page,
			CategoryId: row.Id,
		}
		set = append(set, productRow)
	}
	return set, err
}

func ScrapeCategories(row *model.CategoryRow, opts *preference.Options) []*model.CategoryRow {
	// TODO: Track the error of the empty set scraped from the url.
	set := make([]*model.CategoryRow, 0)
	categoryRow := row
	opts.Doc.Find("#zg_browseRoot .zg_selected").Parent().Next().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "ul" {
			n := 0
			opts.Doc.Find("#zg_browseRoot .zg_selected").Parent().Next().Find("li a").Each(func(i int, s *goquery.Selection) {
				n++
				url, _ := s.Attr("href")
				path := model.BuildPath(n, categoryRow)
				row := &model.CategoryRow{
					Name:     s.Text(),
					Url:      url,
					Path:     path,
					ParentId: categoryRow.Id,
				}
				set = append(set, row)
			})
		}
	})
	return set
}
