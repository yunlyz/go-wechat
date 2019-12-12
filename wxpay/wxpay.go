package wxpay

import (
    "context"

    "github.com/yunlyz/go-wechat/wxpay/client"
    "github.com/yunlyz/go-wechat/wxpay/marketing/favor"
)

type wxpay struct {
    common *client.Service
    fav    *favor.Favor
}

func (pay *wxpay) GetFavor() *favor.Favor {
    return pay.fav
}

func New(mchId int64, serialNo, privateKey, publicKey string) *wxpay {
    common := &client.Service{
        Client: client.New(mchId, serialNo, privateKey, publicKey),
    }
    pay := &wxpay{}
    pay.common = common
    pay.fav = favor.New(context.Background(), pay.common)

    return pay
}
