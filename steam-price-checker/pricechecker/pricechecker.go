package steamchecker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	steamchecker "steam-price-checker/steam-price-checker/items"
)

type MarketPriceOverview struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

type PriceChecker struct {
	Items     []*steamchecker.SteamItem
	CsgoItems []*steamchecker.CsgoItem
	// CsgoPlayersStickers []steamchecker.CsgoPlayerSticker
	// CsgoTeamsStickers   []steamchecker.CsgoTeamSticker
}

func (pc *PriceChecker) SetItemsToCheck() {
	file, err := os.Open("items.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := steamchecker.SteamItem{}
		line := scanner.Text()
		data := strings.Split(line, " || ")
		item.GameID, err = strconv.Atoi(strings.TrimSpace(data[0]))
		if err != nil {
			log.Println("GameID should be integer value, got: ", data[0])
		}
		item.Name = data[1]
		item.HashedName = url.QueryEscape(item.Name)
		if item.GameID == 730 { //730 is appiD of Counter-Strike: Global Offensive
			// how to differentiate stickers types? I mean normal sticker, team or player sticker
			csgoItem := steamchecker.CsgoItem{SteamItem: item}
			pc.CsgoItems = append(pc.CsgoItems, &csgoItem)
		} else {
			pc.Items = append(pc.Items, &item)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	log.Println("data from file items.txt uploaded correctly")
	log.Println(pc.CsgoItems)
	log.Println(pc.Items)
}

func SetPrice(value float64, obj steamchecker.PriceSetter) {
	obj.SetPrice(value)
}

func (pc *PriceChecker) GetPrice(url string) (MarketPriceOverview, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request for %v failed with status code: %v", url, resp.StatusCode)
	}
	result := MarketPriceOverview{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding JSON response:", err)
		return result, err
	}
	result.LowestPrice = result.LowestPrice[:len(result.LowestPrice)-3] // remove € sign which length is 3
	replacer := strings.NewReplacer("-", "0", ",", ".")
	result.LowestPrice = replacer.Replace(result.LowestPrice)
	return result, nil
}

func (pc *PriceChecker) SetPrices() {
	log.Println("data collection has begun")
	for _, elem := range pc.Items {
		// currency: 1 - $, 3 - €, 7 - PLN
		url := fmt.Sprintf("https://steamcommunity.com/market/priceoverview/?currency=3&appid=%v&market_hash_name=%v", elem.GameID, elem.HashedName)
		result, err := pc.GetPrice(url)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
		}
		price, err := strconv.ParseFloat(result.LowestPrice, 64)
		if err != nil {
			log.Println(err)
		}
		SetPrice(price, elem)
		fmt.Println(elem)
	}
	for _, elem := range pc.CsgoItems {
		url := fmt.Sprintf("https://steamcommunity.com/market/priceoverview/?currency=3&appid=730&market_hash_name=%v", elem.HashedName)
		result, err := pc.GetPrice(url)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
		}
		price, err := strconv.ParseFloat(result.LowestPrice, 64)
		if err != nil {
			log.Println(err)
		}
		SetPrice(price, elem)
		fmt.Println(elem)
	}
	log.Println("data collection completed")
	log.Println(pc.CsgoItems)
	log.Println(pc.Items)
}
