package connection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type OMDBInterface interface {
	SearchMovieData(key string, pagination int) ([]Movie, error)
	GetDetailMovie(imdbID string) (Movie, error)
}

type Movie struct {
	Title  string
	Year   string
	Type   string
	Poster string
	ID     string `json:"imdbID"`
}

type resultData struct {
	Search       []Movie
	TotalResults string `json:"totalResults"`
}

type omdbConnection struct {
	baseURL string
	API_KEY string
}

const (
	paramSearch     = "s"
	paramID         = "i"
	paramPagination = "page"
	paramApiKey     = "apikey"
)

func newOmdbConnection(apikey string, baseURL string) OMDBInterface {
	return &omdbConnection{
		API_KEY: apikey,
		baseURL: baseURL,
	}
}

func (outcallOmdb *omdbConnection) SearchMovieData(key string, pagination int) (response []Movie, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	params := make(map[string]string)

	if key != "" {
		params[paramSearch] = key
	}

	if pagination != 0 {
		params[paramPagination] = fmt.Sprintf("%d", pagination)
	}

	params[paramApiKey] = outcallOmdb.API_KEY

	resp, err := DoRequestWithContext(ctx, HTTPAPI{
		Method:    http.MethodGet,
		URL:       outcallOmdb.baseURL,
		URIParams: params,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}

	// fmt.Println(string(content))
	responseData := new(resultData)
	err = json.Unmarshal(content, &responseData)
	if err != nil {
		return response, err
	}

	return responseData.Search, nil
}

func (outcallOmdb *omdbConnection) GetDetailMovie(imdbID string) (response Movie, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	params := make(map[string]string)

	if imdbID != "" {
		params[paramID] = imdbID
	} else {
		return response, errors.New("Invalid ID")
	}

	params[paramApiKey] = outcallOmdb.API_KEY

	resp, err := DoRequestWithContext(ctx, HTTPAPI{
		Method:    http.MethodGet,
		URL:       outcallOmdb.baseURL,
		URIParams: params,
	})

	if resp != nil {
		defer resp.Body.Close()
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	err = json.Unmarshal(content, &response)
	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	return response, nil
}
