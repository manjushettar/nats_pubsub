# CMPE 272 Final Project - Building Secure-Driven Architecture (with Golang and NATS)

## What this: 

A lightweight pub-sub service using NATS and event tokens for streaming validation. I am learning Golang so I thought it would be implement this project such. This is a technical implementation of secure event streaming is distributed microservices. Of course, there are abstractions.

## Problem: 
Distributed event streaming poses problems that aren't solved by JWT authentication. Why?
1. Decoupling caused by the pub-sub model complicates enforcing validation because publishers and subscribers cannot directly interact to authenticate each other. 
2. Asynchronous processing (what the pub-sub model enables) makes time bounded tokens, like JWT access tokens, unstable. A subscriber may not consume data if it is busy with work. Remember - the publisher has no idea about the subscribers who are listening to the event channel and thinks its work has been processed once it sends it to the messaging queue/notification system. 
3. An event may go through several subscribers - and several subscribers can be publishers as well! We might have to use a lot of redundant architecture at different event endpoints that handle different logic, so we need a way to standardize it - that would be very nice. 
4. Relating to the point above, different resources should be able to handle the data differently. We want permissions for our event streaming. We want public and private streams and we want ZERO TRUST!! Why? Because in large software, rogue modules and resources may be able to listen to event channels they are not permitted to!


## What this does:
This is simply an abstraction for event tokens. It is the product of a literature review, combining the best practices from the industry. This library provides:

1. A dual-token system
    - JWT for publisher authentication
    - Extended event timestamps for long-running processes
    - Data hashing to ensure data integrity during its lifecycle

2. Zero-Trust architecture
    - Templated auth at every service boundary
    - Source verification (publisher auth)
    - Decoding the data hash

# 



# Understanding NATS and Event Security
NATS is essentially a medium in Golang to emulate pub-sub. 
## The Problem: Secure Message Delivery
Let's say you have a friend. Let's also say you would like to send that friend a postcard. What do you do? You pack the postcard up into an envelope, write your name on it, and send it to the Post Office. After that, it's the Post Office's problem and you expect them to deliver the message. 

You don't know how many people will interact with your postcard until it reaches your destination, but you wish it was only your friend. Regardless, you know that you definitely want only employees of the Post Office (government employees) to touch your envelope. The postcard is strictly off limits. 

But you can't ensure this! Augh, that's annoying. You could have the post office send you a message to confirm the status of each employee who handled your envelope and maybe they could send a picture of it, showing that it was unopened and delivered. Yeah, but this gets annoying. So many messages, especially if your far away. Think about all the return mail that the Post Office and each employee will have to send! And think about how much mail you will get just for confirmations. If anything, a confirmation at the end showing it was delivered is sufficient. 

## Enter NATS: A message broker 
Publishers are entities that can publish event data as topics. Subscribers are entities that can consume event data from those topics by subscribing to the topics. The topics and the event data are stored in a message distributor or a **message broker**. NATS is a message broker that handles:

1. Message distribution
2. Event Streaming over topics
3. Enables connections to publisher and subscribers

## Another metaphor: scalability and decoupling
Let's say you are radio technician and you want to install a personal system at home for your family. Your mom would like to listen to her TV shows, your dad wishes that he could have the morning news playing, and your newborn brother wants to listen to Cocomelon (you should probably not allow him to -> s/o twitter!). You don't want to have to set up a new device for every connection. 

It's much easier to abstract this. You'll just have one provider on your end. This provider will let you define different topics and send out the radio data (in signals) to some sort of station or hub. From there, if your mom knows she wants to listen to her TV shows, or if one day she wants to watch Cocomelon as well, she changes the channel of the hub to the specified channel. There, she receives the data. Neither you nor your mother interacted with each other directly. The hub is just a broker. 

This is the principle of the publisher-subscriber model. NATS serves to be the hub. NATS is an advanced hub that handles failed deliveries and reconnections. NATS enables decoupling, which means that publishers and subscribers do not need to interact with each other. NATS handles your routing, which means it is scalable. 

## Without NATS (or something like it)
Without NATS, you handle everything in between. Since we're focused on the publisher-subscriber aspect of NATS, one thing you would be concerned with is updating all subscribers when adding a new publisher. 







