package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/SysDevGuru0720/NASA-pic/cache"
	"github.com/SysDevGuru0720/NASA-pic/config"
	"github.com/gomodule/redigo/redis"
)

type Picture struct {
	Date        string `json:"date"`
	URL         string `json:"url"`
	CopyRight   string `json:"copyright"`
	Explanation string `json:"explanation"`
	Title       string `json:"title"`
	HDURL       string `json:"hdurl"`
}

func (p *Picture) GetPic() error {
	conn, err := cache.GetConn()
	if err != nil {
		log.Printf("nasapic: faile to get connection to cahce:%v\n", err)
		return p.fetchPicture(p.Date, conn)
	}

	s, err := cache.Get(conn, p.Date)
	if err != nil {
		log.Printf("nasapic: failed to get from caache: %v\n", err)
		return p.fetchPicture(p.Date, conn)
	}

	err = json.Unmarshal([]byte(s), &p)
	if err != nil {
		return err
	}

	conn.Close()

	return nil
}

func (p *Picture) fetchPicture(date string, conn redis.Conn) error {
	defer conn.Close()
	client := &http.Client{}

	req, err := http.NewRequest("GET", config.Global.Config.NasaURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("date", date)
	q.Add("api_key", os.Getenv("NASA_API_KEY"))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &p)
	if err != nil {
		return err
	}

	err = cache.Set(conn, date, string(respBody))
	if err != nil {
		log.Printf("nasapic: failed to set cache: %v\n", err)
	}

	return nil
}
