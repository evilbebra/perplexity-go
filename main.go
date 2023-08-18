package perplexity

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
