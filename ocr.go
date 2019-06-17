package baiduai

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	ocrBaseURL      = "https://aip.baidubce.com/rest/2.0/ocr/v1/"
	defaultLanguage = "CHN_ENG"
)

type ocr struct {
	r io.Reader
	c *Client
	t string
}

func (o *ocr) respScaner(cQuery map[string]string) (Scaner, error) {
	u := fmt.Sprintf("%s%s", ocrBaseURL, o.t)
	query := make(url.Values)
	token, err := o.c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	query.Set("access_token", token)
	for k, v := range cQuery {
		query.Set(k, v)
	}

	form := make(url.Values)
	b, err := ioutil.ReadAll(o.r)
	if err != nil {
		return nil, err
	}
	form.Set("image", base64.StdEncoding.EncodeToString(b))

	req, err := wrapRequest(http.MethodPost, u, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}, query, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytesScaner(bs), nil
}

type OcrGeneralBasic struct {
	ocr
	languageType string
}

/*
{
	"log_id": 2471272194,
	"words_result_num": 2,
	"words_result":
		[
			{"words": " TSINGTAO"},
			{"words": "青島睥酒"}
		]
}
*/
type GeneralBasicResp struct {
	LogID          int64 `json:"log_id,omitempty"`
	WordsResultNum int   `json:"words_result_num,omitempty"`
	WordsResult    []struct {
		Words string `json:"words,omitempty"`
	} `json:"words_result,omitempty"`
}

func NewOcrGeneralBasic(c *Client, image io.Reader) *OcrGeneralBasic {
	return &OcrGeneralBasic{
		ocr: ocr{
			c: c,
			r: image,
			t: "general_basic",
		},
	}
}

func (o *OcrGeneralBasic) SetLanguage(lt string) *OcrGeneralBasic {
	o.languageType = lt
	return o
}

func (o *OcrGeneralBasic) Resp() (*GeneralBasicResp, error) {
	if o.r == nil {
		return nil, errors.New("")
	}
	if o.languageType == "" {
		o.languageType = defaultLanguage
	}
	s, err := o.respScaner(map[string]string{
		"language_type": o.languageType,
	})
	if err != nil {
		return nil, err
	}
	var result = new(GeneralBasicResp)
	if err = s.Scan(result); err != nil {
		return nil, err
	}
	return result, nil
}
