// Author: Gregor Flajs <flajsg@gmail.com> 2017. 

/*
This is simple chat client that connects on TCP localhost:8080 and relays 
requests (chat messages) to server and checks if server has send any response.

To exit the client just type "exit".
*/
package main

import (
    "net"
    "log"
    "os"
    "fmt"
    "bufio"
    "strings"
    "time"
)

func main() {
    
    log.SetFlags(log.Lshortfile)
 
    servAddr := "localhost:8080"
    
    conn, err := net.Dial("tcp", servAddr)

    if err != nil {
        log.Fatal("Error dialing:", err)
    }
    
    // reader from standard input
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Println("Hello from gochat-client!")

    // read messages from server and display them as they come
    go readFromServer(conn)

    // chat in loop until you send exit command
    for {
        chat, err := reader.ReadString('\n')
        
        if err != nil {
            fmt.Printf("Error reading from stdin: %s\n", err)
            continue
        }
        
        // Remove trailing newline
        chat = strings.Trim(string(chat), "\n")
        
        sendRequest(conn, chat)
        
        // Close client on exit
        if strings.Compare(chat, "exit") == 0 {
            fmt.Println("Exiting...")
            time.Sleep(time.Second)
            break;
        }
    }

    if tcpcon, ok := conn.(*net.TCPConn); ok {
        tcpcon.CloseWrite()
    }
    
    err = conn.Close()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Successfull exit. Have a nice day.")
}

// sendRequest sends chat message to server
func sendRequest(conn net.Conn, text string) {
    message := text;
    
    if _,err := conn.Write([]byte(message + "\n")); err != nil {
        log.Fatal(err)
    }
}


// readFromServer waits for a response from server and displays it
func readFromServer(conn net.Conn) {
    for {
        message, err := bufio.NewReader(conn).ReadString('\n')
        
        if err != nil {
            fmt.Println("Server down:", err)
            return
        }
        
        message = strings.Trim(string(message), "\n")
    
        fmt.Println("\r<<", string(message))
    }
}
