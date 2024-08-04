package upbit

import (
	"coinquant/pkg/upbit/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Client struct {
	token     string
	Url       string
	callCount int
}

const ApiUrl = "https://api.upbit.com"

func MakeClient() (*Client, error) {
	accessKey := os.Getenv("UPBIT_OPEN_API_ACCESS_KEY")
	secretKey := os.Getenv("UPBIT_OPEN_API_SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("one or more environment variables are missing: UPBIT_OPEN_API_ACCESS_KEY, UPBIT_OPEN_API_SECRET_KEY, UPBIT_OPEN_API_SERVER_URL")
	}

	payload := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	client := &Client{
		token: tokenString,
		Url:   ApiUrl,
	}

	return client, nil
}

func (c *Client) GetCandleChart(market string, candleType model.CandleType, end time.Time, count int) ([]model.Candle, error) {
	endpoint := fmt.Sprintf("%s/v1/candles/%s", c.Url, candleType)
	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
		return nil, err
	}
	q := u.Query()
	q.Set("market", market)
	q.Set("count", strconv.Itoa(count))
	q.Set("to", end.Format("2006-01-02T15:04:05"))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
		return nil, err
	}
	c.callCount++
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil, err
	}
	fmt.Printf("call count: %d, %s\n", c.callCount, u.String())
	var candles []model.Candle
	err = json.Unmarshal(body, &candles)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
		return nil, err
	}
	return candles, nil
}
