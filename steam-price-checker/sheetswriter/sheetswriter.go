package steamchecker

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strconv"

	steamchecker "steam-price-checker/steam-price-checker/items"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsWriter struct {
	Items         []*steamchecker.SteamItem
	Service       *sheets.Service
	Context       *context.Context
	SheetID       int64
	SpreadsheetID string
}

func (sw *SheetsWriter) Authenticate() (*sheets.Service, error) {
	ctx := context.Background()
	sw.Context = &ctx
	credentials, err := base64.StdEncoding.DecodeString(os.Getenv("SHEETSKEY"))
	// log.Println(credentials)
	if err != nil {
		log.Println(os.Getenv("SHEETSKEY"))
		log.Println("decode")
		log.Println(err)
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(credentials, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Println("config")
		log.Println(err)
		return nil, err
	}

	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sw.Service = srv
	return srv, nil
}

func (sw *SheetsWriter) Init() {
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
	id, err := strconv.Atoi(os.Getenv("SHEETID"))
	if err != nil {
		log.Println(err)
		return
	}
	sw.SheetID = int64(id)
	sw.SpreadsheetID = os.Getenv("SPREADSHEETID")

	ctx := context.Background()
	sw.Context = &ctx
	credentials, err := base64.StdEncoding.DecodeString(os.Getenv("SHEETSKEY"))
	if err != nil {
		log.Fatalln(err)
	}

	config, err := google.JWTConfigFromJSON(credentials, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalln(err)
	}

	client := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)
	}
	sw.Service = srv
}

func (sw *SheetsWriter) InsertColumn(columnIndex int) {
	insertRequest := &sheets.InsertDimensionRequest{
		Range: &sheets.DimensionRange{
			SheetId:    sw.SheetID,
			Dimension:  "COLUMNS",
			StartIndex: int64(columnIndex),
			EndIndex:   int64(columnIndex + 1),
		},
	}

	updateRequest := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				InsertDimension: insertRequest,
			},
		},
	}

	_, err := sw.Service.Spreadsheets.BatchUpdate(sw.SpreadsheetID, &updateRequest).Do()
	if err != nil {
		log.Fatal("Unable to insert column in Google Sheets:", err)
	}
}
