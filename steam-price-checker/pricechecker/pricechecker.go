package steamchecker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	steamchecker "steam-price-checker/steam-price-checker/items"
)

type MarketPriceOverview struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

type PriceChecker struct {
	APIKey string // probably unnecessary
	Items  []string
	// Items           []steamchecker.SteamItem
	// PlayersStickers []steamchecker.CsgoPlayerSticker
	// TeamsStickers   []steamchecker.CsgoTeamSticker
	Test string
}

func (pc *PriceChecker) SetItemsToCheck() {
	itemsID := []string{"Sticker%20%7C%20ropz%20%28Holo%29%20%7C%20Paris%202023", "Sticker%20%7C%20Twistzz%20%28Holo%29%20%7C%20Paris%202023"}
	pc.Items = append(pc.Items, itemsID...)
}

func (pc *PriceChecker) GetAutographsPrices() {
	log.Println("data collection has begun")
	pc.APIKey = os.Getenv("STEAMKEY")
	if len(pc.APIKey) == 0 {
		log.Fatal("Something went wrong with setting steam API key")
	}
	for _, elem := range pc.Items {
		// currency: 1 - $, 3 - €, 7 - PLN
		url := fmt.Sprintf("https://steamcommunity.com/market/priceoverview/?currency=3&appid=730&market_hash_name=%v", elem) // 730 is AppID of Counter-Strike: Global Offensive
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request for %v failed with status code: %v", elem, resp.StatusCode)
		}
		result := MarketPriceOverview{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
		} else {
			fmt.Println(result)
			sticker := steamchecker.CsgoPlayerSticker{SteamItem: steamchecker.SteamItem{Name: elem, Value: 5.13}, Player: "test"}
			fmt.Println(sticker)
		}
	}
	log.Println("data collection completed")
}
