package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
)

// type UserList {
//     Page   int      `json:"data.user.result.timeline.timeline.instructions"`
// }

type GraphRes struct {
	Data Data `json:"data"`
}

type Data struct {
	User DUser `json:"user"`
}

type DUser struct {
	Result DUResult `json:"result"`
}
type DUResult struct {
	TypeName string      `json:"__typename"`
	Timeline DURTimeline `json:"timeline"`
}

type DURTimeline struct {
	Timeline DURTTimeline `json:"timeline"`
}

type DURTTimeline struct {
	Instructions []TimelineAddEntries `json:"instructions"`
}

type TimelineAddEntries struct {
	Type    string       `json:"type"`
	Entries []TAEEntries `json:"entries"`
}

// Cousin of Cursor
type TAEEntries struct {
	EntryID   string     `json:"entryId"`
	SortIndex string     `json:"sortIndex"`
	Content   TAEContent `json:"content"`
}

type TAEContent struct {
	EntryType string `json:"entryType"`
	TypeName  string `json:"__typename"`
	// Entry can either have a user or a Cursor
	ItemContent TAEItemContent `json:"itemContent,omitempty"`
	Value       string         `json:"value,omitempty"`
	CursorType  string         `json:"cursorType,omitempty"`
}

type TAEItemContent struct {
	ItemType    string      `json:"itemType"`
	UserResults UserResults `json:"user_results,omitempty"`
}

type UserResults struct {
	Result Result `json:"result"`
}

type Result struct {
	ID           string       `json:"rest_id"`
	IsVerified   bool         `json:"is_blue_verified"`
	LegacyResult LegacyResult `json:"legacy"`
}

