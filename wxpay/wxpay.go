package wxpay

import (
    "bytes"
    "crypto"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/json"
    "encoding/pem"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "reflect"
    "strings"
    "sync"
    "time"

    "github.com/google/go-querystring/query"
    "github.com/thanhpk/randstr"
    "github.com/yunlyz/wxpay/marketing/favor"
)

const (
    defaultBaseURL = "https://api.mch.weixin.qq.com/v3/"
    userAgent      = "go-wxpay"

    defaultMediaType = "application/json"

    headerTimestamp = "Wechatpay-Timestamp"
    headerNonce     = "Wechatpay-Nonce"
    headerSignature = "Wechatpay-Signature"

    defaultAuthType = "WECHATPAY2-SHA256-RSA2048"

    fmtSign = "%s\n%s\n%d\n%s\n%s\n"
    fmtAuth = "%s mchid=\"%d\",nonce_str=\"%s\",signature=\"%s\",timestamp=\"%d\",serial_no=\"%s\""
)

type Service struct {
    Client *WxPay
}

type WxPay struct {
    mu     sync.Mutex
    client *http.Client
    common *Service

    BaseURL *url.URL

    MchId      int64
    SerialNo   string
    PrivateKey string
    PublicKey  string
}

func New(mchId int64, serialNo, privateKey, publicKey string) *WxPay {
    baseURL, _ := url.Parse(defaultBaseURL)

    pay := &WxPay{
        mu:         sync.Mutex{},
        client:     &http.Client{},
        BaseURL:    baseURL,
        MchId:      mchId,
        SerialNo:   serialNo,
        PrivateKey: privateKey,
        PublicKey:  publicKey,
    }
    pay.common.Client = pay

    return pay
}

func (pay *WxPay) Favor() *favor.Favor {
    return &favor.Favor{
        Stock:    (*favor.StockService)(pay.common),
        Coupon:   (*favor.CouponService)(pay.common),
        Callback: (*favor.CallbackService)(pay.common),
    }
}

func (pay *WxPay) Busifavor() {

}

func (pay *WxPay) MiniProgram() {

}

func (pay *WxPay) PayScore() {

}

func (pay *WxPay) Other() {

}

func (pay *WxPay) Common() {

}

func (pay *WxPay) NewRequest(method, rawurl string, body interface{}) (req *http.Request, err error) {
    if !strings.HasSuffix(pay.BaseURL.Path, "/") {
        return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", pay.BaseURL)
    }
    u, err := pay.BaseURL.Parse(rawurl)
    if err != nil {
        return
    }

    var buf io.ReadWriter
    var text []byte
    if body != nil {
        text, err = json.Marshal(body)
        buf = bytes.NewBuffer(text)
        if err != nil {
            return
        }
    }

    req, err = http.NewRequest(method, u.String(), buf)
    if err != nil {
        return
    }

    timestamp := time.Now().Unix()
    nonce := randstr.Hex(32)
    signature, err := pay.sign(req, timestamp, nonce, string(text))
    if err != nil {
        return
    }

    req.Header.Set("Content-Type", defaultMediaType)
    req.Header.Set("User-Agent", userAgent)
    req.Header.Set("Accept", defaultMediaType)
    req.Header.Set("Authorization", pay.authorization(timestamp, nonce, signature))

    return
}

func (pay *WxPay) Do(req *http.Request, v interface{}) (err error) {
    response, err := pay.client.Do(req)
    if err != nil {
        return
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }

    // @TODO
    _ = checkResponse(response)
    // verify signature
    if err = pay.verify(response.Header, body); err != nil {
        return
    }

    if v != nil {
        buf := bytes.NewReader(body)
        if w, ok := v.(io.Writer); ok {
            _, _ = io.Copy(w, buf)
        } else {
            err = json.NewDecoder(buf).Decode(v)
            if err == io.EOF {
                err = nil // ignore EOF errors caused by empty response body
            }
            if err != nil {
                return
            }
        }
    }

    return
}

func (pay *WxPay) sign(req *http.Request, timestamp int64, nonce, body string) (signature string, err error) {
    str := fmt.Sprintf(fmtSign, req.Method, req.URL.Path, timestamp, nonce, body)

    hash := sha256.New()
    hash.Write([]byte(str))
    block, _ := pem.Decode([]byte(pay.PrivateKey))
    key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return
    }
    sign, err := rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hash.Sum(nil))
    if err != nil {
        return
    }
    signature = base64.StdEncoding.EncodeToString(sign)

    return
}

func (pay *WxPay) verify(header http.Header, body []byte) (err error) {
    timestamp := header.Get(headerTimestamp)
    nonce := header.Get(headerNonce)
    signature := header.Get(headerSignature)

    signStr := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, string(body))
    decodeString, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        return
    }

    hash := sha256.New()
    hash.Write([]byte(signStr))

    block, _ := pem.Decode([]byte(pay.PublicKey))
    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return
    }

    return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), crypto.SHA256, hash.Sum(nil), decodeString)
}

func (pay *WxPay) authorization(timestamp int64, nonce, signature string) string {
    return fmt.Sprintf(fmtAuth, defaultAuthType, pay.MchId, nonce, signature, timestamp, pay.SerialNo)
}

func checkResponse(resp *http.Response) error {
    return nil
}

type ErrorMessage struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func AddOptions(s string, opt interface{}) (string, error) {
    v := reflect.ValueOf(opt)
    if v.Kind() == reflect.Ptr && v.IsNil() {
        return s, nil
    }

    u, err := url.Parse(s)
    if err != nil {
        return s, err
    }

    qs, err := query.Values(opt)
    if err != nil {
        return s, err
    }

    u.RawQuery = qs.Encode()
    return u.String(), nil
}

// CertificateDecrypt using AES-256-GCM algorithm to decrypt ciphertext
// Notice: ciphertext is decrypted after base64
func CertificateDecrypt(ciphertext []byte, key, nonce, associated string) (str string, err error) {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return
    }
    plaintext, err := gcm.Open(nil, []byte(nonce), ciphertext, []byte(associated))
    if err != nil {
        return
    }
    str = string(plaintext)
    return
}
