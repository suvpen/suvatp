package builder

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/api/chat"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/suvpen/suvatp/api"
	"github.com/suvpen/suvatp/atperr"
	"github.com/suvpen/suvatp/util"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"
	"time"
)

type PostData struct {
	Text               string
	ImagePaths         []string
	QuoteUrl, EmbedUrl string
	CreatedAt          time.Time
}

type MessageData struct {
	ConvoId                   string
	Text                      string
	PostCid, PostUri, PostUrl string
}

func BuildPost(atpClient *api.ATPClient, postData PostData) (*bsky.FeedPost, error) {
	var createdAtStr string

	if postData.CreatedAt.IsZero() {
		createdAtStr = time.Now().Local().Format(time.RFC3339)
	} else {
		createdAtStr = postData.CreatedAt.Format(time.RFC3339)
	}

	post := &bsky.FeedPost{
		Text:      postData.Text,
		CreatedAt: createdAtStr,
	}

	if postData.QuoteUrl != "" {
		resp, err := atpClient.
			GetPost(util.GetHandleFromURL(postData.QuoteUrl), util.GetRecordKeyFromUrl(postData.QuoteUrl))
		if err != nil {
			return nil, fmt.Errorf("error building post: invalid QuoteUrl record: %w", err)
		}

		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}

		post.Embed.EmbedRecord = &bsky.EmbedRecord{
			Record: &atproto.RepoStrongRef{Cid: *resp.Cid, Uri: resp.Uri},
		}
	}

	if postData.EmbedUrl != "" {
		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}

		if post.Embed.EmbedExternal == nil {
			addLink(atpClient.Client, post, postData.EmbedUrl)
		}
	}

	injectedFacets, err := injectFacets(atpClient, post.Text)
	if err != nil {
		return nil, fmt.Errorf("error injecting facets: %w", err)
	}

	post.Facets = injectedFacets

	if len(postData.ImagePaths) > 0 {
		images, err := atpClient.UploadImages(postData.ImagePaths)
		if err != nil {
			return nil, err
		}

		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}

		post.Embed.EmbedImages = &bsky.EmbedImages{
			Images: images,
		}
	}

	post.Langs = []string{"id"}

	return post, nil
}

func BuildMessage(atpClient *api.ATPClient, msgData MessageData) (*chat.ConvoSendMessage_Input, error) {
	msgInput := &chat.ConvoSendMessage_Input{
		ConvoId: msgData.ConvoId,
		Message: &chat.ConvoDefs_MessageInput{
			Text: msgData.Text,
		},
	}

	if msgData.PostCid != "" && msgData.PostUri != "" {
		msgInput.Message.Embed = &chat.ConvoDefs_MessageInput_Embed{
			EmbedRecord: &bsky.EmbedRecord{
				Record: &atproto.RepoStrongRef{Cid: msgData.PostCid, Uri: msgData.PostUri},
			}}
	} else if msgData.PostUrl != "" {
		postRecord, err := atpClient.
			GetPost(util.GetHandleFromURL(msgData.PostUrl), util.GetRecordKeyFromUrl(msgData.PostUrl))
		if err == nil {
			msgInput.Message.Embed = &chat.ConvoDefs_MessageInput_Embed{
				EmbedRecord: &bsky.EmbedRecord{
					Record: &atproto.RepoStrongRef{Cid: *postRecord.Cid, Uri: postRecord.Uri},
				}}
		}
	}

	injectedFacets, err := injectFacets(atpClient, msgData.Text)
	if err != nil {
		return nil, fmt.Errorf("error building message: error injecting facets: %w", err)
	}

	msgInput.Message.Facets = injectedFacets

	return msgInput, nil
}

