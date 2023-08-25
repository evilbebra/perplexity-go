# perplexity-go
Go Unofficial API for Perplexity.ai. Alternative to chat-gt for free and without api keys.

# Usage
Just <code>go get -u github.com/evilbebra/perplexity-go</code> and use it like this:
```go
package main

import (
	"fmt"
	perplexity "github.com/evilbebra/perplexity-go"
)

func main() {
	resp, err := perplexity.NewQuestion("How old is Golang?")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Your uuid: %s\n", resp.Response.UUID)
		fmt.Printf("Text of your question was: %s\n", resp.Response.QueryStr)

		for i, query := range resp.Response.RelatedQueries {
			fmt.Printf("Related Query %d: %s\n", i, query)
		}
		fmt.Printf("SearchFocus: %s\n", resp.Response.SearchFocus)
		fmt.Printf("Mode: %s\n", resp.Response.Mode)
		fmt.Printf("Mode: %v\n", resp.Response.Gpt4)
		fmt.Printf("Web Result for query: %s\n", resp.Answer.ExtraWebResults)
		fmt.Printf("Answer: %s\n", resp.Answer.Text)
	}

  /*
    OUTPUT:
      Your uuid: 3d1e0e45-9ae2-4b0b-bdc9-4e42c457fbdb
      Text of your question was: How old is Golang?  
      Related Query 0: what are the unique features of Golang
      Related Query 1: how has Golang evolved since its release
      Related Query 2: what are some popular applications of Golang
      SearchFocus: internet
      Mode: concise
      Mode: false
      Web Result for query: [{Go Turns 10 - The Go Programming Language
      https://go.dev/blog/10years Celebrating 10 years of Go makes me think
      back to early November 2009, when we were getting ready to share Go with
      the world. We didn't know ...  <nil>}]
      Answer: The Go programming language, also known as Golang, was developed
      at Google by Robert Griesemer, Rob Pike, and Ken Thompson[1]. It was first
      publicly announced in November 2009[1]. The initial version, Go 1.0, was released
       in March 2012[1]. Therefore, as of the current date, Go is approximately 11 years old.
      Despite its age, Go continues to gain popularity and is widely used in production
      at Google and many other organizations and open-source projects[1][2].
  */

	resp, err = perplexity.NewQuestion("What is UDP?")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Your uuid: %s\n", resp.Response.UUID)
		fmt.Printf("Text of your question was: %s\n", resp.Response.QueryStr)

		for i, query := range resp.Response.RelatedQueries {
			fmt.Printf("Related Query %d: %s\n", i, query)
		}
		fmt.Printf("SearchFocus: %s\n", resp.Response.SearchFocus)
		fmt.Printf("Mode: %s\n", resp.Response.Mode)
		fmt.Printf("Web Result for query: %s\n", resp.Answer.ExtraWebResults)
		fmt.Printf("Answer: %s\n", resp.Answer.Text)
	}

  /*
    OUTPUT:
      Your uuid: a34b40c2-9029-42fc-96af-cd8bfba31c72
      Text of your question was: What is UDP?
      Related Query 0: what is the difference between UDP and TCP
      Related Query 1: what are some common applications that use UDP
      Related Query 2: how does UDP handle packet loss
      SearchFocus: internet
      Mode: concise
      Web Result for query: []
      Answer: User Datagram Protocol (UDP) is a communication protocol
      used across the Internet for time-sensitive transmissions such as
      video playback, gaming, and Domain Name System (DNS) lookups[2][5].
      UDP is a simple message-oriented transport layer protocol that is documented
      in RFC 768[1]. Unlike Transmission Control Protocol (TCP), UDP is an unreliable
      and connectionless protocol, meaning that it does not require prior communication
      to set up communication channels or data paths[1][3][4]. UDP provides two services
      not provided by the IP layer: port numbers to help distinguish different user
      requests and an optional checksum capability to verify that the data has not been
      corrupted[3]. UDP is beneficial in time-sensitive communications because it enables
      the transfer of data before an agreement is provided by the receiving party,
      resulting in faster data transfer speeds[3][5]. However, UDP provides no guarantees
      to the upper layer protocol for message delivery, and the UDP layer retains no state
      of UDP messages once sent, making it an unreliable protocol[1].
      If transmission reliability is desired, it must be implemented in the user's application[1].
  */
}

```
