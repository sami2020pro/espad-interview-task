package main

import (
  "log"
  "task/config" // configuration.go | at task/config
  "task/handler" // handler.go | at task/handler
  "task/storage" // storage.go | at task/storage 

  "github.com/valyala/fasthttp" // library
)

// main function for run the codes
func main() {
  // Read the configuration file that's mean task/configuration.json
  configuration, err := config.FromFile("./configuration.json")
  if err != nil { // We check that if *err* has an error, we can log that error 
     log.Fatal(err)
  }

  // Create new service for connect to Redis
  service, err := storage.New(configuration.Redis.Host, configuration.Redis.Port, configuration.Redis.Password)
  if err != nil { // We check that if *err* has an error, we can log that error 
     log.Fatal(err)
  }
  defer service.Close() // defer is available in golang

  // Create the new router from handler and New function
  router := handler.New(configuration.Options.Schema, configuration.Options.Prefix, service)

  log.Fatal(fasthttp.ListenAndServe(":" + configuration.Server.Port, router.Handler))
}

/* ('Sami Ghasemi) */
