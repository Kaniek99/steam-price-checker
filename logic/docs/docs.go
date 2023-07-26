package steamchecker

type SteamItem struct {
	Name  string
	Value float64
}

type CsgoTeamSticker struct {
	SteamItem
	Team string
}

type CsgoPlayerSticker struct {
	SteamItem
	Player string
}
