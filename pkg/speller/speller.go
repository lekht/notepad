package speller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/lekht/notepad/config"
)

var ErrSpelling = errors.New("there are spelling errors in the text")

type SpellChecker struct {
	address  string
	attempts int
	client   *http.Client

	languages []string
	options   string
	format    string
}

func New(cfg *config.Speller) *SpellChecker {
	c := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}

	return &SpellChecker{
		address:  cfg.Address,
		attempts: cfg.Attempts,
		client:   c,

		languages: cfg.Langs,
		options:   cfg.Options,
		format:    cfg.Format,
	}
}

type SpellRequest struct {
	Text    string `json:"text"`
	Lang    string `json:"lang"`
	Options string `json:"options"`
	Format  string `json:"format"`
}

type Mistakes struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func (sc *SpellChecker) Check(title, body string) (bool, error) {
	if len(sc.languages) < 1 {
		return false, errors.New("language not specified")
	}

	language := sc.languages[0]

	if len(sc.languages) > 1 {
		for i := 1; i < len(sc.languages); i++ {
			language += fmt.Sprintf(",%s", sc.languages[i])
		}
	}
	text := title + " " + body
	newText := strings.ReplaceAll(text, " ", "+")

	resp, err := http.PostForm(sc.address, url.Values{
		"text":    {newText},
		"lang":    {language},
		"format":  {sc.format},
		"options": {sc.options},
	})
	if err != nil {
		return false, errors.Wrap(err, " failed to make post req")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to read response body")
	}

	var mistakes []Mistakes
	if err = json.Unmarshal(respBody, &mistakes); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal response to struct")
	}

	if len(mistakes) == 0 {
		return true, nil
	}
	return false, nil
}
