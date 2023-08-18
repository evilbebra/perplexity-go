package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func (s *Session) InitWss() error {
	if err := s.switchToWebSocket(); err != nil {
		return err
	}

	if err := s.handleWssHandshake(); err != nil {
		return err
	}

	return nil
}

func (s *Session) switchToWebSocket() error {
	s.BaseApiURI.Scheme = "wss"
	defer func(u *url.URL) { u.Scheme = "https" }(s.BaseApiURI)

	cookie, err := concatenateCookies(s.Client.Jar)
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Add("Cookie", *cookie)

	s.BaseApiURI.RawQuery = addParams(map[string]string{
		"EIO":       "4",
		"transport": "websocket",
		"sid":       s.Sid,
	})

	wssURI := s.BaseApiURI.String()
	conn, resp, err := websocket.DefaultDialer.Dial(wssURI, header)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		log.Println("Error switching protocols:", resp.StatusCode)
		return err
	}

	s.Wss = conn

	return nil
}

func (s *Session) handleWssHandshake() error {
	s.Wss.WriteMessage(websocket.TextMessage, []byte("2probe"))
	for {
		_, message, err := s.Wss.ReadMessage()
		if err != nil {
			continue
		}

		if string(message) == "6" {
			break
		}

		if string(message) == "3probe" {
			s.Wss.WriteMessage(websocket.TextMessage, []byte("5"))
		}
	}

	return nil
}
