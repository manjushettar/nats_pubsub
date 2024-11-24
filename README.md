# Building Secure Event-Driven Architecture (with Golang and NATS)

## What this is: 

A lightweight pub-sub service using NATS and event tokens for streaming validation. I am learning Golang so I thought it would be fun to implement this project such. This is a technical implementation of secure event streaming in distributed microservices. Of course, there are abstractions.

## Problem: 
Distributed event streaming poses problems that aren't solved by JWT authentication. Why?
1. Decoupling caused by the pub-sub model complicates validation because publishers and subscribers cannot directly authenticate each other. 
2. Asynchronous processing (what the pub-sub model enables) makes time bounded tokens, like JWT access tokens, unstable. A subscriber may not consume data if it is busy with work. Remember - the publisher has no idea about the subscribers who are listening to the event channel and thinks its work has been processed once it sends it to the messaging queue/notification system. In fact, a publisher does not even care what happens to data after it's done generating it. 
3. An event may go through several subscribers - and several subscribers can be publishers as well! We might have to use a lot of redundant architecture at different event endpoints that handle different logic, so we need a way to standardize it. 
4. Relating to the point above, different resources should be able to handle the data differently. We want permissions for our event streaming. We want public and private streams and we want ZERO TRUST!! Why? Because in large software, rogue modules and resources may be able to listen to event channels they are not permitted to!


## What repo provides:
This is simply an abstraction for event tokens. It is the product of a literature review, combining the best practices from the industry. This library provides:

1. A dual-token system
    - JWT access tokens for publisher authentication
    - Extended event timestamps for long-running processes
    - Data hashing to ensure data integrity during its lifecycle

2. Zero-Trust architecture
    - Templated auth at every service boundary implemented using a Token Manager
        - Source verification (publisher auth)
        - Consumer verification

# Understanding NATS and Event Security
NATS is a medium in Golang to emulate pub-sub. 
## The Problem: Secure Message Delivery
Let's say you have a friend. Let's also say you would like to send that friend a postcard. What do you do? You pack the postcard up into an envelope, write your name on it, and send it to the Post Office. After that, it's the Post Office's problem and you expect them to deliver the message. 

You don't know how many people will interact with your postcard until it reaches your destination, but you wish it was only your friend. You are realistic, however, and know that many people will likely interact with your postcard, BUT you do know that you want only employees of the Post Office (government employees) to touch your envelope. The postcard is strictly off limits. 

But you can't ensure this! Augh, that's annoying. You could have the post office send you a message to confirm the status of each employee who handled your envelope and maybe they could send a picture of it, showing that it was unopened and delivered. Yeah, but this gets annoying. So many messages, especially if you're far away. Think about all the return mail that the Post Office and each employee will have to send! And think about how much mail you will get just for confirmations. If anything, a confirmation at the end showing it was delivered is sufficient. 

## Enter NATS: A message broker 
Publishers are entities that can publish event data as topics. Subscribers are entities that can consume event data from those topics by subscribing to the topics. The topics and the event data are stored in a message distributor or a **message broker**. Think of NATS as a message broker that handles:

1. Message distribution
2. Event Streaming over topics
3. Enables connections to publisher and subscribers

## Another metaphor: scalability and decoupling
Let's say you are radio technician and you want to install a personal system at home for your family. Your mom would like to listen to her TV shows, your dad wishes that he could have the morning news playing, and your newborn brother wants to listen to Cocomelon (you should probably not allow him to). You don't want to have to set up a new device for every connection. 

It's much easier to abstract this. You'll just have one provider on your end. This provider will let you define different topics and send out the radio data (in signals) to some sort of station or hub. From there, if your mom knows she wants to listen to her TV shows, or if one day she wants to watch Cocomelon as well, she changes the channel of the hub to the specified channel. There, she receives the data. Neither you nor your mother interacted with each other directly. The hub is just a broker. 

This is the principle of the publisher-subscriber model. NATS serves to be the hub. NATS is an advanced hub that handles failed deliveries and reconnections. NATS enables decoupling, which means that publishers and subscribers do not need to interact with each other. NATS handles your routing, which means it is scalable. 

## Without NATS (or something like it)
Without NATS, you handle everything in between. Since we're focused on the publisher-subscriber aspect of NATS, one thing you would be concerned with is updating all subscribers when adding a new publisher. Your setup would have to start from 0 every time you were adding a new connection. 

# Implementation:

## Events

1. **auth/** -> this subdirectory contains the implementation of the token manager, handling both JWT access tokens and event tokens
2. **events/** -> this subdirectory contains the implementation of the publisher and subscriber structures
3. **main.go** -> the set up file that defines a publisher and two subscribers, using a NATS server to distribute messages

Now that we understand what events are, this is how they are implemented in this repo.

```
type Event struct {
    ID        string          `json:"id"`
    Data      json.RawMessage `json:"data"`
    JWT       string          `json:"jwt"`
    EventAuth auth.EventToken `json:"event_auth"`
}
```

Every event has an ID and some data. Also associated with an event is a JWT access token to provide a base level of publisher verification and an event authorization token to ensure that the event ends up at an authorized subscriber. Events are defined at the publisher/subscriber level in order to abstract it away from the Token Manager, which doesn't need to know about the implementation of different events.

The main Publisher functions can be found in **/events/pubsub.go**.
```
type Publisher struct {
    ID           string
    natsConn     *nats.Conn
    tokenManager *auth.TokenManager
}
```

A publisher is identified by an ID. A publisher also has a NATS connection that lets it connect to the NATS server, along with a Token Manager to authenticate events. In the rest of the file, a constructor is defined as well as a Publish() function that allows publishers to publish event data to the NATS server.

Similarily, a Subscriber structure is implemented in **/events/pubsub.go**.
```
type Subscriber struct {
    ID           string
    natsConn     *nats.Conn
    tokenManager *auth.TokenManager
}
```

Subscribers behave very similary to publishers, sharing a NATS connection. Subscribers have a Subscribe() function that allows them to read event data being streamed from the NATS connection.

The **token manager** is one of the key implementations in this module. It is defined as such:

```
type TokenManager struct {
    jwtSecret   []byte
    eventSecret []byte
}
```

All the token manager has to do is verify the JWT access token and event token secrets that are used by the publishersiand subscribers. A TokenManager has functions to create both types of secrets, as well as handle the verification of events that are passed through it. Publishers and Subscribers must use the same TokenManager instance in order to verify signatures. Publishers and subscribers using TokenManagers with different hashing methods or timestamps cannot interact with each other and will fail at the TokenManager level. 

**This means that a subscriber MUST know the event secret AND the JWT access token to read any data from the message broker.** A subscriber can still connect to the NATS server, but it cannot read any data from publishers that are using a different pair of secrets. This enforces two-way verfication without the individual publisher/subscriber instances being directly aware of the other. 

This is also where the benefits of decoupling become kind of blurred. The subscriber still can listen to the channel, but fails to read any data. Invalidating the connection would be a NATS-level modifcation and is out of the scope of this implementation.

## IMPORTANT:
In **main.go**, one subscriber uses the same event and token signature as the publisher, so it has access to the event data being written. The other subscriber (subscriber 2) uses a different pair of tokens.

## This implementation can be continued by:
Adding more publishers and subscribers that *do something*. Currently, publishers just create hello messages and subscribers simply consume those messages and say they received them. In real world applications, these are resources and instances that do something and spit out more things. Chaining these publishers and subscribers is a common functionality and can be implemented as such by just changing the interface functions.






