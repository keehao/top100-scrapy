package crawler_test

import (
	"fmt"
	"net/http"
	"testing"
	"top100-scrapy/pkg/crawler"
	"top100-scrapy/pkg/model/product"
	"top100-scrapy/pkg/test"

	"github.com/PuerkitoBio/goquery"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

var (
	cassetteName = "crawler/base"
	url          = "https://www.amazon.com/Best-Sellers/zgbs/amazon-devices/ref=zg_bs_nav_0"
	doc          *goquery.Document
	t            *testing.T
)

func init() {
	cassettePath := fmt.Sprintf("%s/%s", test.FixturesUri, cassetteName)
	r, err := recorder.New(cassettePath)
	if err != nil {
		t.Errorf("Could not instantiate a recorder, error: %v", err)
	}
	defer r.Stop()

	// Create an HTTP client and inject the transport with the recorder.
	client := &http.Client{
		Transport: r, // Inject as transport!
	}
	resp, err := client.Get(url)
	if err != nil {
		t.Errorf("Failed to get the url, error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		factors := map[string]interface{}{
			"status_code": resp.StatusCode,
			"status":      resp.Status,
		}
		t.Errorf("The status of the code error occurs! Error: %v, factors: %v", err, factors)
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Errorf("Failed to return a document, error: %v", err)
	}
}

func TestScrapeProductNames(t *testing.T) {
	// Test the names of the top 5 products.
	expected := test.CannedScrapedProducts
	actual := crawler.New().WithDoc(doc).ScrapeProductNames()[:5]
	failedMsg := fmt.Sprintf("Failed, expected the names of the top 5 products: %s, got the names of the top 5 products: %s", expected, actual)
	assert.Equal(t, expected, actual, failedMsg)
}

func TestScrapeProducts(t *testing.T) {
	// Test the top 5 products
	products := product.NewRows()
	products.Set = test.CannedProductSet
	expected := products.RemovePointers(products.Set)
	products = crawler.New().WithDoc(doc).ScrapeProducts()
	actual := products.RemovePointers(products.Set)[:5]
	failedMsg := fmt.Sprintf("Failed, expected the top 5 products: %v, got the top 5 products: %v", expected, actual)
	assert.Equal(t, expected, actual, failedMsg)
}
