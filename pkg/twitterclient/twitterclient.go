package twitterclient

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/niceb5y/heart-extractor/pkg/auth"
	"github.com/niceb5y/heart-extractor/pkg/crypto"
	"github.com/niceb5y/heart-extractor/pkg/db"
)

const fetchCount = 200

func GetClient() *twitter.Client {
	consumerKey, _ := crypto.Decrypt(auth.EncryptedConsumerKey)
	consumerSecret, _ := crypto.Decrypt(auth.EncryptedConsumerSecret)

	var token *oauth1.Token
	token = db.GetAuthToken()

	if token == nil {
		token = auth.CreateAuthToken()
		db.SetAuthToken(token)
	}

	config := oauth1.NewConfig(consumerKey,
		consumerSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}

func Fetch(twitterClient *twitter.Client, favoriteListParams *twitter.FavoriteListParams, mediaURLs *[]string) {
	if favoriteListParams == nil {
		favoriteListParams = &twitter.FavoriteListParams{Count: fetchCount}
	}

	favoriteList, _, err := twitterClient.Favorites.List(favoriteListParams)
	if err != nil {
		panic(err)
	}

	for listIdx, favoriteListItem := range favoriteList {
		extendedEntities := favoriteListItem.ExtendedEntities
		if extendedEntities != nil {
			for _, media := range extendedEntities.Media {
				mediaURL := media.MediaURLHttps
				if len(media.VideoInfo.Variants) > 0 {
					bitrate := -1
					for _, variant := range media.VideoInfo.Variants {
						if variant.Bitrate > bitrate {
							bitrate = variant.Bitrate
							mediaURL = variant.URL
						}
					}
				} else {
					if media.Type == "photo" {
						mediaURL += "?name=orig"
					}
				}
				*mediaURLs = append(*mediaURLs, mediaURL)
			}
		}
		if listIdx == len(favoriteList)-1 {
			if favoriteListParams.MaxID != favoriteListItem.ID {
				newFavoriteListParams := &twitter.FavoriteListParams{Count: fetchCount, MaxID: favoriteListItem.ID - 1}
				Fetch(twitterClient, newFavoriteListParams, mediaURLs)
			}
		}
	}
}
