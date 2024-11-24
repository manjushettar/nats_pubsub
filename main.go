package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/nats-io/nats.go"
    "implementation/auth"
    "implementation/events"
)

func main() {
    // Configuration
    // jwtSecret and eventSecret are the predefined secrets (think passports) that we are defining
    jwtSecret := []byte("random-asdfghjkl")
    eventSecret := []byte("event-random-asdfghjkl")
    natsURL := "nats://localhost:4222"

    // connect to the nats connection (broker)
    nc, err := nats.Connect(natsURL)
    if err != nil {
        log.Fatalf("Failed to connect to NATS: %v", err)
    }
    defer nc.Close()

    // create a tokenmanager instance with our base jwt-secrets and eventSecret
    // pass this tokenmanager to any publisher and subscriber and it will authorize communication between the two
    tokenManager := auth.NewTokenManager(jwtSecret, eventSecret)
   

    jwtSecret2 := []byte("11111")
    eventSecret2 := []byte("111111")
    tokenManager2 := auth.NewTokenManager(jwtSecret2 , eventSecret2)
    
    publisher := events.NewPublisher("publisher-1", nc, tokenManager)
    subscriber := events.NewSubscriber("subscriber-1", nc, tokenManager)
    
    subscriber2 := events.NewSubscriber("subscriber-2", nc, tokenManager2)

    // create a subscription
    // whenever an event is received, the callback function gets run
    // this callback func prints the event data
    err = subscriber.Subscribe("events", func(data interface{}) error {
        fmt.Printf("Subscriber 1 received event: %+v\n", data)
        return nil
    })

    if err != nil {
        log.Fatalf("Failed to sub: %v", err)
    }

    err = subscriber2.Subscribe("chats", func(data interface{}) error {
        fmt.Printf("Subscriber 2 received event: %+v\n", data)
        return nil
    })

    if err != nil {
        log.Fatalf("Failed to sub: %v", err)
    }


    // this is a go routine, executing an example anonymous func
    // do below indefinitely:
        // wait every 5 seconds for a ticker value
        // wait every 3 seconds for a ticker2 value
        // for each ticker, create an event with a message that says hello and the timestamp
        // once that ticker is ready (created at the 3s and 5s time values), publish the event to the specified topic using the publisher

    go func() {
        ticker := time.NewTicker(5 * time.Second)
        ticker2 := time.NewTicker(3 * time.Second)
        for {
            select {
            case <-ticker.C:
                event := map[string]string{
                    "message": fmt.Sprintf("Hello, i am at %v", time.Now()),
                }
                if err := publisher.Publish("events", "greeting", event); err != nil {
                    log.Printf("Failed to pub: %v", err)
                }
            case <-ticker2.C:
                event := map[string]string{
                    "message": fmt.Sprintf("Hello, i am at %v", time.Now()),
                }
                if err := publisher.Publish("chats", "greeting", event); err != nil {
                    log.Printf("Failed to pub: %v", err)
                }
            }
        }
    }()
    
    //terminates nats and whole program with ctrl c
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    fmt.Println("\nShutting down...")
}
