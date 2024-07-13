package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/chat"
	"github.com/suvpen/suvatp/atperr"
	"time"
)

func (atpClient *ATPClient) GetConvoForMembers(targetDid string) (*chat.ConvoGetConvoForMembers_Output, error) {
	resp, err := chat.ConvoGetConvoForMembers(context.TODO(), atpClient.PdsClient, []string{targetDid})
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetConvoForMembers(targetDid)
			} else {
				return nil, fmt.Errorf("error getting convo for members %s : %w", targetDid, err)
			}
		} else {
			return nil, fmt.Errorf("error getting convo for members %s : %w", targetDid, err)
		}
	}

	return resp, nil
}

func (atpClient *ATPClient) ListConvos(cursor string, limit int64) (*chat.ConvoListConvos_Output, error) {
	resp, err := chat.ConvoListConvos(
		context.TODO(), atpClient.PdsClient, cursor, limit)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.ListConvos(cursor, limit)
			} else {
				return nil, fmt.Errorf("error getting chat list: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting chat list: %w", err)
		}
	}

	return resp, nil
}

func (atpClient *ATPClient) GetLog(cursor string) (*chat.ConvoGetLog_Output, error) {
	resp, err := chat.ConvoGetLog(context.TODO(), atpClient.PdsClient, cursor)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.GetLog(cursor)
			} else {
				return nil, fmt.Errorf("error getting chat log: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting chat log: %w", err)
		}
	}

	return resp, nil
}

func (atpClient *ATPClient) SendMessage(msgInput *chat.ConvoSendMessage_Input) (*chat.ConvoDefs_MessageView, error) {
	resp, err := chat.ConvoSendMessage(context.TODO(), atpClient.PdsClient, msgInput)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SendMessage(msgInput)
			} else {
				return nil, fmt.Errorf("error sending message: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error sending message: %w", err)
		}
	}

	return resp, nil
}

func (atpClient *ATPClient) SendMessageBatch(msgBatchInput *chat.ConvoSendMessageBatch_Input) (*chat.ConvoSendMessageBatch_Output, error) {
	resp, err := chat.ConvoSendMessageBatch(context.TODO(), atpClient.PdsClient, msgBatchInput)
	if err != nil {
		if atperr.IsUpstreamFailureError(err) || atperr.IsUpstreamTimeoutError(err) || atperr.IsInternalServerError(err) {
			if atpClient.RetryCount != atpClient.Config.Retries {
				atpClient.RetryCount++
				time.Sleep(time.Second * 3)
				return atpClient.SendMessageBatch(msgBatchInput)
			} else {
				return nil, fmt.Errorf("error sending message batch: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error sending message batch: %w", err)
		}
	}

	return resp, nil
}
