
package main

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "encoding/json"
  "net/http"
)

type APIResponse struct {
  Result bool
  Message string
}

type ConfigurationEndpoint struct {
  RemoteIP string
  RemotePort string
  RemoteName string
}

type Configuration struct {
  Endpoints []ConfigurationEndpoint
}

var json_config Configuration

func get_api_name(endpoint *ConfigurationEndpoint) (response *APIResponse, err error) {
  resp, err := http.Get(fmt.Sprintf("http://%v:%v/", endpoint.RemoteIP, endpoint.RemotePort))
  defer resp.Body.Close()

  if err != nil {
    return nil, err
  }

  if resp.StatusCode == 200 {
    var resp_json APIResponse

    raw, err := ioutil.ReadAll(resp.Body)

    if err != nil {
      return nil, err
    }

    err = json.Unmarshal(raw, &resp_json)

    if err != nil {
      return nil, err
    }

    return &resp_json, nil
  } else {
    return nil, errors.New("Unable to reach endpoint")
  }
}

func first_healthy_endpoint() (use_endpoint *ConfigurationEndpoint, err error) {
  use_endpoint = nil

  if len(json_config.Endpoints) > 0 {
    for _, ep := range json_config.Endpoints {
      resp, err := http.Get(fmt.Sprintf("http://%v:%v/health", ep.RemoteIP, ep.RemotePort))

      if err != nil {
        // If we get an error, we assume the endpoint is down
        // This COULD be more intelligent, but why care about the API being
        // down, and instead just move onto a working one - sysad will
        // have the broken API alerting in monitoring, right?
        continue
      }

      if resp.StatusCode == 200 {
        // fmt.Println("Got 200")
        var resp_json APIResponse

        raw, err := ioutil.ReadAll(resp.Body)

        if err != nil {
          return nil, err
        }

        // fmt.Println("Got JSON")
        err = json.Unmarshal(raw, &resp_json)

        if err != nil {
          return nil, err
        }


        if resp_json.Result {
          // fmt.Println("Endpoint Healthy")
          use_endpoint = &ep
          break
        } else {
          // fmt.Println("Endpoint Unhealthy")
        }
      }

      resp.Body.Close()
    }

    if use_endpoint == nil {
      // fmt.Println("Returning NIL")
      return nil, errors.New("No working endpoints")
    } else {
      return use_endpoint, nil
    }
  } else {
    return nil, errors.New("No endpoints configured.")
  }
}

func main() {
  raw_config, err := ioutil.ReadFile(os.Args[1])

  if err != nil {
    panic(err)
  }

  err = json.Unmarshal(raw_config, &json_config)

  if err != nil {
    panic(err)
  }

  fmt.Println("Configuration loaded.")

  for _, ep := range json_config.Endpoints {
    fmt.Printf("Configured Endpoint: %v, Port: %v, Name: %v\n", ep.RemoteIP, ep.RemotePort, ep.RemoteName)
  }

  healthy_endpoint, err := first_healthy_endpoint()

  if healthy_endpoint == nil {
    fmt.Println("No endpoints to work with. Exiting.")
    os.Exit(1)
  }

  fmt.Printf("Using endpoint: %v\n", healthy_endpoint.RemoteName)

  response, err := get_api_name(healthy_endpoint)

  if err != nil {
    panic(err)
  }

  fmt.Printf("Response Message (Result): %v (%v)\n", response.Message, response.Result)
}
