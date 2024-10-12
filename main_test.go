package main

import (
	"testing"
)

func TestToParam(t *testing.T) {
	var tcs = []struct {
		variables Variables
		want      string
	}{
		{
			variables: Variables{
				UserId:                 "1234567890",
				Count:                  20,
				Cursor:                 "1234567891234567890|1234567891234567890",
				IncludePromotedContent: false,
			},
			want: "%7B%22userId%22%3A%221234567890%22%2C%22count%22%3A20%2C%22cursor%22%3A%221234567891234567890%7C1234567891234567890%22%2C%22includePromotedContent%22%3Afalse%7D",
		},
	}

	for _, tc := range tcs {
		got := tc.variables.ToParams()
		if got != tc.want {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}

func TestDefaultFeatures(t *testing.T) {
	var tcs = []struct {
		features map[string]bool
		want     string
	}{
		{
			features: map[string]bool{
				"articles_preview_enabled":                                                true,
				"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
				"communities_web_enable_tweet_community_results_fetch":                    true,
				"creator_subscriptions_quote_tweet_preview_enabled":                       false,
				"creator_subscriptions_tweet_preview_api_enabled":                         true,
				"freedom_of_speech_not_reach_fetch_enabled":                               true,
				"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
				"longform_notetweets_consumption_enabled":                                 true,
				"longform_notetweets_inline_media_enabled":                                true,
				"longform_notetweets_rich_text_read_enabled":                              true,
				"responsive_web_edit_tweet_api_enabled":                                   true,
				"responsive_web_enhance_cards_enabled":                                    false,
				"responsive_web_graphql_exclude_directive_enabled":                        true,
				"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
				"responsive_web_graphql_timeline_navigation_enabled":                      true,
				"responsive_web_twitter_article_tweet_consumption_enabled":                true,
				"rweb_tipjar_consumption_enabled":                                         true,
				"rweb_video_timestamps_enabled":                                           true,
				"standardized_nudges_misinfo":                                             true,
				"tweet_awards_web_tipping_enabled":                                        false,
				"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
				"verified_phone_label_enabled":                                            false,
				"view_counts_everywhere_api_enabled":                                      true,
			},
			want: "%7B%22articles_preview_enabled%22%3Atrue%2C%22c9s_tweet_anatomy_moderator_badge_enabled%22%3Atrue%2C%22communities_web_enable_tweet_community_results_fetch%22%3Atrue%2C%22creator_subscriptions_quote_tweet_preview_enabled%22%3Afalse%2C%22creator_subscriptions_tweet_preview_api_enabled%22%3Atrue%2C%22freedom_of_speech_not_reach_fetch_enabled%22%3Atrue%2C%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%22%3Atrue%2C%22longform_notetweets_consumption_enabled%22%3Atrue%2C%22longform_notetweets_inline_media_enabled%22%3Atrue%2C%22longform_notetweets_rich_text_read_enabled%22%3Atrue%2C%22responsive_web_edit_tweet_api_enabled%22%3Atrue%2C%22responsive_web_enhance_cards_enabled%22%3Afalse%2C%22responsive_web_graphql_exclude_directive_enabled%22%3Atrue%2C%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%22%3Afalse%2C%22responsive_web_graphql_timeline_navigation_enabled%22%3Atrue%2C%22responsive_web_twitter_article_tweet_consumption_enabled%22%3Atrue%2C%22rweb_tipjar_consumption_enabled%22%3Atrue%2C%22rweb_video_timestamps_enabled%22%3Atrue%2C%22standardized_nudges_misinfo%22%3Atrue%2C%22tweet_awards_web_tipping_enabled%22%3Afalse%2C%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%22%3Atrue%2C%22verified_phone_label_enabled%22%3Afalse%2C%22view_counts_everywhere_api_enabled%22%3Atrue%7D",
		},
	}
	for _, tc := range tcs {

		got := featsToParams(tc.features)
		if got != tc.want {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}

func TestMagicURLBuilder(t *testing.T) {
	var tcs = []struct {
		variables Variables
		features  map[string]bool
		want      string
	}{
		{
			variables: Variables{
				UserId:                 "1234567890",
				Count:                  20,
				Cursor:                 "1234567891234567890|1234567891234567890",
				IncludePromotedContent: false,
			},
			features: map[string]bool{
				"rweb_tipjar_consumption_enabled":                                         true,
				"responsive_web_graphql_exclude_directive_enabled":                        true,
				"verified_phone_label_enabled":                                            false,
				"creator_subscriptions_tweet_preview_api_enabled":                         true,
				"responsive_web_graphql_timeline_navigation_enabled":                      true,
				"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
				"communities_web_enable_tweet_community_results_fetch":                    true,
				"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
				"articles_preview_enabled":                                                true,
				"responsive_web_edit_tweet_api_enabled":                                   true,
				"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
				"view_counts_everywhere_api_enabled":                                      true,
				"longform_notetweets_consumption_enabled":                                 true,
				"responsive_web_twitter_article_tweet_consumption_enabled":                true,
				"tweet_awards_web_tipping_enabled":                                        false,
				"creator_subscriptions_quote_tweet_preview_enabled":                       false,
				"freedom_of_speech_not_reach_fetch_enabled":                               true,
				"standardized_nudges_misinfo":                                             true,
				"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
				"rweb_video_timestamps_enabled":                                           true,
				"longform_notetweets_rich_text_read_enabled":                              true,
				"longform_notetweets_inline_media_enabled":                                true,
				"responsive_web_enhance_cards_enabled":                                    false,
			},
			want: "https://x.com/i/api/graphql/OSXFkKmGvfw_6pGgGtkWFg/Followers?variables=%7B%22userId%22%3A%221234567890%22%2C%22count%22%3A20%2C%22cursor%22%3A%221234567891234567890%7C1234567891234567890%22%2C%22includePromotedContent%22%3Afalse%7D&features=%7B%22articles_preview_enabled%22%3Atrue%2C%22c9s_tweet_anatomy_moderator_badge_enabled%22%3Atrue%2C%22communities_web_enable_tweet_community_results_fetch%22%3Atrue%2C%22creator_subscriptions_quote_tweet_preview_enabled%22%3Afalse%2C%22creator_subscriptions_tweet_preview_api_enabled%22%3Atrue%2C%22freedom_of_speech_not_reach_fetch_enabled%22%3Atrue%2C%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%22%3Atrue%2C%22longform_notetweets_consumption_enabled%22%3Atrue%2C%22longform_notetweets_inline_media_enabled%22%3Atrue%2C%22longform_notetweets_rich_text_read_enabled%22%3Atrue%2C%22responsive_web_edit_tweet_api_enabled%22%3Atrue%2C%22responsive_web_enhance_cards_enabled%22%3Afalse%2C%22responsive_web_graphql_exclude_directive_enabled%22%3Atrue%2C%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%22%3Afalse%2C%22responsive_web_graphql_timeline_navigation_enabled%22%3Atrue%2C%22responsive_web_twitter_article_tweet_consumption_enabled%22%3Atrue%2C%22rweb_tipjar_consumption_enabled%22%3Atrue%2C%22rweb_video_timestamps_enabled%22%3Atrue%2C%22standardized_nudges_misinfo%22%3Atrue%2C%22tweet_awards_web_tipping_enabled%22%3Afalse%2C%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%22%3Atrue%2C%22verified_phone_label_enabled%22%3Afalse%2C%22view_counts_everywhere_api_enabled%22%3Atrue%7D",
		},
	}

	baseurl := "https://x.com/i/api/graphql/OSXFkKmGvfw_6pGgGtkWFg/Followers"

	for _, tc := range tcs {
		got := magicUrlBuilder(baseurl, tc.variables, tc.features)
		if got != tc.want {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}
