package steamchecker

type SteamItem struct {
	Name       string
	HashedName string
	Price      float64
	GameID     int
}

func (item *SteamItem) SetPrice(value float64) {
	item.Price = value
}

type CsgoItem struct {
	SteamItem
	Wear string
}

func (item *CsgoItem) SetPrice(value float64) {
	item.SteamItem.Price = value
}
