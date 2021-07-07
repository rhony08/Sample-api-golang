package connection

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"os"
	"strconv"
	"sync"
)

//HTTPAPI for creating request or do request
type HTTPAPI struct {
	Method            string            // method
	URL               string            // url
	URIParams         map[string]string // params of API
	Body              []byte            // content body of API
	Headers           map[string]string // request header
	RestfulParams     map[string]string // restful params for restful API
	BasicAuthUserName string
	BasicAuthPassword string
}

// Options for client
type Options struct {
	Timeout     int
	Environment string
}

// Client struct
type Client struct {
	client  http.Client
	options *Options

	mtx sync.RWMutex
}

//ErrDefaultClientIsNil client is nil
var ErrDefaultClientIsNil = errors.New("Default client is nil")

//ErrResponseIsNil Response is nil
var ErrResponseIsNil = errors.New("Response is nil")

var defaultClient *Client

func getDefaultClient() *Client {
	if defaultClient == nil {
		timeoutConfig := os.Getenv("TIMEOUT")
		var err error
		var timeoutConfigNumber int

		if timeoutConfig != "" {
			timeoutConfigNumber, err = strconv.Atoi(os.Getenv("TIMEOUT"))
			if err != nil {
				timeoutConfigNumber = 3
			}
		}

		defaultClient = &Client{
			client: http.Client{Timeout: time.Second * time.Duration(timeoutConfigNumber)},
		}
	}
	return defaultClient
}

func doRequest(h HTTPAPI, ctxParam context.Context) (*http.Response, error) {
	var response *http.Response

	request, err := getDefaultClient().NewRequest(h)
	if err != nil {
		return response, err
	}

	request = request.WithContext(ctxParam)
	response, err = getDefaultClient().client.Do(request)

	//
	// avoid reponse and err nil booth (if response nil and no error)
	// will generate new error
	//
	if nil == response && nil == err {
		err = ErrResponseIsNil
	}
	return response, err
}

// NewRequest for creating new http request
func (c *Client) NewRequest(h HTTPAPI) (*http.Request, error) {
	var (
		request *http.Request
		buff    *bytes.Buffer
		rawURL  string
	)

	if h.URL == "" || len(h.URL) == 0 {
		return request, fmt.Errorf("URL is required")
	}

	h.Method = strings.ToUpper(h.Method)
	if h.Method != "POST" && h.Method != "GET" && h.Method != "PUT" && h.Method != "DELETE" {
		return request, fmt.Errorf("Unsupported method %s", h.Method)
	}

	u, err := url.Parse(h.URL)
	if err != nil {
		return request, err
	}

	val := u.Query()
	for key, value := range h.URIParams {
		val.Add(key, value)
	}
	encodedVal := val.Encode()
	rawURL = u.String()

	if len(h.URIParams) > 0 {
		rawURL += "?" + encodedVal
	}

	switch h.Method {
	case "POST", "PUT", "DELETE":
		buff = bytes.NewBuffer(h.Body)
	}

	//if buff == nil, will produce nil pointer interface
	if buff != nil {
		request, err = http.NewRequest(h.Method, rawURL, buff)
	} else {
		request, err = http.NewRequest(h.Method, rawURL, nil)
	}

	if err != nil {
		return request, err
	}

	for key, value := range h.Headers {
		request.Header.Add(key, value)
	}

	if h.BasicAuthUserName != "" && h.BasicAuthPassword != "" {
		request.SetBasicAuth(h.BasicAuthUserName, h.BasicAuthPassword)
	}

	request.Header.Add("Content-Length", strconv.Itoa(len(encodedVal)))

	// add monit url to context
	request = request.WithContext(context.Background())
	return request, err
}

// DoRequestWithContext will request data using provided context
func DoRequestWithContext(ctx context.Context, h HTTPAPI) (*http.Response, error) {
	return doRequest(h, ctx)
}
