
package loadbalancer

import (
  "io/ioutil"
  "encoding/json"
)

type Configuration struct {
  BindIP string
  BindPort string
  APIName string
  Health bool
}

var Cfg Configuration

func Load(file_in string) (err error){
  var configuration []byte

  configuration, err = ioutil.ReadFile(file_in)

  if err != nil {
      return err
  }

  err = json.Unmarshal(configuration, &Cfg)

  if err != nil {
    return err
  }

  return nil
}
