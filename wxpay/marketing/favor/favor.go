package favor

import (
    "context"

    "github.com/yunlyz/go-wechat/wxpay/client"
)

// Favor API V3文档-营销分类
type Favor struct {
    Stock    *StockService
    Coupon   *CouponService
    Callback *CallbackService
}

func New(ctx context.Context, srv *client.Service) *Favor {
    favor := &Favor{}
    favor.Stock = (*StockService)(srv)
    favor.Coupon = (*CouponService)(srv)
    favor.Callback = (*CallbackService)(srv)
    
    return favor
}
