package favor

import (
    "context"

    "github.com/yunlyz/go-wechat/wxpay"
)

// Favor API V3文档-营销分类
type Favor struct {
    Stock    *StockService
    Coupon   *CouponService
    Callback *CallbackService
}

func New(ctx context.Context, client *wxpay.Client) *Favor {
    return &Favor{
        Stock:    (*StockService)(client.Common),
        Coupon:   (*CouponService)(client.Common),
        Callback: (*CallbackService)(client.Common),
    }
}
