# Chatty

Chatty is a web service for chat rooms built with Go.

## Why?

This project is made mainly for learning purposes. I wanted to practice websockets and goroutines in Go.
In order to do that, I decided to build a simple chat room service.

## Is it production ready?

No. This project is not production ready. Although it is completely functional it lacks of 
authentication, authorization and instrumentation. From scalability perspective, the rooms
live as long as the chat service is running. This means that if the service is restarted, all
rooms will be lost. This is not a problem for a small chat service, but it is not ideal.
Lastly, it is a good starting point for a chat service.

## Functionalities

- [x] Create a room
- [x] Join a room using websocket
- [x] Send a message to a room
- [x] Receive messages from a room
- [x] Graceful shutdown
- [x] Intermittent ping to keep connection alive

## Architecture

This project is built with primarily two layers: business and platform.
Business is responsible for application logic such as chat room creation and participant actions.
On the other hand, platform is responsible for persisting data, handling network connections and 
orchestrating message pipelines.