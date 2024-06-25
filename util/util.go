package util

import (
	"fmt"
	"strings"
)

type FollowRecord struct {
	Did       string
	Schema    string
	RecordKey string
}

func DecodeGraphRecord(uri string) (*FollowRecord, error) {
	parts := strings.Split(uri, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("error: invalid post uri: %q", uri)
	}

	return &FollowRecord{
		Did:       parts[len(parts)-3],
		Schema:    parts[len(parts)-2],
		RecordKey: parts[len(parts)-1],
	}, nil
}

func CreateBskyProfileURL(handle string) string {
	return fmt.Sprintf("https://bsky.app/profile/%s", handle)
}

func CreateBskyPostURL(handle, rkey string) string {
	if handle == "" || rkey == "" {
		return "-"
	}

	return fmt.Sprintf("https://bsky.app/profile/%s/post/%s", handle, rkey)
}

func GetHandleFromURL(atpUrl string) string {
	atpUrl = strings.Replace(atpUrl, "@", "", -1)
	atpUrl = strings.Split(strings.ToLower(atpUrl), "?")[0]
	linkParts := strings.Split(strings.ToLower(atpUrl), "/")

	if len(linkParts) >= 5 {
		return linkParts[4]
	} else if len(linkParts) == 3 {
		return linkParts[2]
	} else {
		return ""
	}
}

func GetRecordKeyFromUrl(url string) string {
	uriSplit := strings.Split(url, "/")
	return uriSplit[len(uriSplit)-1]
}
