package favor

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/yunlyz/wxpay"
)

type CouponService wxpay.Service

// Coupon represent a wechat pay coupon
type Coupon struct {
    StockCreatorMchid string `json:"stock_creator_mchid"`
    StockID           string `json:"stock_id"`
    CouponID          string `json:"coupon_id"`
    CutToMessage      struct {
        SinglePriceMax int `json:"single_price_max"`
        CutToPrice     int `json:"cut_to_price"`
    } `json:"cut_to_message"`
    CouponName              string      `json:"coupon_name"`
    Status                  string      `json:"status"`
    Description             string      `json:"description"`
    CreateTime              time.Time   `json:"create_time"`
    CouponType              string      `json:"coupon_type"`
    NoCash                  interface{} `json:"no_cash"`
    AvailableBeginTime      time.Time   `json:"available_begin_time"`
    AvailableEndTime        time.Time   `json:"available_end_time"`
    Singleitem              interface{} `json:"singleitem"`
    NormalCouponInformation struct {
        CouponAmount       int `json:"coupon_amount"`
        TransactionMinimum int `json:"transaction_minimum"`
    } `json:"normal_coupon_information"`
    ConsumeInformation struct {
        ConsumeTime   time.Time `json:"consume_time"`
        ConsumeMchid  string    `json:"consume_mchid"`
        TransactionID string    `json:"transaction_id"`
        GoodsDetail   []struct {
            GoodsID        string `json:"goods_id"`
            Quantity       int    `json:"quantity"`
            Price          int    `json:"price"`
            DiscountAmount int    `json:"discount_amount"`
        } `json:"goods_detail"`
    } `json:"consume_information"`
}

type CreateCouponRequest struct {
    StockID           string `json:"stock_id"`
    OutRequestNo      string `json:"out_request_no"`
    Appid             string `json:"appid"`
    StockCreatorMchid string `json:"stock_creator_mchid"`
    CouponValue       int    `json:"coupon_value"`
    CouponMinimum     int    `json:"coupon_minimum"`
}

type CreateCouponResponse struct {
    wxpay.ErrorMessage
    CouponID string `json:"coupon_id"`
    TraceNo  string `json:"trace_no"`
}

func (srv *CouponService) Create(ctx context.Context, openid string, req *CreateCouponRequest) (rsp *CreateCouponResponse, err error) {
    path := fmt.Sprintf("marketing/srv/users/%s/coupons", openid)
    request, err := srv.Client.NewRequest(http.MethodPost, path, req)
    if err != nil {
        return
    }
    rsp = &CreateCouponResponse{}
    if err = srv.Client.Do(request, rsp); err != nil {
        return
    }
    return
}

type GetCouponOptions struct {
    Appid string `url:"appid"`
}

func (srv *CouponService) Get(ctx context.Context, appid, couponId, openid string) (rsp *Coupon, err error) {
    path := fmt.Sprintf("marketing/srv/users/%s/coupons/%s", openid, couponId)
    rawurl, err := wxpay.AddOptions(path, &GetCouponOptions{Appid: appid})
    if err != nil {
        return
    }
    request, err := srv.Client.NewRequest(http.MethodGet, rawurl, nil)
    if err != nil {
        return
    }
    rsp = &Coupon{}
    if err = srv.Client.Do(request, rsp); err != nil {
        return
    }

    return
}
