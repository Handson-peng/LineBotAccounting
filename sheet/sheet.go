package sheet

import (
	"context"
	"log"

	"google.golang.org/api/sheets/v4"
)

var SpreadsheetId string

type Service sheets.Service

func (srv *Service) CreatSheetDemo(title string) {
	srv.CreatSheet(title)
	srv.ValueUpdate(title, "A1:G1", []string{"時間", "金額", "分類", "內容", "", "總計", "=SUM(B:B)"})
}

func (srv *Service) CreatSheet(title string) {

	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: title,
			},
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&req},
	}

	resp, err := srv.Spreadsheets.BatchUpdate(SpreadsheetId, rbb).Context(context.Background()).Do()
	if err != nil {
		log.Fatal(err)
	}

	req = sheets.Request{
		UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
			Properties: &sheets.SheetProperties{
				SheetId: resp.Replies[0].AddSheet.Properties.SheetId,
				Index:   0,
			},
			Fields: "Index",
		},
	}
	_, err = srv.Spreadsheets.BatchUpdate(SpreadsheetId, rbb).Context(context.Background()).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Service) AppendRow(title string, str []string) {

	interface1D := make([]interface{}, len(str))
	for i, v := range str {
		interface1D[i] = v
	}

	interface2D := [][]interface{}{}
	interface2D = append(interface2D, interface1D)
	vr := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         interface2D,
	}
	_, err := srv.Spreadsheets.Values.Append(SpreadsheetId, title+"!A2", &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		srv.CreatSheetDemo(title)
		srv.AppendRow(title, str)
	}
}

func (srv *Service) ValueGet(title, coord string) [][]any {
	resp, err := srv.Spreadsheets.Values.Get(SpreadsheetId, title+"!"+coord).Do()
	if err != nil {
		srv.CreatSheetDemo(title)
        return srv.ValueGet(title, coord)
	}
	return resp.Values
}

func (srv *Service) ValueUpdate(title, coord string, str []string) {
	interface1D := make([]interface{}, len(str))
	for i, v := range str {
		interface1D[i] = v
	}

	interface2D := [][]interface{}{}
	interface2D = append(interface2D, interface1D)
	vr := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         interface2D,
	}
	_, err := srv.Spreadsheets.Values.Update(SpreadsheetId, title+"!"+coord, &vr).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatal(err)
	}
}
