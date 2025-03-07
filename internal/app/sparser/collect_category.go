package sparser

import (
	"alice088/sparser/internal/pkg/dto"
	"alice088/sparser/internal/pkg/samokat"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
)

func CollectCategories(log *zerolog.Logger, sessionData *dto.SessionData, geo *dto.GEO) [17]*dto.Category {
	categories := [17]*dto.Category{}
	client := &http.Client{}
	jsonGeo, err := json.Marshal(geo)

	if err != nil {
		log.Fatal().Err(err).
			Int("Geo id", geo.ID).
			Msg("Failed to marshal geography")
	}

	log.Debug().Str("JsonGeo", string(jsonGeo)).Send()

	req, err := http.NewRequest("GET", fmt.Sprintf(samokat.CATEGORIES_LIST, sessionData.ShowcaseID), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", samokat.USER_AGENT)
	req.Header.Set("authorization", sessionData.AuthToken)
	req.Header.Set("Cookie", "SELECTED_ADDRESS_KEY="+url.QueryEscape(string(jsonGeo)))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed send request")
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal().Int("StatusCode", resp.StatusCode).Msg("Not ok status code from response during collect category")
	}
	defer closeBody(log, err, resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed read response body during collect category")
	}

	rawCategoryList := gjson.GetBytes(body, "data.0")

	rawCategoryList.ForEach(func(index, value gjson.Result) bool {
		if index.Int() >= 1 && index.Int() <= 17 {
			category := &dto.Category{
				Id:       index.Get("id").String(),
				Name:     index.Get("name").String(),
				ParentId: index.Get("parent_id").String(),
			}
			categories[index.Int()] = category
		}
		return true
	})

	return categories
}

func closeBody(log *zerolog.Logger, err error, resp *http.Response) {
	func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed close response body during collect category")
		}
	}(resp.Body)
}
