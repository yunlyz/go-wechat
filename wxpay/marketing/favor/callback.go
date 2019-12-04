package favor

import (
    "context"
    "net/http"
    "time"

    "github.com/yunlyz/wxpay"
)

type CallbackService wxpay.Service

type SetCallbackRequest struct {
    Mchid     string `json:"mchid"`
    NotifyURL string `json:"notify_url"`
    Switch    bool   `json:"switch"`
}

type SetCallbackResponse struct {
    UpdateTime time.Time `json:"update_time"`
    NotifyURL  string    `json:"notify_url"`
    Mchid      string    `json:"mchid"`
}

func (srv *CallbackService) SetCallback(ctx context.Context, req *SetCallbackRequest) (rsp *SetCallbackResponse, err error) {
    request, err := srv.Client.NewRequest(http.MethodPost, "marketing/srv/callbacks", req)
    if err != nil {
        return
    }
    rsp = &SetCallbackResponse{}
    if err = srv.Client.Do(request, rsp); err != nil {
        return
    }

    return
}
