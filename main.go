// Author: Gregor Flajs <flajsg@gmail.com> 2017. 

/*
This is simple chat server that listens on port 8080. Any new connection is 
handled by its own go routine.

Every client that connects must enter their nickname in order to chat.

To login client must sent "NICK [nickname]" command. Then they can chat.

Every chat message is send to all connected clients. 

TODOs (something you can try for yourself):
- private chat,
- chat rooms,
- command for a list of connected clients,
- ..... 
*/
package main

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "strings"
    "errors"
)

// Logged in users and their connections
type user struct {
    nick string
    conn net.Conn
}

// All logged users (clients)
var nicknames []user

// Listening port
var port = "8080"

// Establish TCP server and listen for incomming connections
func main() {
    
    log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate) 
    
    fmt.Println("Welcome to GoChat Server 1.0");
    
    ln, err := net.Listen("tcp", ":" + port)
    
    if err != nil {
        log.Fatalf("Failed to listen: %s", err)
    }
    
    fmt.Println("Listening on port:", port);

    for {
        if conn, err := ln.Accept(); err == nil {
            go handleConnection(conn)
        }
    }
    
    if err := ln.Close(); err == nil {
        log.Fatalf("Error closing connectio: %s", err)
    }
}

// handleConnection is spawned once per connection from a client, 
// and ends when the client is done sending requests.
func handleConnection(conn net.Conn) {
    defer conn.Close();
    
    // Logged in user
    var client = user{"", nil}

    sendRes(conn, "Welcome to gochat-server 1.0", client)
    
    for {
        
        // will listen for message to process ending in newline (\n)
        message, err := bufio.NewReader(conn).ReadString('\n')
        
        if err != nil {
            log.Println("Connection closed : " + client.nick)
            cleanClient(client)
            return
        }
        
        // remove trailing newline
        message = strings.Trim(string(message), "\n")

        // output message received
        log.Println("Message Received:", message)
        
        if len(message) == 0 {
            continue
        }
        
        // NICK operation -- login
        if strings.HasPrefix(message, "NICK") {
            new_client, err := login(conn, strings.Split(message, " ")[1])
            
            if err != nil {
                sendRes(conn, err.Error(), client)
            } else {
                client = new_client;
            }
            
        // exit -- close connection on client side
        } else if strings.Compare(message, "exit") == 0 {
            logout(conn, client)
        } else {
            
            if len(client.nick) == 0 {
                sendRes(conn, "You must first login by: NICK [nickname]", client)
            } else {
                sendToAllButMe(client.nick + ": "+ message, client)
            }
        }
    }
}

// Login checks if nickname is valid and adds it to global nicknames.
// It return client if all is good and error if login failed.
func login(conn net.Conn, nick string) (client user, err error) {
    if len(nick) > 0 {
        
        if isNickTaken(nicknames, nick) == true {
            return user{}, errors.New("nickname already taken")
        }
        
        client := user{nick: nick, conn: conn}
        nicknames = append(nicknames, client)
        
        sendRes(conn, fmt.Sprintf("Welcome %s", client.nick), client)
        sendToAllButMe( fmt.Sprintf("%s joind the chat", client.nick) , client)
        
        log.Println("Joined: "+client.nick)
        log.Println(nicknames)
        
        return client, nil
    }
    
    return user{}, errors.New("invalid nick name: " + nick)
}

// Logout client and close connection
func logout(conn net.Conn, client user) {
    sendRes(conn, "Bye!", client)
    cleanClient(client)
}

// sendRes sends message to client
func sendRes(conn net.Conn, message string, client user) {
    message = strings.Trim(message, "\n")
    log.Printf("Sending '%s' >> %s\n", message, client.nick)
    _, err := conn.Write([]byte(message + "\n"))
    if err != nil {
        log.Printf("Error writing to client: %s\n", err)
    }
}


// sendToAll sends message to all connected users
func sendToAll(message string, client user) {
    for _, c := range nicknames {
        if c.conn != nil {
            sendRes(c.conn, message, c)
        }
    }
}

// sendToAllButMe sends message to all clients except the one sending the requests
func sendToAllButMe(message string, client user) {
    for _, c := range nicknames {
        if c.conn != nil && c.nick != client.nick {
            sendRes(c.conn, message, c)
        }
    }
}

// isNickTaken checks if nickname is already taken and return 
// true if it is or fales if it isn't
func isNickTaken(nicknames []user, nick string) (bool) {
    for _, c := range nicknames {
        if c.nick == nick {
            return true
        }
    }
    return false;
}

// cleanClient removes client from a list of nicknames 
// and notifies all other clients that he left.
func cleanClient(client user) {
    log.Println("Left: "+client.nick)
    
    sendToAllButMe(fmt.Sprintf("%s left the chat", client.nick), client)
    
    nicknames = removeNickname(nicknames, client.nick)
    
    log.Println(nicknames)
}

// removeNickname removes cliend by nickname and returns 
// an array without hit client.
func removeNickname(slice []user, n string) ([]user) {
    slice_i := -1
    for i,u := range slice {
        if u.nick == n {
            slice_i = i
            break;
        }
    }
    
    if slice_i >= 0 {
        return append(slice[:slice_i], slice[slice_i+1:]...)
    }
    return slice
}
