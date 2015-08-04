
package main

import (
  "fmt"
  "os"
  "net"

  swlb "github.com/mrcrilly/software-aware-balancing"
  "github.com/zenazn/goji"
)

func main() {
  err := swlb.Load(os.Args[1])

  if err != nil {
    panic(err)
  }

  fmt.Println("Starting API with the following configuration...")
  fmt.Printf("\tBindPort: %v\n\tBindIP: %v\n\tAPIName: %v\n", swlb.Cfg.BindIP, swlb.Cfg.BindPort, swlb.Cfg.APIName)

  bind_to := fmt.Sprintf("%v:%v", swlb.Cfg.BindIP, swlb.Cfg.BindPort)
  listener, err := net.Listen("tcp", bind_to)

  goji.Get("/", swlb.RouteIndex)
  goji.Get("/health", swlb.RouteHealth)
  goji.Get("/health/flip", swlb.RouteFlipHealth)
  
  goji.ServeListener(listener)
}
