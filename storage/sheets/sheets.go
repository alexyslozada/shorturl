package sheets

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/alexyslozada/shorturl/model"
)

type Sheets struct {
	service *sheets.Service
	logger  *zap.SugaredLogger
}

func New(cf string, logger *zap.SugaredLogger) (Sheets, error) {
	svc, err := createSheetsClient(cf)
	if err != nil {
		return Sheets{}, err
	}

	return Sheets{
		service: svc,
		logger:  logger,
	}, nil
}

func (s Sheets) AddRow(short *model.ShortURL, createdAt int64, spreadsheetID string) error {
	lastRow, err := getLastRow(s.service, spreadsheetID)
	if err != nil {
		s.logger.Errorw(fmt.Sprintf("Error getting last row, error was: %v", err))
		return err
	}

	sheetRange := fmt.Sprintf("datos!A%d:J%d", lastRow, lastRow)

	createdTime := time.Unix(createdAt, 0)

	values := []interface{}{
		short.ID,
		short.Short,
		short.RedirectTo,
		createdTime,
		createdTime.Year(),
		createdTime.Month(),
		createdTime.Day(),
		createdTime.Hour(),
		createdTime.Minute(),
		createdTime.Second(),
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}
	insertReq := s.service.Spreadsheets.Values.Append(spreadsheetID, sheetRange, valueRange)
	insertReq.ValueInputOption("USER_ENTERED")

	// Ejecutar la solicitud de inserción
	_, err = insertReq.Do()
	if err != nil {
		s.logger.Errorw(fmt.Sprintf("Error inserting record: %v", err))
		return err
	}

	return nil
}

func createSheetsClient(credentialsFile string) (*sheets.Service, error) {
	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets client: %v", err)
	}

	return client, nil
}

func getLastRow(client *sheets.Service, spreadsheetID string) (int64, error) {
	resp, err := client.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get values: %v", err)
	}

	if len(resp.Sheets) == 0 {
		return 0, fmt.Errorf("cannot get sheets from spreadsheetID %s, error: there are not sheets", spreadsheetID)
	}

	sheetsPropeties := resp.Sheets[0].Properties
	rowCount := sheetsPropeties.GridProperties.RowCount
	rowCount++

	return rowCount, nil
}
