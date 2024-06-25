package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/ozone"
)

func (atpClient *ATPClient) LabelAccount(adminDid, targetDid, label string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventLabel: &ozone.ModerationDefs_ModEventLabel{
				CreateLabelVals: []string{label},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			AdminDefs_RepoRef: &atproto.AdminDefs_RepoRef{
				Did: targetDid,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.Client, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error labeling %s: %w", targetDid, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) NegateAccountLabel(adminDid, targetDid, label string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventLabel: &ozone.ModerationDefs_ModEventLabel{
				NegateLabelVals: []string{label},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			AdminDefs_RepoRef: &atproto.AdminDefs_RepoRef{
				Did: targetDid,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.Client, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error labeling %s: %w", targetDid, err)
	}

	return resp, nil
}