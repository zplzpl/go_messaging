package model

import (
	"git.dev.tencent.com/pangolin/cmd/pangolin_portal/util"
	"git.dev.tencent.com/pangolin/pkg/shopee_bi_shop/model"
)

type ShopListRequest struct {
	Sn            string   `form:"sn" binding:"required"`
	Offset        int64    `form:"offset"`
	Limit         int64    `form:"limit" binding:"required"`
	ShopName      string   `form:"shopName"`
	AvgPriceStart float64  `form:"avgPriceStart"`
	AvgPriceEnd   float64  `form:"avgPriceEnd"`
	ShopLocation  string   `form:"shopLocation"`
	Catids        []uint64 `form:"catids"`
	CrawWord      []string `form:"crawWord"`
	AdsKeyword    []string `form:"adsKeyword"`
	Order         string   `form:"order"`
}

func (p *ShopListRequest) Validate() error {
	return nil
}

func (p *ShopListRequest) GetSn() string {
	return p.Sn
}

func (p *ShopListRequest) GetOffset() int64 {
	return p.Offset
}

func (p *ShopListRequest) GetLimit() int64 {
	return p.Limit
}

func (p *ShopListRequest) GetShopName() string {
	return p.ShopName
}

func (p *ShopListRequest) GetAvgPriceStart() int64 {
	return util.MYConvertToPrice(p.AvgPriceStart)
}

func (p *ShopListRequest) GetAvgPriceEnd() int64 {
	return util.MYConvertToPrice(p.AvgPriceEnd)
}

func (p *ShopListRequest) GetCatids() []uint64 {
	cids := make([]uint64, 0, len(p.Catids))
	for _, item := range p.Catids {
		if item > 0 {
			cids = append(cids, item)
		}
	}
	return cids
}

func (p *ShopListRequest) GetCrawWord() []string {
	cw := make([]string, 0, len(p.CrawWord))
	for _, item := range p.CrawWord {
		if item != "" {
			cw = append(cw, item)
		}
	}
	return cw
}

func (p *ShopListRequest) GetAdsKeyword() []string {
	akw := make([]string, 0, len(p.AdsKeyword))
	for _, item := range p.AdsKeyword {
		if item != "" {
			akw = append(akw, item)
		}
	}
	return akw
}

func (p *ShopListRequest) GetShopLocation() string {
	return p.ShopLocation
}

func (p *ShopListRequest) GetOrder() string {
	switch p.Order {
	case "ths_a":
		return "total_history_sold asc"
	case "ths_d":
		return "total_history_sold desc"
	case "hg_a":
		return "history_gmv asc"
	case "hg_d":
		return "history_gmv desc"
	case "ap_a":
		return "avg_price asc"
	case "ap_d":
		return "avg_price desc"
	case "fc_a":
		return "follower_count asc"
	case "fc_d":
		return "follower_count desc"
	}

	return ""
}

type ShopListResponse struct {
	Code    int64
	Message string
	Data    []*ShopListItem
	Pager   *Pager
}

type ShopListItem struct {
	ShopInfo     *ShopInfo     `json:"shop_info"`
	ShopSnapshot *ShopSnapshot `json:"shop_snapshot"`
}

type ShopInfo struct {
	UserId          uint64   `json:"user_id"`
	ShopId          uint64   `json:"shop_id"`
	UserName        string   `json:"user_name"`
	Name            string   `json:"name"`
	Portrait        string   `json:"portrait"`
	ShopUrl         string   `json:"shop_url"`
	FollowingCount  int64    `json:"following_count"`
	FollowerCount   int64    `json:"follower_count"`
	PreparationTime int64    `json:"preparation_time"`
	ItemCount       int64    `json:"item_count"`
	Description     string   `json:"description"`
	RatingGood      int64    `json:"rating_good"`
	RatingBad       int64    `json:"rating_bad"`
	RatingStar      float64  `json:"rating_star"`
	ShopLocation    string   `json:"shop_location"`
	CTime           int64    `json:"c_time"`
	TotalAvgStar    float64  `json:"total_avg_star"`
	Cover           string   `json:"cover"`
	ShopCovers      []string `json:"shop_covers"`
}

