package main

import (
	"fmt"
	"github.com/evilbebra/perplexity-go/internal"
)

func NewQuestion(question string) (internal.AnswerResponse, error) {
	s, err := internal.NewSession()
	if err != nil {
		return internal.AnswerResponse{}, err
	}

	if !initConnection(s) {
		return internal.AnswerResponse{}, fmt.Errorf("cant init connection")
	}

	defer s.Close()

	messageChannel := make(chan internal.AnswerResponse)
	s.MessageChannel = messageChannel

	go func() {
		err := s.ReadAnswerLoop()
		if err != nil {
			fmt.Errorf("read answer in AskMe: %v", err)
		}
	}()

	err = s.Ask(question)
	if err != nil {
		return internal.AnswerResponse{}, err
	}

	resp := <-s.MessageChannel
	return resp, nil
}

func initConnection(s *internal.Session) bool {
	var err error

	if err = s.CheckConnection(); err != nil {
		return false
	}

	if err = s.GetSid(); err != nil {
		return false
	}

	if err = s.InitWss(); err != nil {
		return false
	}
	return true
}

func main() {
	resp, err := NewQuestion("How old is Golang?")
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

	resp, err = NewQuestion("What is UDP?")
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

}
