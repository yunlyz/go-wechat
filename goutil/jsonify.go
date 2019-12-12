package goutil

import (
    "encoding/json"
)

func Jsonify(obj interface{}) string {
    bytes, _ := json.Marshal(obj)
    return string(bytes)
}
