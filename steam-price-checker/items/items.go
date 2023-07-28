package steamchecker

type PriceSetter interface {
	SetPrice(value float64)
}

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

type CsgoTeamSticker struct {
	CsgoItem
	Team string
}

type CsgoPlayerSticker struct {
	CsgoItem
	Player string
}
