package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type SearchFocus string

const (
	Internet     SearchFocus = "internet"
	Writing      SearchFocus = "writing"
	Academic     SearchFocus = "scholar"
	WolframAlpha SearchFocus = "wolfram"
	YouTube      SearchFocus = "youtube"
	Reddit       SearchFocus = "reddit"
)

type SearchSource string

const (
	Default SearchSource = "default"
)

type SearchMode string

const (
	Concise SearchMode = "concise"
	Copilot SearchMode = "copilot"
)

type Session struct {
	Sid            string
	Wss            *websocket.Conn
	Client         *http.Client
	UserAgent      string
	BaseApiURI     *url.URL
	AskSeqNum      int
	MessageChannel chan AnswerResponse
	Mutex          sync.Mutex // Добавляем мьютекс
}

func NewSession() (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	rawURL := "https://www.perplexity.ai/socket.io/"
	baseApiURI, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	session := Session{
		Client: &http.Client{
			Jar: jar,
		},
		UserAgent:  userAgent(),
		BaseApiURI: baseApiURI,
		AskSeqNum:  1,
	}

	baseApiURI.RawQuery = addParams(map[string]string{
		"EIO":       "4",
		"transport": "polling",
	})

	req, err := http.NewRequest("GET", baseApiURI.String(), nil)

	bodyBytes, err := session.handleRequest(req)
	if err != nil {
		return nil, err
	}

	var result GetSidResponse
	if err := json.Unmarshal(bodyBytes[1:], &result); err != nil {
		return nil, err
	}

	session.Sid = result.Sid

	return &session, nil
}

func (s *Session) CheckConnection() error {
	s.BaseApiURI.RawQuery = addParams(map[string]string{
		"EIO":       "4",
		"transport": "polling",
		"sid":       s.Sid,
	})

	body := bytes.NewBufferString("40")

	req, err := http.NewRequest("POST", s.BaseApiURI.String(), body)

	bodyBytes, err := s.handleRequest(req)
	if err != nil {
		return err
	}

	if string(bodyBytes) != "OK" {
		return errors.New("Session Check Failed")
	}

	return nil
}

func (s *Session) GetSid() error {
	s.BaseApiURI.RawQuery = addParams(map[string]string{
		"EIO":       "4",
		"transport": "polling",
		"sid":       s.Sid,
	})

	req, err := http.NewRequest("GET", s.BaseApiURI.String(), nil)
	if err != nil {
		return err
	}

	bodyBytes, err := s.handleRequest(req)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bodyBytes), "40{\"sid\":") {
		return errors.New("Get Sid Failed")
	}
	return nil
}

func (s *Session) Close() {
	s.Wss.WriteMessage(websocket.CloseMessage, []byte{})
	s.Wss.Close()
}

func (s *Session) Ask(question string) error {
	code := s.AskSeqNum
	defer func(s *Session) { s.AskSeqNum += 1 }(s)

	askReq := AskRequest{
		Source:      Default,
		Language:    "en-US",
		Timezone:    "Asia/Shanghai",
		SearchFocus: Internet,
		Gpt4:        false,
		Mode:        Concise,
	}

	marshalled, err := json.Marshal(askReq)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(
		"%d%d[%q,%q,%v]",
		42, code, "perplexity_ask", question, string(marshalled),
	)

	s.Mutex.Lock()         // Захватываем мьютекс перед операцией записи
	defer s.Mutex.Unlock() // Освобождаем мьютекс после операции записи

	err = s.Wss.WriteMessage(websocket.TextMessage, []byte(q))
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ReadAnswerLoop() error {
	for {
		_, message, err := s.Wss.ReadMessage()
		if err != nil {
			return err
		}

		if bytes.Contains(message, []byte("query_answered")) {
			var result AskResponse
			if err := parseMessage(message, &result); err != nil {
				return err
			}

			var answer AnswerDetails
			if err := json.Unmarshal([]byte(result.Text), &answer); err != nil {
				return err
			}

			_, _, _ = s.Wss.ReadMessage()

			s.MessageChannel <- AnswerResponse{
				Answer:   answer,
				Response: result,
			}

		} else {
			log.Println("Something is wrong with the execution status message", string(message[:50]))
		}
	}
}

func (s *Session) handleRequest(req *http.Request) ([]byte, error) {
	setupCommonHeaders(req, s.UserAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func setupCommonHeaders(req *http.Request, userAgent string) {
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("user-agent", userAgent)
}

func parseMessage(message []byte, v any) error {
	start := strings.Index(string(message), ",") + 1
	respBytes := message[start : len(message)-1]

	if err := json.Unmarshal(respBytes, v); err != nil {
		return err
	}

	return nil
}

func concatenateCookies(jar http.CookieJar) (*string, error) {
	url, err := url.Parse("https://www.perplexity.ai")
	if err != nil {
		return nil, err
	}
	cookies := jar.Cookies(url)

	var cookieStrings []string
	for _, cookie := range cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}

	cookieString := strings.Join(cookieStrings, "; ")

	return &cookieString, nil
}

func addParams(params map[string]string) string {
	p := url.Values{}
	for key, value := range params {
		p.Add(key, value)
	}
	return p.Encode()
}

func userAgent() string {
	return fmt.Sprintf("")
}