type ShopSnapshot struct {
	ShopId                uint64   `json:"shop_id"`
	CrawWord              []string `json:"craw_word"`
	CrawCatIds            []uint64 `json:"craw_cat_ids"`
	SoldProductCnt        uint64   `json:"sold_product_cnt"`
	TotalSold             uint64   `json:"total_sold"`
	Gmv                   float64  `json:"gmv"`
	AvgPrice              float64  `json:"avg_price"`
	CNYGmv                uint64   `json:"cny_gmv"`
	CNYAvgPrice           uint64   `json:"cny_avg_price"`
	ItemCount             uint64   `json:"item_count"`
	AdsKeyword            []string `json:"ads_keyword"`
	CrawCatNames          []string `json:"craw_cat_names"`
	HistorySoldProductCnt uint64   `json:"history_sold_product_cnt"`
	TotalHistorySold      uint64   `json:"total_history_sold"`
	HistoryGmv            float64  `json:"history_gmv"`
	TotalViewCnt          uint64   `json:"total_view_cnt"`
	FollowerCount         uint64   `json:"follower_count"`
}

type Pager struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

func RootToShopInfo(info *model.ShopInfo) *ShopInfo {

	p := &ShopInfo{
		UserId:          info.UserId,
		ShopId:          info.ShopId,
		UserName:        info.UserName,
		Name:            info.Name,
		Portrait:        util.ImgPathToImgUrl("", info.Portrait),
		ShopUrl:         util.GenShopUrl(info.UserName),
		FollowingCount:  info.FollowingCount,
		FollowerCount:   info.FollowerCount,
		PreparationTime: info.PreparationTime,
		ItemCount:       info.ItemCount,
		Description:     info.Description,
		RatingGood:      info.RatingGood,
		RatingBad:       info.RatingBad,
		RatingStar:      util.Round(info.RatingStar, 1),
		ShopLocation:    info.ShopLocation,
		CTime:           info.CTime,
		TotalAvgStar:    info.TotalAvgStar,
		Cover:           util.ImgPathToImgUrl("", info.Cover),
		ShopCovers:      info.ShopCovers,
	}

	for i, v := range p.ShopCovers {
		p.ShopCovers[i] = util.ImgPathToImgUrl("", v)
	}

	return p
}

func RootToShopSnapshot(sn *model.ShopSnapshot) *ShopSnapshot {

	p := &ShopSnapshot{
		ShopId:                sn.ShopId,
		CrawWord:              sn.CrawWord,
		CrawCatIds:            sn.CrawCatIds,
		SoldProductCnt:        sn.SoldProductCnt,
		TotalSold:             sn.TotalSold,
		Gmv:                   util.Round(util.PriceConvertToMY(sn.Gmv), 2),
		AvgPrice:              util.Round(util.PriceConvertToMY(sn.AvgPrice), 2),
		CNYGmv:                util.PriceConvertMYToCNY(sn.Gmv),
		CNYAvgPrice:           util.PriceConvertMYToCNY(sn.AvgPrice),
		ItemCount:             sn.ItemCount,
		AdsKeyword:            sn.AdsKeyword,
		CrawCatNames:          sn.CrawCatNames,
		HistorySoldProductCnt: sn.HistorySoldProductCnt,
		TotalHistorySold:      sn.TotalHistorySold,
		HistoryGmv:            util.Round(util.PriceConvertToMY(sn.HistoryGmv), 2),
		TotalViewCnt:          sn.TotalViewCnt,
		FollowerCount:         sn.FollowerCount,
	}

	return p
}

func RootToShopListItem(items []*model.ShopListItem) []*ShopListItem {

	list := make([]*ShopListItem, 0, len(items))

	for _, item := range items {

		tmp := &ShopListItem{
			ShopInfo:     RootToShopInfo(item.ShopInfo),
			ShopSnapshot: RootToShopSnapshot(item.ShopSnapshot),
		}

		list = append(list, tmp)
	}

	return list
}
