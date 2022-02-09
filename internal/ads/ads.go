package ads

import (
	"fmt"
	"os"
	"path"

	"hz.code/neugls/ads/internal/config"
	"hz.code/neugls/ads/internal/database"
)

const AdsStatusValid = int(1)
const AdsStatusExpired = int(0)

type Ads struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` //webpage, image, video
	Position string `json:"position"`
	Content  string `json:"content"` //url, image url, video url
	Status   int    `json:"status"`
	Sort     int    `json:"sort"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
}

func List(start, limit uint) ([]*Ads, uint, error) {
	ads := []*Ads{}
	dbo := database.GetDBO()
	dbo.Clear().Select("*").From("#__ads").Order("sort desc")
	count, err := dbo.LoadList(&ads, start, limit)
	return ads, count, err
}

func ListValid() ([]*Ads, uint, error) {
	ads := []*Ads{}
	dbo := database.GetDBO()
	dbo.Clear().Select("*").From("#__ads").Where(fmt.Sprintf("status=%d", AdsStatusValid)).Order("sort desc")
	count, err := dbo.LoadList(&ads, 0, 10000)
	return ads, count, err
}

func Insert(ad *Ads) error {
	_, e := database.Insert("insert into #__ads(name, type, position, content, status, sort, created, updated) values(?,?,?,?,?,?,?,?)", ad.Name, ad.Type, ad.Position, ad.Content, ad.Status, ad.Sort, ad.Created, ad.Updated)
	return e
}

func Update(ad *Ads) error {
	_, e := database.DB().Exec("update #__ads set name=?, type=?, position=?, content=?, status=?, sort=?, updated=? where id=?", ad.Name, ad.Type, ad.Position, ad.Content, ad.Status, ad.Sort, ad.Updated, ad.ID)
	return e
}

func Load(id uint) (*Ads, error) {
	ad := &Ads{}
	e := database.DBGetByIDWithMemCache("ad", ad, "select * from #__ads where id=?", id)
	return ad, e
}

func Delete(id uint) error {
	ad, err := Load(id)
	if err != nil {
		return fmt.Errorf("load ad error %s", err.Error())
	}

	//remove ads content
	if ad.Type != "webpage" {
		if er := os.Remove(path.Join(config.V.DataDir, ad.Content)); er != nil {
			return fmt.Errorf("remove ads content error %s", er.Error())
		}
	}

	_, e := database.DB().Exec(database.Prefix("delete from #__ads where id=?"), id)
	return e
}

func Invalid(id uint) error {
	_, e := database.DB().Exec(database.Prefix("update #__ads set status=? where id=?"), AdsStatusExpired, id)
	return e
}

func Valid(id uint) error {
	_, e := database.DB().Exec(database.Prefix("update #__ads set status=? where id=?"), AdsStatusValid, id)
	return e
}
