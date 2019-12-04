package favor

// Favor API V3文档-营销分类
type Favor struct {
    Stock    *StockService
    Coupon   *CouponService
    Callback *CallbackService
}
