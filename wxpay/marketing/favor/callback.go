package favor

import (
    "context"
    "net/http"
    "time"

    "github.com/yunlyz/go-wechat/wxpay/client"
)

type CallbackService client.Service

type SetCallbackRequest struct {
    Mchid     string `json:"mchid"`
    NotifyURL string `json:"notify_url"`
    Switch    bool   `json:"switch,omitempty"`
}

type SetCallbackResponse struct {
    UpdateTime time.Time `json:"update_time"`
    NotifyURL  string    `json:"notify_url"`
    Mchid      string    `json:"mchid"`
}

func (srv *CallbackService) SetCallback(ctx context.Context, req *SetCallbackRequest) (rsp *SetCallbackResponse, err error) {
    request, err := srv.Client.NewRequest(http.MethodPost, "marketing/favor/callbacks", req)
    if err != nil {
        return
    }
    rsp = &SetCallbackResponse{}
    if err = srv.Client.Do(request, rsp); err != nil {
        return
    }

    return
}
