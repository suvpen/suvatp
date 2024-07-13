package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/suvpen/suvatp/atperr"
	"net/http"
	"os"
	"time"
)

func (atpClient *ATPClient) GetPost(didOrHandle, rKey string) (*atproto.RepoGetRecord_Output, error) {
	resp, err := atproto.RepoGetRecord(
		context.TODO(), atpClient.Client, "", atpClient.Config.PostsCollection, didOrHandle, rKey)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetPost(didOrHandle, rKey)
			} else {
				return nil, fmt.Errorf("error getting post record: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting post record: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) GetPostThread(
	didOrHandle, rKey string, depth, parentHeight int64) (*bsky.FeedGetPostThread_Output, error) {

	postRecord, err := atpClient.GetPost(didOrHandle, rKey)
	if err != nil {
		return nil, fmt.Errorf("error getting post thread: %w", err)
	}

	resp, err := bsky.FeedGetPostThread(
		context.TODO(), atpClient.Client, depth, parentHeight, postRecord.Uri)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetPostThread(didOrHandle, rKey, depth, parentHeight)
			} else {
				return nil, fmt.Errorf("error getting post thread: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting post thread: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) GetAuthorFeed(did, cursor, filter string, limit int64) (*bsky.FeedGetAuthorFeed_Output, error) {
	filters := []string{FilterPostsWithReplies, FilterPostsNoReplies, FilterPostsWithMedia, FilterPostsAndAuthorThreads}
	var found bool
	for _, fil := range filters {
		if fil == filter {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("error getting %s feed: invalid filter", did)
	}

	resp, err := bsky.FeedGetAuthorFeed(context.TODO(), atpClient.Client, did, cursor, filter, limit)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetAuthorFeed(did, cursor, filter, limit)
			} else {
				return nil, fmt.Errorf("error getting %s feed: %w", did, err)
			}
		} else {
			return nil, fmt.Errorf("error getting %s feed: %w", did, err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) GetRepostedBy(didOrHandle, rKey, cursor string) (*bsky.FeedGetRepostedBy_Output, error) {
	postRecord, err := atpClient.GetPost(didOrHandle, rKey)
	if err != nil {
		return nil, fmt.Errorf("error getting repostedby: %w", err)
	}

	resp, err := bsky.FeedGetRepostedBy(
		context.TODO(), atpClient.Client, *postRecord.Cid, cursor, 100, postRecord.Uri)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetRepostedBy(didOrHandle, rKey, cursor)
			} else {
				return nil, fmt.Errorf("error getting repostedby: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting repostedby: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) GetLikes(didOrHandle, rKey, cursor string) (*bsky.FeedGetLikes_Output, error) {
	postRecord, err := atpClient.GetPost(didOrHandle, rKey)
	if err != nil {
		return nil, fmt.Errorf("error getting likes: %w", err)
	}

	resp, err := bsky.FeedGetLikes(
		context.TODO(), atpClient.Client, *postRecord.Cid, cursor, 100, postRecord.Uri)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetLikes(didOrHandle, rKey, cursor)
			} else {
				return nil, fmt.Errorf("error getting %s likes: %w", postRecord.Uri, err)
			}
		} else {
			return nil, fmt.Errorf("error getting %s likes: %w", postRecord.Uri, err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) SearchPost(q, cursor string, limit int64) (*bsky.FeedSearchPosts_Output, error) {
	resp, err := bsky.FeedSearchPosts(
		context.TODO(), atpClient.Client, "", cursor, "", "",
		limit, "", q, "", "",
		nil, "", "")
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SearchPost(q, cursor, limit)
			} else {
				return nil, fmt.Errorf("error searching post: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error searching post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) Post(post *bsky.FeedPost) (*atproto.RepoCreateRecord_Output, error) {
	resp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.PostsCollection,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Post(post)
			} else {
				return nil, fmt.Errorf("error creating post: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error creating post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) ReplyPost(cid, uri string, post *bsky.FeedPost) (*atproto.RepoCreateRecord_Output, error) {
	post.Reply = &bsky.FeedPost_ReplyRef{
		Root:   &atproto.RepoStrongRef{Cid: cid, Uri: uri},
		Parent: &atproto.RepoStrongRef{Cid: cid, Uri: uri},
	}

	resp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.PostsCollection,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.ReplyPost(cid, uri, post)
			} else {
				return nil, fmt.Errorf("error replying post: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error replying post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) DeletePost(rKey string) error {
	err := atproto.RepoDeleteRecord(context.TODO(), atpClient.Client, &atproto.RepoDeleteRecord_Input{
		Collection: atpClient.Config.PostsCollection,
		Repo:       atpClient.Client.Auth.Did,
		Rkey:       rKey,
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.DeletePost(rKey)
			} else {
				return fmt.Errorf("error deleting post: %w", err)
			}
		} else {
			return fmt.Errorf("error deleting post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) Repost(didOrHandle, rKey string) (*atproto.RepoCreateRecord_Output, error) {
	postRecord, err := atpClient.GetPost(didOrHandle, rKey)
	if err != nil {
		return nil, fmt.Errorf("error reposting post: %w", err)
	}

	repostResp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.RepostsCollection,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.FeedRepost{
				CreatedAt: time.Now().Local().Format(time.RFC3339),
				Subject: &atproto.RepoStrongRef{
					Uri: postRecord.Uri,
					Cid: *postRecord.Cid,
				},
			},
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Repost(didOrHandle, rKey)
			} else {
				return nil, fmt.Errorf("error reposting post: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error reposting post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return repostResp, nil
}

func (atpClient *ATPClient) UndoRepost(rKey string) error {
	err := atproto.RepoDeleteRecord(context.TODO(), atpClient.Client, &atproto.RepoDeleteRecord_Input{
		Collection: atpClient.Config.RepostsCollection,
		Repo:       atpClient.Client.Auth.Did,
		Rkey:       rKey,
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.UndoRepost(rKey)
			} else {
				return fmt.Errorf("error undoing repost: %w", err)
			}
		} else {
			return fmt.Errorf("error undoing repost: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) Like(didOrHandle, rKey string) (*atproto.RepoCreateRecord_Output, error) {
	postRecord, err := atpClient.GetPost(didOrHandle, rKey)
	if err != nil {
		return nil, fmt.Errorf("error liking post: %w", err)
	}

	repostResp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.LikesCollection,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.FeedLike{
				CreatedAt: time.Now().Local().Format(time.RFC3339),
				Subject: &atproto.RepoStrongRef{
					Cid: *postRecord.Cid,
					Uri: postRecord.Uri,
				},
			},
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Like(didOrHandle, rKey)
			} else {
				return nil, fmt.Errorf("error liking post: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error liking post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return repostResp, nil
}

func (atpClient *ATPClient) Unlike(rKey string) error {
	err := atproto.RepoDeleteRecord(context.TODO(), atpClient.Client, &atproto.RepoDeleteRecord_Input{
		Collection: atpClient.Config.LikesCollection,
		Repo:       atpClient.Client.Auth.Did,
		Rkey:       rKey,
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Unlike(rKey)
			} else {
				return fmt.Errorf("error unliking post: %w", err)
			}
		} else {
			return fmt.Errorf("error unliking post: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) UploadImages(imagePaths []string) ([]*bsky.EmbedImages_Image, error) {
	if len(imagePaths) == 0 {
		return nil, nil
	}

	var images []*bsky.EmbedImages_Image
	for _, imgPath := range imagePaths {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			return nil, fmt.Errorf("error uploading image: cannot read image file: %w", err)
		}

		resp, err := atproto.RepoUploadBlob(context.TODO(), atpClient.Client, bytes.NewReader(imgData))
		if err != nil {
			if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
				time.Sleep(time.Second * 3)

				resp, err = atproto.RepoUploadBlob(context.TODO(), atpClient.Client, bytes.NewReader(imgData))
				if err != nil {
					return nil, fmt.Errorf("error uploading image: cannot upload image: %w", err)
				}
			} else {
				return nil, fmt.Errorf("error uploading image: cannot upload image: %w", err)
			}
		}

		images = append(images, &bsky.EmbedImages_Image{
			Image: &lexutil.LexBlob{
				Ref:      resp.Blob.Ref,
				MimeType: http.DetectContentType(imgData),
				Size:     resp.Blob.Size,
			},
		})
	}

	atpClient.RetryCount = 0

	return images, nil
}
