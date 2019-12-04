package favor

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/yunlyz/wxpay"
)

type StockService wxpay.Service

// Stock represents a Wechat merchant stock
type Stock struct {
    StockId            string    `json:"stock_id"`
    StockName          string    `json:"stock_name"`
    Comment            string    `json:"comment"`
    BelongMerchant     string    `json:"belong_merchant"`
    AvailableBeginTime time.Time `json:"available_begin_time"`
    AvailableEndTime   time.Time `json:"available_end_time"`
    DistributedCoupons int64     `json:"distributed_coupons"`
    StockUseRule       struct {
        MaxCoupons         int  `json:"max_coupons"`
        MaxAmount          int  `json:"max_amount"`
        MaxAmountByDay     int  `json:"max_amount_by_day"`
        MaxCouponsPerUser  int  `json:"max_coupons_per_user"`
        NaturalPersonLimit bool `json:"natural_person_limit"`
        PreventAPIAbuse    bool `json:"prevent_api_abuse"`
    } `json:" stock_use_rule"`
    PatternInfo struct {
        Description     string `json:"description"`
        MerchantLogo    string `json:"merchant_logo"`
        MerchantName    string `json:"merchant_name"`
        BackgroundColor string `json:" background_color"`
        CouponImage     string `json:" coupon_image "`
    } `json:"pattern_info"`
    CouponUseRule struct {
        CouponAvailableTime struct {
            FixAvailableTime struct {
                AvailableWeekDay []int `json:" available_week_day "`
                BeginTime        int   `json:" begin_time "`
                EndTime          int   `json:" end_time "`
            } `json:"fix_available_time"`
            SecondDayAvailable        bool `json:"second_day_available"`
            AvailableTimeAfterReceive int  `json:"available_time_after_receive"`
        } `json:"coupon_available_time"`
        FixedNormalCoupon struct {
            CouponAmount       int `json:" coupon_amount"`
            TransactionMinimum int `json:"transaction_minimum"`
        } `json:" fixed_normal_coupon "`
        DisscountCoupon struct {
            DiscountAmountMax  int `json:"discount_amount_max"`
            DiscountPercent    int `json:"discount_percent"`
            TransactionMinimum int `json:"transaction_minimum"`
        } `json:"disscount_coupon"`
        ExchangeCoupon struct {
            SinglePriceMax int `json:"single_price_max"`
            ExchangePrice  int `json:"exchange_price"`
        } `json:"exchange_coupon"`
        GoodsTag           []string `json:"goods_tag"`
        TradeType          string   `json:"trade_type"`
        CombineUse         bool     `json:"combine_use"`
        AvailableItems     []string `json:"available_items"`
        UnavailableItems   []string `json:"unavailable_items"`
        AvailableMerchants []string `json:"available_merchants"`
    } `json:"coupon_use_rule"`
    NoCash       bool      `json:"no_cash"`
    StartTime    time.Time `json:"start_time"`
    StopTime     time.Time `json:"stop_time"`
    CutToMessage struct {
        SinglePriceMax int64 `json:"single_price_max"`
        CutToPrice     int64 `json:"cut_to_price"`
    } `json:"cut_to_message"`
    Singleitem   bool   `json:"singleitem"`
    StockType    string `json:"stock_type"`
    OutRequestNo string `json:"out_request_no"`
}

type CreatorMchOptions struct {
    StockCreatorMchid string `url:"stock_creator_mchid"`
}

type CreateStockResponse struct {
    StockID    string    `json:"stock_id"`
    CreateTime time.Time `json:"create_time"`
    wxpay.ErrorMessage
}

// CreateStock-创建代金券批次
// https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_1.shtml
// 通过此接口可创建代金券批次，包括预充值&免充值类型
func (srv *StockService) CreateStock(ctx context.Context, stock *Stock) (result *CreateStockResponse, err error) {
    request, err := srv.Client.NewRequest(http.MethodPost, "marketing/srv/coupon-stocks", stock)
    if err != nil {
        return
    }
    result = &CreateStockResponse{}
    if err = srv.Client.Do(request, result); err != nil {
        return
    }
    return
}

