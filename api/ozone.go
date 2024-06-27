package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/ozone"
)

func (atpClient *ATPClient) SearchRepos(q, cursor string, limit int64) (*ozone.ModerationSearchRepos_Output, error) {
	resp, err := ozone.ModerationSearchRepos(context.TODO(), atpClient.LabelerClient, cursor, limit, q, "")
	if err != nil {
		return nil, fmt.Errorf("error while searching repos of %s: %w", q, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) QueryLabel(cursor string, limit int64) (*ozone.ModerationQueryEvents_Output, error) {
	resp, err := ozone.ModerationQueryEvents(
		context.TODO(), atpClient.LabelerClient,
		nil, nil, "", "", "",
		"", cursor, false, false, limit,
		nil, nil, nil, "", "",
		[]string{"tools.ozone.moderation.defs#modEventLabel"})
	if err != nil {
		return nil, fmt.Errorf("error querying label events: %w", err)
	}

	return resp, nil
}

func (atpClient *ATPClient) LabelAccount(adminDid, targetDid, label string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventLabel: &ozone.ModerationDefs_ModEventLabel{
				CreateLabelVals: []string{label},
				NegateLabelVals: []string{},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			AdminDefs_RepoRef: &atproto.AdminDefs_RepoRef{
				Did: targetDid,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
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
				CreateLabelVals: []string{},
				NegateLabelVals: []string{label},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			AdminDefs_RepoRef: &atproto.AdminDefs_RepoRef{
				Did: targetDid,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error labeling %s: %w", targetDid, err)
	}

	return resp, nil
}
