package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/suvpen/suvatp/atperr"
	"github.com/suvpen/suvatp/util"
	"strings"
	"time"
)

func (atpClient *ATPClient) GetPreferences() (*bsky.ActorGetPreferences_Output, error) {
	resp, err := bsky.ActorGetPreferences(context.TODO(), atpClient.Client)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetPreferences()
			} else {
				return nil, fmt.Errorf("error getting preferences: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting preferences: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) SubscribeLabeler(did string) error {
	lastPrefs, err := atpClient.GetPreferences()
	if err != nil {
		return fmt.Errorf("error subscribing preferences: %w", err)
	}

	input := &bsky.ActorPutPreferences_Input{
		Preferences: lastPrefs.Preferences,
	}

	input.Preferences = append(input.Preferences, bsky.ActorDefs_Preferences_Elem{
		ActorDefs_LabelersPref: &bsky.ActorDefs_LabelersPref{
			Labelers: []*bsky.ActorDefs_LabelerPrefItem{
				{
					Did: did,
				},
			},
		},
	})

	err = bsky.ActorPutPreferences(context.TODO(), atpClient.Client, input)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SubscribeLabeler(did)
			} else {
				return fmt.Errorf("error subscribing preferences: %w", err)
			}
		} else {
			return fmt.Errorf("error subscribing preferences: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) LikeLabeler(cid, did string) (*atproto.RepoCreateRecord_Output, error) {
	repostResp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.LikesCollection,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.FeedLike{
				CreatedAt: time.Now().Local().Format(time.RFC3339),
				Subject: &atproto.RepoStrongRef{
					Cid: cid,
					Uri: fmt.Sprintf("at://%s/%s/self", did, atpClient.Config.LabelerService),
				},
			},
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.LikeLabeler(cid, did)
			} else {
				return nil, fmt.Errorf("error liking labeler: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error liking labeler: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return repostResp, nil
}

func (atpClient *ATPClient) ResolveHandle(handle string) (string, error) {
	resp, err := atproto.IdentityResolveHandle(context.TODO(), atpClient.Client, handle)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.ResolveHandle(handle)
			} else {
				return "", fmt.Errorf("error resolving %s handle: %w", handle, err)
			}
		} else {
			return "", fmt.Errorf("error resolving %s handle: %w", handle, err)
		}
	}

	atpClient.RetryCount = 0

	return resp.Did, nil
}

func (atpClient *ATPClient) GetProfile(didOrHandle string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
	profile, err := bsky.ActorGetProfile(context.TODO(), atpClient.Client, didOrHandle)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetProfile(didOrHandle)
			} else {
				return nil, fmt.Errorf("error getting %s profile: %w", didOrHandle, err)
			}
		} else {
			return nil, fmt.Errorf("error getting %s profile: %w", didOrHandle, err)
		}
	}

	atpClient.RetryCount = 0

	return profile, nil
}

func (atpClient *ATPClient) SearchActors(q, cursor string, limit int64) (*bsky.ActorSearchActors_Output, error) {
	profile, err := bsky.ActorSearchActors(context.TODO(), atpClient.Client, cursor, limit, q, "")
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SearchActors(q, cursor, limit)
			} else {
				return nil, fmt.Errorf("error searching actors with q=%s: %w", q, err)
			}
		} else {
			return nil, fmt.Errorf("error searching actors with q=%s: %w", q, err)
		}
	}

	atpClient.RetryCount = 0

	return profile, nil
}

