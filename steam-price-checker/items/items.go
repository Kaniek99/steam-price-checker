package steamchecker

type SteamItem struct {
	Name  string
	Value float64
	Wear  string
}

type CsgoTeamSticker struct {
	SteamItem
	Team string
}

type CsgoPlayerSticker struct {
	SteamItem
	Player string
}
