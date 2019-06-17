package baiduai

import (
	"errors"
	"io"
	"net/url"
	"strconv"
)

const (
	defaultDevPid = 1536
	defaultCuid   = "golang-baiduai"
	voiceURL      = "https://vop.baidu.com/server_api"
)

type Voice struct {
	r      io.Reader
	devPid int
	c      *Client
	cuid   string
}

func NewVoice(client *Client) *Voice {
	return &Voice{c: client}
}

func (v *Voice) SetDevPid(pid int) *Voice {
	v.devPid = pid
	return v
}

func (v *Voice) SetCuid(cuid string) *Voice {
	v.cuid = cuid
	return v
}

func (v *Voice) SetRaw(r io.Reader) *Voice {
	v.r = r
	return v
}

type VoiceResp struct {
	ErrNo    int      `json:"err_no,omitempty"`
	ErrMsg   string   `json:"err_msg,omitempty"`
	SN       string   `json:"sn,omitempty"`
	CorpusNo string   `json:"corpus_no,omitempty"`
	Result   []string `json:"result,omitempty"`
}

func (v *Voice) Resp() ([]string, error) {
	if v.devPid == 0 {
		v.devPid = 1536
	}
	if v.cuid == "" {
		v.cuid = defaultCuid
	}
	if v.r == nil {
		return nil, errors.New("")
	}
	query := make(url.Values)
	query.Set("dev_pid", strconv.Itoa(v.devPid))
	query.Set("cuid", v.cuid)
	token, err := v.c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	query.Set("token", token)

	s, err := httpPost(voiceURL, map[string]string{
		"Content-Type": "audio/pcm; rate=16000",
	}, query, v.r)
	if err != nil {
		return nil, err
	}
	var resp = new(VoiceResp)
	if err = s.Scan(v); err != nil {
		return nil, err
	}
	if resp.ErrMsg != "success." {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp.Result, nil
}