func BuildMessageBatch(atpClient *api.ATPClient, msgsData []MessageData) (*chat.ConvoSendMessageBatch_Input, error) {
	var msgItems []*chat.ConvoSendMessageBatch_BatchItem

	for _, msgData := range msgsData {
		msg, err := BuildMessage(atpClient, msgData)
		if err != nil {
			return nil, fmt.Errorf("error building message batch: %w", err)
		}

		msgItems = append(msgItems, &chat.ConvoSendMessageBatch_BatchItem{ConvoId: msg.ConvoId, Message: msg.Message})
	}

	return &chat.ConvoSendMessageBatch_Input{Items: msgItems}, nil
}

func injectFacets(atpClient *api.ATPClient, text string) ([]*bsky.RichtextFacet, error) {
	var facets []*bsky.RichtextFacet

	for _, ent := range util.ExtractLinksBytes(text) {
		facets = append(facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Link: &bsky.RichtextFacet_Link{
						Uri: ent.Text,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: ent.Start,
				ByteEnd:   ent.End,
			},
		})
	}

	for _, ent := range util.ExtractMentionsBytes(text) {
		if !strings.Contains(ent.Text, ".") {
			continue
		}

		profile, err := atpClient.GetProfile(ent.Text)
		if err != nil {
			if atperr.ErrorInvalidActorDidOrHandle(err) || atperr.ErrorProfileNotFound(err) {
				continue
			} else {
				return nil, err
			}
		}
		facets = append(facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Mention: &bsky.RichtextFacet_Mention{
						Did: profile.Did,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: ent.Start,
				ByteEnd:   ent.End,
			},
		})
	}

	for _, ent := range util.ExtractTagsBytes(text) {
		facets = append(facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
						Tag: ent.Text,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: ent.Start,
				ByteEnd:   ent.End,
			},
		})
	}

	return facets, nil
}

func addLink(xrpcc *xrpc.Client, post *bsky.FeedPost, link string) {
	res, _ := http.Get(link)
	if res != nil {
		defer res.Body.Close()

		br := bufio.NewReader(res.Body)
		var reader io.Reader = br

		data, err2 := br.Peek(1024)
		if err2 == nil {
			enc, name, _ := charset.DetermineEncoding(data, res.Header.Get("content-type"))
			if enc != nil {
				reader = enc.NewDecoder().Reader(br)
			} else if len(name) > 0 {
				enc := util.GetEncoding(name)
				if enc != nil {
					reader = enc.NewDecoder().Reader(br)
				}
			}
		}

		var title, description, imgUrl string

		doc, err := goquery.NewDocumentFromReader(reader)
		if err == nil {
			title = doc.Find(`title`).Text()
			description, _ = doc.Find(`meta[property="description"]`).Attr("content")
			imgUrl, _ = doc.Find(`meta[property="og:image"]`).Attr("content")
			if title == "" {
				title, _ = doc.Find(`meta[property="og:title"]`).Attr("content")
				if title == "" {
					title = link
				}
			}
			if description == "" {
				description, _ = doc.Find(`meta[property="og:description"]`).Attr("content")
				if description == "" {
					description = link
				}
			}
			post.Embed.EmbedExternal = &bsky.EmbedExternal{
				External: &bsky.EmbedExternal_External{
					Description: description,
					Title:       title,
					Uri:         link,
				},
			}
		} else {
			post.Embed.EmbedExternal = &bsky.EmbedExternal{
				External: &bsky.EmbedExternal_External{
					Uri: link,
				},
			}
		}
		if imgUrl != "" && post.Embed.EmbedExternal != nil {
			resp, err := http.Get(imgUrl)
			if err == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				b, err := io.ReadAll(resp.Body)
				if err == nil {
					blobResp, err := atproto.RepoUploadBlob(context.TODO(), xrpcc, bytes.NewReader(b))
					if err == nil {
						post.Embed.EmbedExternal.External.Thumb = &lexutil.LexBlob{
							Ref:      blobResp.Blob.Ref,
							MimeType: http.DetectContentType(b),
							Size:     blobResp.Blob.Size,
						}
					}
				}
			}
		}
	}
}
