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

func (atpClient *ATPClient) QueryOpenReports(cursor string, limit int64) (*ozone.ModerationQueryStatuses_Output, error) {
	resp, err := ozone.ModerationQueryStatuses(
		context.TODO(), atpClient.LabelerClient,
		false, "", cursor, nil, nil,
		true, "", limit, "", "",
		"tools.ozone.moderation.defs#reviewOpen", "", "", "desc", "lastReportedAt",
		"", nil, false)
	if err != nil {
		return nil, fmt.Errorf("error querying open reports: %w", err)
	}

	return resp, nil
}

func (atpClient *ATPClient) QueryEventDetail(subject string) (*ozone.ModerationQueryEvents_Output, error) {
	resp, err := ozone.ModerationQueryEvents(
		context.TODO(), atpClient.LabelerClient,
		nil, nil, "", "", "",
		"", "", false, false, 2,
		nil, nil, nil, "", subject,
		nil)
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

func (atpClient *ATPClient) LabelPost(adminDid, cid, uri, label string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventLabel: &ozone.ModerationDefs_ModEventLabel{
				CreateLabelVals: []string{label},
				NegateLabelVals: []string{},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			RepoStrongRef: &atproto.RepoStrongRef{
				Cid: cid,
				Uri: uri,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error labeling post %s: %w", uri, err)
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

func (atpClient *ATPClient) NegatePostLabel(adminDid, cid, uri, label string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventLabel: &ozone.ModerationDefs_ModEventLabel{
				CreateLabelVals: []string{},
				NegateLabelVals: []string{label},
			},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			RepoStrongRef: &atproto.RepoStrongRef{
				Cid: cid,
				Uri: uri,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error labeling post %s: %w", uri, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) AcknowledgeAccountRecord(adminDid, targetDid string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventAcknowledge: &ozone.ModerationDefs_ModEventAcknowledge{},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			AdminDefs_RepoRef: &atproto.AdminDefs_RepoRef{
				Did: targetDid,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error acknowledging %s account record: %w", targetDid, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) AcknowledgePostRecord(adminDid, cid, uri string) (*ozone.ModerationDefs_ModEventView, error) {
	eventInput := &ozone.ModerationEmitEvent_Input{
		CreatedBy: adminDid,
		Event: &ozone.ModerationEmitEvent_Input_Event{
			ModerationDefs_ModEventAcknowledge: &ozone.ModerationDefs_ModEventAcknowledge{},
		},
		Subject: &ozone.ModerationEmitEvent_Input_Subject{
			RepoStrongRef: &atproto.RepoStrongRef{
				Cid: cid,
				Uri: uri,
			},
		},
	}

	resp, err := ozone.ModerationEmitEvent(context.TODO(), atpClient.LabelerClient, eventInput)
	if err != nil {
		return nil, fmt.Errorf("error acknowledging %s post record: %w", uri, err)
	}

	return resp, nil
}
