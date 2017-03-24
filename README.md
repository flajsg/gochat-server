# gochat-server
A simple TCP chat server/client in GoLang.

## Introduction

This is just an example how easy it is to use built in GO packages, to create 
a TCP connections and send requests between client and a server. I've created
this examples as a learning experiment for myself.

I'm publishing it to the public for anyone to use/alter/contribute. You can use this as a fundation to extend it beyond
what it was initially intended to be. The code is probably not perfect but it works as it is :)

## Basic usage

First run `go run server.go` in your first terminal to create a server, then run `go run client.go` in another terminal
to connect to server. Server will immediately send you a welcome message:

    << Welcome to gochat-server 1.0
    
To login simply type:

    NICK mynickname
    
Server will send you a welcome message:

    << Welcome mynickname
    
You can now start sending messages. Every message is automatically sent when you hit enter and it is distributed to all connecter clients. 
    
To exit the client type:

    exit
        
That's it!

You can run another client in a third terminal to chat between them.

## TODOs

Since this is just a simple showcase on what you can do with Go, I didn't implement any other common chat commands.
This is something that you can do it yourself. For example:

- sending private chat,
- creating and entering chat rooms,
- coommands to display logged clients,
- ...
