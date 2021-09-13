package eth

import (
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/utils"
)

var ErrInvalidTopics = errors.New("Invalid topics")

/**
translateTopics takes in an ethReq's topics field and translates it to a it's equivalent QtumReq
topics (optional) has a max lenght of 4

Topics are order-dependent. A transaction with a log with topics [A, B] will be matched by the following topic filters:

    [] “anything”
    [A] “A in first position (and anything after)”
    [null, B] “anything in first position AND B in second position (and anything after)”
    [A, B] “A in first position AND B in second position (and anything after)”
    [[A, B], [A, B]] “(A OR B) in first position AND (A OR B) in second position (and anything after)”
*/
func TranslateTopics(ethTopics []interface{}) ([][]string, error) {

	var topics [][]string
	nilCount := 0

	for _, topic := range ethTopics {
		switch topic.(type) {
		case []string:
			stringTopics := []string{}
			for _, t := range topic.([]string) {
				stringTopics = append(stringTopics, utils.RemoveHexPrefix(t))
			}
			topics = append(topics, stringTopics)
		case string:
			topics = append(topics, []string{utils.RemoveHexPrefix(topic.(string))})
		case nil:
			nilCount++
			topics = append(topics, nil)
		case []interface{}:
			stringTopics := []string{}
			for _, t := range topic.([]interface{}) {
				if stringTopic, ok := t.(string); ok {
					stringTopics = append(stringTopics, utils.RemoveHexPrefix(stringTopic))
				} else {
					return nil, ErrInvalidTopics
				}
			}
			topics = append(topics, stringTopics)
		}
	}

	// QTUM rpc requires at least one topic to match
	// null topics don't count as matching they count as skipping
	// so if all topics are null, QTUM rpc will return an empty result as nothing will match
	if nilCount == len(topics) {
		// return an empty list of topics
		return [][]string{}, nil
	}

	return topics, nil

}
