package util

import (
	"regexp"
	"strings"
)

const (
	urlPattern     = `https?://[-A-Za-z0-9+&@#\/%?=~_|!:,.;\(\)]+`
	mentionPattern = `@[a-zA-Z0-9.]+`
	tagPattern     = `\B#\S+`
)

var (
	urlRegex     = regexp.MustCompile(urlPattern)
	mentionRegex = regexp.MustCompile(mentionPattern)
	tagRegex     = regexp.MustCompile(tagPattern)
)

type FacetEntity struct {
	Start int64
	End   int64
	Text  string
}

type URIComponent struct {
	Handle string
	RKey   string
}

func ExtractLinksBytes(text string) []FacetEntity {
	var result []FacetEntity
	matches := urlRegex.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		result = append(result, FacetEntity{
			Text:  text[m[0]:m[1]],
			Start: int64(len(text[0:m[0]])),
			End:   int64(len(text[0:m[1]]))},
		)
	}
	return result
}

func ExtractURIComponent(uri string) URIComponent {
	textParts := strings.Split(uri, "/")
	if len(textParts) != 7 {
		return URIComponent{}
	}

	return URIComponent{
		Handle: textParts[4],
		RKey:   textParts[6],
	}
}

func ExtractMentionsBytes(text string) []FacetEntity {
	var result []FacetEntity
	matches := mentionRegex.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		result = append(result, FacetEntity{
			Text:  strings.TrimPrefix(text[m[0]:m[1]], "@"),
			Start: int64(len(text[0:m[0]])),
			End:   int64(len(text[0:m[1]]))},
		)
	}
	return result
}

func ExtractTagsBytes(text string) []FacetEntity {
	var result []FacetEntity
	matches := tagRegex.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		result = append(result, FacetEntity{
			Text:  strings.TrimPrefix(text[m[0]:m[1]], "#"),
			Start: int64(len(text[0:m[0]])),
			End:   int64(len(text[0:m[1]]))},
		)
	}
	return result
}
