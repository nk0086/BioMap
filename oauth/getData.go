package oauth

import (
	"net/http"

	"google.golang.org/api/calendar/v3"
)

// GetCalendarEventsはGoogle Calendar APIを使用して、
// カレンダーのイベントを取得する関数です
func GetCalendarEvents(client *http.Client) ([]*calendar.Event, error) {
	// Calendar APIのサービスを生成する
	service, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	// カレンダーのイベントを取得する
	events, err := service.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin("2022-01-01T00:00:00Z").TimeMax("2022-12-31T00:00:00Z").OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}