func (atpClient *ATPClient) GetFollows(cursor string) (*bsky.GraphGetFollows_Output, error) {
	follows, err := bsky.GraphGetFollows(context.TODO(), atpClient.Client, atpClient.Client.Auth.Did, cursor, 100)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetFollows(cursor)
			} else {
				return nil, fmt.Errorf("error getting follows: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting follows: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return follows, nil
}

func (atpClient *ATPClient) FollowDid(did string) (*atproto.RepoCreateRecord_Output, error) {
	if !strings.Contains(did, "did:plc:") {
		return nil, fmt.Errorf("error following DID %s: DID must contain 'did:plc:'", did)
	}

	resp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.GraphFollowLexicon,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.GraphFollow{
				LexiconTypeID: atpClient.Config.GraphFollowLexicon,
				CreatedAt:     time.Now().Local().Format(time.RFC3339),
				Subject:       did,
			},
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.FollowDid(did)
			} else {
				return nil, fmt.Errorf("error following DID %s: %w", did, err)
			}
		} else {
			return nil, fmt.Errorf("error following DID %s: %w", did, err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) FollowHandle(handle string) (*atproto.RepoCreateRecord_Output, error) {
	did, err := atpClient.ResolveHandle(handle)
	if err != nil {
		return nil, fmt.Errorf("error following handle %s: %w", handle, err)
	}

	return atpClient.FollowDid(did)
}

func (atpClient *ATPClient) Unfollow(didOrHandle string) error {
	profile, err := atpClient.GetProfile(didOrHandle)
	if err != nil {
		return fmt.Errorf("error unfollowing %s: %w", didOrHandle, err)
	}

	if profile.Viewer.Following == nil {
		return nil
	}

	folRecord, err := util.DecodeGraphRecord(*profile.Viewer.Following)
	if err != nil {
		return fmt.Errorf("error unfollowing DID %s: %w", didOrHandle, err)
	}

	err = atproto.RepoDeleteRecord(context.TODO(), atpClient.Client, &atproto.RepoDeleteRecord_Input{
		Repo:       atpClient.Client.Auth.Did,
		Collection: folRecord.Schema,
		Rkey:       folRecord.RecordKey,
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Unfollow(didOrHandle)
			} else {
				return fmt.Errorf("error unfollowing DID %s: %w", didOrHandle, err)
			}
		} else {
			return fmt.Errorf("error unfollowing DID %s: %w", didOrHandle, err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) MuteDid(did string) error {
	err := bsky.GraphMuteActor(context.TODO(), atpClient.Client, &bsky.GraphMuteActor_Input{Actor: did})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Unfollow(did)
			} else {
				return fmt.Errorf("error muting DID %s: %w", did, err)
			}
		} else {
			return fmt.Errorf("error muting DID %s: %w", did, err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}

func (atpClient *ATPClient) MuteHandle(handle string) error {
	did, err := atpClient.ResolveHandle(handle)
	if err != nil {
		return fmt.Errorf("error muting handle %s: %w", handle, err)
	}

	return atpClient.MuteDid(did)
}

func (atpClient *ATPClient) BlockDid(did string) (*atproto.RepoCreateRecord_Output, error) {
	resp, err := atproto.RepoCreateRecord(context.TODO(), atpClient.Client, &atproto.RepoCreateRecord_Input{
		Collection: atpClient.Config.GraphBlockLexicon,
		Repo:       atpClient.Client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.GraphBlock{
				LexiconTypeID: atpClient.Config.GraphBlockLexicon,
				CreatedAt:     time.Now().Local().Format(time.RFC3339),
				Subject:       did,
			},
		},
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.BlockDid(did)
			} else {
				return nil, fmt.Errorf("error blocking DID %s: %w", did, err)
			}
		} else {
			return nil, fmt.Errorf("error blocking DID %s: %w", did, err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) BlockHandle(handle string) (*atproto.RepoCreateRecord_Output, error) {
	did, err := atpClient.ResolveHandle(handle)
	if err != nil {
		return nil, fmt.Errorf("error blocking handle %s: %w", handle, err)
	}

	return atpClient.BlockDid(did)
}

func (atpClient *ATPClient) Unblock(didOrHandle string) error {
	profile, err := atpClient.GetProfile(didOrHandle)
	if err != nil {
		return fmt.Errorf("error unblocking %s: %w", didOrHandle, err)
	}

	if profile.Viewer.Blocking == nil {
		return nil
	}

	blockRecord, err := util.DecodeGraphRecord(*profile.Viewer.Blocking)
	if err != nil {
		return fmt.Errorf("error unblocking DID %s: %w", didOrHandle, err)
	}

	err = atproto.RepoDeleteRecord(context.TODO(), atpClient.Client, &atproto.RepoDeleteRecord_Input{
		Repo:       atpClient.Client.Auth.Did,
		Collection: blockRecord.Schema,
		Rkey:       blockRecord.RecordKey,
	})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.Unblock(didOrHandle)
			} else {
				return fmt.Errorf("error unblocking %s: %w", didOrHandle, err)
			}
		} else {
			return fmt.Errorf("error unblocking %s: %w", didOrHandle, err)
		}
	}

	atpClient.RetryCount = 0

	return nil
}
