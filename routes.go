
package loadbalancer

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/zenazn/goji/web"
)

type TypeRespsonse struct {
  Result bool
  Message string
}

func RouteIndex(c web.C, w http.ResponseWriter, r *http.Request) {
  return_data := TypeRespsonse{Result: true, Message: fmt.Sprintf("This is API: %v", Cfg.APIName)}
  return_data_json, err := json.Marshal(return_data)

  if err != nil {
    panic(err)
  }

  fmt.Fprint(w, string(return_data_json))
}

func RouteFlipHealth(c web.C, w http.ResponseWriter, r *http.Request) {
  Cfg.Health = ! Cfg.Health

  return_data := TypeRespsonse{Result: Cfg.Health, Message: "OK"}
  return_data_json, err := json.Marshal(return_data)

  if err != nil {
    panic(err)
  }

  fmt.Fprint(w, string(return_data_json))
}

func RouteHealth(c web.C, w http.ResponseWriter, r *http.Request) {
  return_data := TypeRespsonse{Result: Cfg.Health, Message: "OK"}
  return_data_json, err := json.Marshal(return_data)

  if err != nil {
    panic(err)
  }

  fmt.Fprint(w, string(return_data_json))
}
