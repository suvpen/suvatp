package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/suvpen/suvatp/util"
	"strings"
	"time"
)

func (atpClient *ATPClient) GetProfile(didOrHandle string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
	profile, err := bsky.ActorGetProfile(context.TODO(), atpClient.Client, didOrHandle)
	if err != nil {
		return nil, fmt.Errorf("error get %s profile: %w", didOrHandle, err)
	}

	return profile, nil
}

func (atpClient *ATPClient) SearchActors(q, cursor string, limit int64) (*bsky.ActorSearchActors_Output, error) {
	profile, err := bsky.ActorSearchActors(context.TODO(), atpClient.Client, cursor, limit, q, "")
	if err != nil {
		return nil, fmt.Errorf("error searching actors with q=%s: %w", q, err)
	}

	return profile, nil
}

func (atpClient *ATPClient) GetFollows(cursor string) (*bsky.GraphGetFollows_Output, error) {
	follows, err := bsky.GraphGetFollows(context.TODO(), atpClient.Client, atpClient.Client.Auth.Did, cursor, 100)
	if err != nil {
		return nil, fmt.Errorf("error getting follows: %w", err)
	}

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
		return nil, fmt.Errorf("error following DID %s: %w", did, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) FollowHandle(handle string) (*atproto.RepoCreateRecord_Output, error) {
	profile, err := atpClient.GetProfile(handle)
	if err != nil {
		return nil, fmt.Errorf("error following handle %s: %w", handle, err)
	}

	return atpClient.FollowDid(profile.Did)
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
		return fmt.Errorf("error unfollowing DID %s: %w", didOrHandle, err)
	}

	return nil
}

func (atpClient *ATPClient) MuteDid(did string) error {
	err := bsky.GraphMuteActor(context.TODO(), atpClient.Client, &bsky.GraphMuteActor_Input{Actor: did})
	if err != nil {
		return fmt.Errorf("error muting DID %s: %w", did, err)
	}

	return nil
}

func (atpClient *ATPClient) MuteHandle(handle string) error {
	profile, err := atpClient.GetProfile(handle)
	if err != nil {
		return fmt.Errorf("error muting handle %s: %w", handle, err)
	}

	return atpClient.MuteDid(profile.Did)
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
		return nil, fmt.Errorf("error blocking DID %s: %w", did, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) BlockHandle(handle string) (*atproto.RepoCreateRecord_Output, error) {
	profile, err := atpClient.GetProfile(handle)
	if err != nil {
		return nil, fmt.Errorf("error blocking handle %s: %w", handle, err)
	}

	return atpClient.BlockDid(profile.Did)
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
		return fmt.Errorf("error unblocking DID %s: %w", didOrHandle, err)
	}

	return nil
}
