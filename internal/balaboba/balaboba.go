package balaboba

import (
	"encoding/json"
	"fmt"
	"github.com/mispon/stewart-bot/internal/utils"
	"io/ioutil"
)

type (
	Balaboba interface {
		GetText(query string, intro int) (string, error)
	}

	balaboba struct {
		url string
	}
)

// New creates new instance
func New(url string) Balaboba {
	return &balaboba{
		url: url,
	}
}

// GetText returns generated text for query and intro
func (b balaboba) GetText(query string, intro int) (string, error) {
	var (
		request, _ = json.Marshal(map[string]interface{}{
			"query":  query,
			"intro":  intro,
			"filter": 0,
		})
	)

	res, err := utils.SendPost(b.url, request)
	if err != nil {
		return "", err
	}
	defer utils.Close(res.Body.Close)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	quote := struct {
		BadQuery int
		Error    int
		Query    string
		Text     string
	}{}
	if err = json.Unmarshal(body, &quote); err != nil {
		return "", err
	}

	var text string
	if quote.Error != 0 {
		text = "не могу, error > 0 :("
	} else if quote.BadQuery > 0 || len(quote.Text) == 0 {
		text = "мне нечего сказать на это"
	} else {
		text = fmt.Sprintf(`%s`, quote.Text)
	}

	return text, nil
}
