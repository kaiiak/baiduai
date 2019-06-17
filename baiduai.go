// Package baiduai restful api wrapper
package baiduai

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	baiduOauthAPI = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s&"
)

// ErrorResp of api
type ErrorResp struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// Token response of baidu ai auth
type Token struct {
	RefreshToken  string `json:"refresh_token,omitempty"`
	ExpiresIn     int64  `json:"expires_in,omitempty"`
	Scope         string `json:"scope,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
	*ErrorResp
}

// Client of baidu ai
type Client struct {
	id, secret  string
	accessToken string
	tokenLock   sync.RWMutex
	expriesAt   time.Time
}

// New baidu ai client
func New(id, secret string) *Client {
	c := &Client{}
	c.id = id
	c.secret = secret

	return c
}

func (c *Client) requestAccessToken() error {
	c.tokenLock.Lock()
	defer c.tokenLock.Unlock()
	s, err := httpGet(fmt.Sprintf(baiduOauthAPI, c.id, c.secret), nil, nil)
	if err != nil {
		return err
	}
	token := &Token{}
	if err = s.Scan(token); err != nil {
		return err
	}
	if token.ErrorResp != nil {
		return errors.New(token.ErrorDescription)
	}
	c.accessToken = token.AccessToken
	c.expriesAt = c.expriesAt.Add(time.Second * time.Duration(token.ExpiresIn))
	return nil
}

// GetAccessToken if token expries will refresh token
func (c *Client) GetAccessToken() (string, error) {
	c.tokenLock.RLock()
	defer c.tokenLock.RUnlock()
	if time.Now().After(c.expriesAt) {
		if err := c.requestAccessToken(); err != nil {
			return "", err
		}
	}
	return c.accessToken, nil
}
