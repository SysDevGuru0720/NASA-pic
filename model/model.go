package model

import (
	"encoding/json"
	"log"

	"github.com/SysDevGuru0720/NASA-pic/cache"
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

}
