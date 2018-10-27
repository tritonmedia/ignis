package trello

import (
	"log"

	"github.com/adlio/trello"
)

var client *trello.Client

var idlist string

// Init the trello client
func Init(appKey string, token string, list string) {
	client = trello.NewClient(appKey, token)
	idlist = list

	log.Printf("[trello/init] using list: %s", idlist)
}

// CreateCard create a request card.
func CreateCard(title string, source string, isMovie bool) (string, error) {
	labels := []string{"ignis"}

	if isMovie {
		labels = append(labels, "Movie")
	}

	desc := "DOWNLOAD: [fixme]()"
	if source != "" {
		desc = "DOWNLOAD: [http}(" + source + ")"
	}

	list, err := client.GetList(idlist, trello.Defaults())
	if err != nil {
		log.Printf("[trello/createcard]: failed to get list: %s", err.Error())
		return "", err
	}

	card := trello.Card{
		Name: title,
		Desc: desc,
		Pos:  1,
	}

	err = list.AddCard(&card, trello.Defaults())
	if err != nil {
		return "", err
	}

	return "https://trello.com/c/" + card.ID, err
}