type ActivateStockResponse struct {
    StartTime string `json:"start_time"`
    StockID   string `json:"stock_id"`
}

// ActivateStock 激活代金券批次
// https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_3.shtml
// 制券成功后，可调用此接口激活代金券批次
func (srv *StockService) ActivateStock(ctx context.Context, stockCreatorMchId, stockID string) (
result *ActivateStockResponse, err error) {
    opt := &CreatorMchOptions{StockCreatorMchid: stockCreatorMchId}
    path := fmt.Sprintf("marketing/srv/stocks/%s/start", stockID)
    rawurl, err := wxpay.AddOptions(path, opt)
    if err != nil {
        return
    }

    req, err := srv.Client.NewRequest(http.MethodGet, rawurl, nil)
    if err != nil {
        return
    }

    result = &ActivateStockResponse{}
    err = srv.Client.Do(req, result)
    if err != nil {
        return
    }

    return
}

type PauseStockResponse struct {
    PauseTime string `json:"pause_time"`
    StockID   string `json:"stock_id"`
}

// PauseStock 暂停代金券批次
// https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_13.shtml
// 通过此接口可暂停指定代金券批次。暂停后，该代金券批次暂停发放。
func (srv *StockService) PauseStock(stockCreatorMchId, stockID string) (result *PauseStockResponse, err error) {
    opt := &CreatorMchOptions{StockCreatorMchid: stockCreatorMchId}
    path := fmt.Sprintf("marketing/srv/stocks/%s/pause", stockID)
    rawurl, err := wxpay.AddOptions(path, opt)
    if err != nil {
        return
    }

    req, err := srv.Client.NewRequest(http.MethodGet, rawurl, nil)
    if err != nil {
        return
    }

    result = &PauseStockResponse{}
    err = srv.Client.Do(req, result)
    if err != nil {
        return
    }

    return
}

// RestartStock 重启代金券批次
// https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_14.shtml
// 通过此接口可重启指定代金券批次。重启后，该代金券批次可以再次发放。
func (srv *StockService) RestartStock(stockCreatorMchId, stockID string) (result *PauseStockResponse, err error) {
    opt := &CreatorMchOptions{StockCreatorMchid: stockCreatorMchId}
    path := fmt.Sprintf("marketing/srv/stocks/%s/pause", stockID)
    rawurl, err := wxpay.AddOptions(path, opt)
    if err != nil {
        return
    }

    req, err := srv.Client.NewRequest(http.MethodGet, rawurl, nil)
    if err != nil {
        return
    }

    result = &PauseStockResponse{}
    err = srv.Client.Do(req, result)
    if err != nil {
        return
    }

    return
}

type QueryStocksOptions struct {
    Offset            uint32 `url:"offset"`
    Limit             uint32 `url:"limit"`
    StockCreatorMchid string `url:"stock_creator_mchid"`
    CreateStartTime   string `url:"create_start_time"`
    CreateEndTime     string `url:"create_end_time"`
    Status            string `url:"status"`
}

type QueryStocksResponse struct {
    TotalCount int      `json:"total_count"`
    Data       []*Stock `json:"data"`
    Limit      int      `json:"limit"`
    Offset     int      `json:"offset"`
}

func (srv *StockService) QueryStocks(opts *QueryStocksOptions) (result *QueryStocksResponse, err error) {
    req, err := srv.Client.NewRequest(http.MethodGet, "marketing/srv/stocks", nil)
    if err != nil {
        return
    }

    result = &QueryStocksResponse{}
    err = srv.Client.Do(req, result)
    if err != nil {
        return
    }

    return
}

func (srv *StockService) GetStock(stockCreatorMchId, stockID string) (result *Stock, err error) {
    opt := &CreatorMchOptions{StockCreatorMchid: stockCreatorMchId}
    path := fmt.Sprintf("marketing/srv/stocks/%s", stockID)
    rawurl, err := wxpay.AddOptions(path, opt)
    if err != nil {
        return
    }

    req, err := srv.Client.NewRequest(http.MethodGet, rawurl, nil)
    if err != nil {
        return
    }

    result = &Stock{}
    err = srv.Client.Do(req, result)
    if err != nil {
        return
    }

    return
}
