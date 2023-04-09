package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ResultProduct struct {
	KehManufacturer   string `json:"keh_manufacturer"`
	KehCoverage       string `json:"keh_coverage"`
	Discount          string `json:"discount"`
	HideGroupPrices   string `json:"hideGroupPrices"`
	ItemGroupId       string `json:"itemGroupId"`
	KehMaxFocalLength string `json:"keh_max_focal_length"`
	FreeShipping      string `json:"freeShipping"`
	StoreBaseCurrency string `json:"storeBaseCurrency"`
	Price             string `json:"price"`
	ToPrice           string `json:"toPrice"`
	ImageUrl          string `json:"imageUrl"`
	InStock           string `json:"inStock"`
	Currency          string `json:"currency"`
	ID                string `json:"id"`
	ImageHover        string `json:"imageHover"`
	SKU               string `json:"sku"`
	KehMaxAperture    string `json:"keh_max_aperture"`
	BasePrice         string `json:"basePrice"`
	KehProductType    string `json:"keh_product_type"`
	StartPrice        string `json:"startPrice"`
	KehMount          string `json:"keh_mount"`
	Image             string `json:"image"`
	DeliveryInfo      string `json:"deliveryInfo"`
	KehZoomPrime      string `json:"keh_zoom_prime"`
	HideAddToCart     string `json:"hideAddToCart"`
	SalePrice         string `json:"salePrice"`
	OldPrice          string `json:"oldPrice"`
	KehFilterSize     string `json:"keh_filter_size"`
	Swatches          struct {
		Swatch                     []string `json:"swatch"`
		LowestPrice                string   `json:"lowestPrice"`
		NumberOfAdditionalVariants string   `json:"numberOfAdditionalVariants"`
	} `json:"swatches"`
	Weight            string `json:"weight"`
	KlevuCategory     string `json:"klevu_category"`
	TotalVariants     string `json:"totalVariants"`
	GroupPrices       string `json:"groupPrices"`
	KehMinFocalLength string `json:"keh_min_focal_length"`
	URL               string `json:"url"`
	KehSystem         string `json:"keh_system"`
	Name              string `json:"name"`
	ShortDesc         string `json:"shortDesc"`
	Category          string `json:"category"`
	KehLensType       string `json:"keh_lens_type"`
}

type SearchResponse struct {
	Result []ResultProduct
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading config variables: ", err.Error())
	}

	url := os.Getenv("SEARCH_URL")

	for {
		result, err := search(url)

		if err != nil {
			log.Fatal(err)
		}

		if result != nil {
			sendEmail(result.URL)
			break
		} else {
			log.Println("The item was not found this time")
			time.Sleep(30 * time.Second)
		}
	}
}

func search(url string) (*ResultProduct, error) {
	log.Println("Searching...")
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		log.Println("Error decoding: ", err)
		return nil, err
	}

	wantedItem := ResultProduct{
		InStock:     "yes",
		KehCoverage: "aps-c & dx",
		KehMount:    "sony e mount",
	}

	for _, result := range searchResponse.Result {
		if result.InStock == wantedItem.InStock && result.KehCoverage == wantedItem.KehCoverage && result.KehMount == wantedItem.KehMount {
			log.Println("Item Found!")
			log.Printf("Manufacturer: %s\nCoverage: %s\nMax Aperture: %s\nIn Stock: %s\nMax Focal Length: %s\nMin Focal Length: %s\nURL: %s\nMount Type: %s\nSalePrice: %s\nOldPrice: %s\n\n",
				result.KehManufacturer, result.KehCoverage, result.KehMaxAperture, result.InStock, result.KehMaxFocalLength, result.KehMinFocalLength, result.URL, result.KehMount, result.SalePrice, result.OldPrice)
			return &result, nil
		}
	}
	return nil, nil
}

func sendEmail(itemUrl string) {
	log.Println("sending email: ", itemUrl)
	from := mail.NewEmail("GolangApp", os.Getenv("SENDER_EMAIL")) // Change to your verified sender
	subject := "The item is available"
	to := mail.NewEmail("Golang App", os.Getenv("RECEIVER_EMAIL")) // Change to your recipient
	plainTextContent := "The item you wanted is back online and available!"
	htmlContent := fmt.Sprintf("<strong>The item you wanted is back online and available! follow: <a href=\"%s\">This link</a></strong>", itemUrl)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent. ", response.StatusCode)
	}
}