type LegacyResult struct {
	FollowedBy     bool   `json:"followed_by"`
	Following      bool   `json:"following"`
	Description    string `json:"description"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	ScreenName     string `json:"screen_name"`
	Name           string `json:"name"`
}

var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}

func envPreflightCheck() error {
	keys := []string{
		"CURSOR_FLOOR",
		"CURSOR_CIEL",
		"X_CSRF_TOKEN",
		"X_BEARER_AUTH_TOKEN",
		"X_AUTH_TOKEN",
		"X_KDT",
		"X_USER_ID",
	}

	for _, key := range keys {
		if os.Getenv(key) == "" {
			return errors.New("expected " + key + " Environment Variable to be set.")
		}
	}

	return nil
}

type Variables struct {
	UserId                 string `json:"userId"`
	Count                  int    `json:"count"`
	Cursor                 string `json:"cursor"`
	IncludePromotedContent bool   `json:"includePromotedContent"`
}

func (v Variables) ToString() string {
	vars, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
		return ""
	}
	return string(vars)
}

func (v Variables) ToParams() string {
	return url.QueryEscape(v.ToString())
}

func prepareQuery(cursorFloor string, cursorCiel string) Variables {

	cursor := strings.Join([]string{cursorFloor, cursorCiel}, "|")

	userID := os.Getenv("X_USER_ID")

	var variables Variables = Variables{
		UserId:                 userID,
		Count:                  20,
		Cursor:                 cursor,
		IncludePromotedContent: false,
	}
	return variables
}

func prepareReq(req http.Request) {
	csrfToken := os.Getenv("X_CSRF_TOKEN")
	bearerAuth := os.Getenv("X_BEARER_AUTH_TOKEN")
	authToken := os.Getenv("X_AUTH_TOKEN")
	kdtVal := os.Getenv("X_KDT")

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "zstd")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-twitter-auth-type", "OAuth2Session")
	req.Header.Add("x-csrf-token", csrfToken)
	req.Header.Add("x-twitter-client-language", "en")
	req.Header.Add("x-client-transaction-id", "en")
	req.Header.Add("x-twitter-active-user", "yes")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("authorization", "Bearer "+bearerAuth)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", `"kdt=`+kdtVal+`; ct0=`+csrfToken+`; auth_token=`+authToken+`; lang=en"`)
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "trailers")

}

func featsToParams(features map[string]bool) string {
	feats, err := json.Marshal(features)

	if err != nil {
		log.Println(err)
		return ""
	}
	return url.QueryEscape(string(feats))
}

func magicUrlBuilder(baseUrl string, vars Variables, features map[string]bool) string {

	features_params := featsToParams(features)

	magicUrl := baseUrl + "?variables=" + vars.ToParams() + "&features=" + features_params

	return magicUrl
}

func rangeGet(cursorFloor string, cursorCiel string) GraphRes {
	variables := prepareQuery(cursorFloor, cursorCiel)

	features := map[string]bool{"rweb_tipjar_consumption_enabled": true, "responsive_web_graphql_exclude_directive_enabled": true, "verified_phone_label_enabled": false, "creator_subscriptions_tweet_preview_api_enabled": true, "responsive_web_graphql_timeline_navigation_enabled": true, "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false, "communities_web_enable_tweet_community_results_fetch": true, "c9s_tweet_anatomy_moderator_badge_enabled": true, "articles_preview_enabled": true, "responsive_web_edit_tweet_api_enabled": true, "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true, "view_counts_everywhere_api_enabled": true, "longform_notetweets_consumption_enabled": true, "responsive_web_twitter_article_tweet_consumption_enabled": true, "tweet_awards_web_tipping_enabled": false, "creator_subscriptions_quote_tweet_preview_enabled": false, "freedom_of_speech_not_reach_fetch_enabled": true, "standardized_nudges_misinfo": true, "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true, "rweb_video_timestamps_enabled": true, "longform_notetweets_rich_text_read_enabled": true, "longform_notetweets_inline_media_enabled": true, "responsive_web_enhance_cards_enabled": false}

	baseurl := "https://x.com/i/api/graphql/OSXFkKmGvfw_6pGgGtkWFg/Followers"

	magicUrl := magicUrlBuilder(baseurl, variables, features)

	gurl, err := url.Parse(magicUrl)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", gurl.String(), nil)
	if err != nil {
		panic(err)
	}

	prepareReq(*req)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	// Rate limited! Take a break and continue with another request.
	if res.StatusCode == 429 {

		resetHeader := res.Header.Get("x-rate-limit-reset")
		resetEpoch, err := strconv.ParseInt(resetHeader, 10, 64)

		if err != nil {
			log.Printf("Failed to convert rate limit reset timestamp to number: %s", err)
		}

		t := time.Unix(resetEpoch, 0)

		relativeTime := time.Until(t)

		log.Printf("Rate limit reached! Gonna chill out for %d seconds", int(relativeTime.Seconds()))

		time.Sleep(relativeTime)
		return rangeGet(cursorFloor, cursorCiel)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	decomp, err := Decompress(body)

	if err != nil {
		panic(err)
	}
	var dat GraphRes

	if err := json.Unmarshal(decomp, &dat); err != nil {
		panic(err)
	}

	return dat
}

func main() {
	err := envPreflightCheck()

	if err != nil {
		log.Fatal(err)
	}

	cursorFloor := os.Getenv("CURSOR_FLOOR")
	cursorCiel := os.Getenv("CURSOR_CIEL")

	var users []Result = make([]Result, 0)

	var res GraphRes

	i := 0
	for i < 1000 {
		res = rangeGet(cursorFloor, cursorCiel)

		entries := res.Data.User.Result.Timeline.Timeline.Instructions[0].Entries

		if len(entries) == 0 {
			wait, err := time.ParseDuration("10s")
			if err != nil {
				log.Println("Your hardcoded duration is broken. Oops.")
			}
			log.Printf("Something fell off the rails.... Let's try again in %d seconds", int(wait.Seconds()))

			time.Sleep(wait)
			continue
		}

		usersList := entries[:len(entries)-2]

		bottom := entries[len(entries)-2]
		top := entries[len(entries)-1]

		log.Println(bottom)
		log.Println(top)

		nextCursor := strings.Split(bottom.Content.Value, "|")

		cursorFloor = nextCursor[0]
		cursorCiel = nextCursor[1]

		for _, e := range usersList {
			record := e.Content.ItemContent.UserResults.Result

			users = append(users, record)
		}
		i++
	}

	for _, user := range users {
		log.Println(user.ID)
	}
}
