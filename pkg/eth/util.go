package eth

import (
	"github.com/qtumproject/janus/pkg/utils"
)

// translateTopics takes in an ethReq's topics field and translates it to a it's equivalent QtumReq
// topics (optional) has a max lenght of 4
func TranslateTopics(ethTopics []interface{}) ([]interface{}, error) {

	var topics []interface{}

	for _, topic := range ethTopics {
		switch topic.(type) {
		case []interface{}:
			topic, err := TranslateTopics(topic.([]interface{}))
			if err != nil {
				return nil, err
			}
			topics = append(topics, topic...)
		case string:
			topics = append(topics, utils.RemoveHexPrefix(topic.(string)))
		case nil:
			topics = append(topics, nil)
		}
	}

	return topics, nil

}
