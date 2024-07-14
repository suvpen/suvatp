package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/ozone"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/suvpen/suvatp/atperr"
	"time"
)

func (atpClient *ATPClient) SearchRepos(q, cursor string, limit int64) (*ozone.ModerationSearchRepos_Output, error) {
	resp, err := ozone.ModerationSearchRepos(context.TODO(), atpClient.LabelerClient, cursor, limit, q, "")
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SearchRepos(q, cursor, limit)
			} else {
				return nil, fmt.Errorf("error while searching repos of %s: %w", q, err)
			}
		} else {
			return nil, fmt.Errorf("error while searching repos of %s: %w", q, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.QueryLabel(cursor, limit)
			} else {
				return nil, fmt.Errorf("error querying label events: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error querying label events: %w", err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}

func (atpClient *ATPClient) QueryOpenReports(cursor string, limit int64) (*ozone.ModerationQueryStatuses_Output, error) {
	//TODO: uncomment if fixed
	//resp, err := ozone.ModerationQueryStatuses(
	//	context.TODO(), atpClient.LabelerClient,
	//	false, "", cursor, nil, nil,
	//	true, "", limit, false, "", "",
	//	"tools.ozone.moderation.defs#reviewOpen", "", "", "desc", "lastReportedAt",
	//	"", nil, false)
	//if err != nil {
	//	if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
	//		if atpClient.RetryCount != atpClient.Config.Retries {
	//			atpClient.RetryCount++
	//			time.Sleep(time.Second * 3)
	//			return atpClient.QueryOpenReports(cursor, limit)
	//		} else {
	//			return nil, fmt.Errorf("error querying open reports: %w", err)
	//		}
	//	} else {
	//		return nil, fmt.Errorf("error querying open reports: %w", err)
	//	}
	//}

	params := map[string]interface{}{
		"cursor":        cursor,
		"includeMuted":  true,
		"limit":         limit,
		"onlyMuted":     false,
		"reviewState":   "tools.ozone.moderation.defs#reviewOpen",
		"sortDirection": "desc",
		"sortField":     "lastReportedAt",
	}

	var resp *ozone.ModerationQueryStatuses_Output

	err := atpClient.LabelerClient.Do(
		context.TODO(), xrpc.Query, "", "tools.ozone.moderation.queryStatuses", params, nil, resp)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.QueryOpenReports(cursor, limit)
			} else {
				return nil, fmt.Errorf("error querying open reports: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error querying open reports: %w", err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.QueryEventDetail(subject)
			} else {
				return nil, fmt.Errorf("error querying label events: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error querying label events: %w", err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.LabelAccount(adminDid, targetDid, label)
			} else {
				return nil, fmt.Errorf("error labeling %s: %w", targetDid, err)
			}
		} else {
			return nil, fmt.Errorf("error labeling %s: %w", targetDid, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.LabelPost(adminDid, cid, uri, label)
			} else {
				return nil, fmt.Errorf("error labeling post %s: %w", uri, err)
			}
		} else {
			return nil, fmt.Errorf("error labeling post %s: %w", uri, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.NegateAccountLabel(adminDid, targetDid, label)
			} else {
				return nil, fmt.Errorf("error unlabeling %s: %w", targetDid, err)
			}
		} else {
			return nil, fmt.Errorf("error unlabeling %s: %w", targetDid, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.NegatePostLabel(adminDid, cid, uri, label)
			} else {
				return nil, fmt.Errorf("error unlabeling post %s: %w", uri, err)
			}
		} else {
			return nil, fmt.Errorf("error unlabeling post %s: %w", uri, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.AcknowledgeAccountRecord(adminDid, targetDid)
			} else {
				return nil, fmt.Errorf("error acknowledging %s account record: %w", targetDid, err)
			}
		} else {
			return nil, fmt.Errorf("error acknowledging %s account record: %w", targetDid, err)
		}
	}

	atpClient.RetryCount = 0

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
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.AcknowledgePostRecord(adminDid, cid, uri)
			} else {
				return nil, fmt.Errorf("error acknowledging %s post record: %w", uri, err)
			}
		} else {
			return nil, fmt.Errorf("error acknowledging %s post record: %w", uri, err)
		}
	}

	atpClient.RetryCount = 0

	return resp, nil
}
