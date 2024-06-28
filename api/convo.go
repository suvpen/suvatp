package api

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/chat"
)

type MessageInput struct {
	Text string `json:"text"`
}

func (atpClient *ATPClient) GetConvoForMembers(targetDid string) (*chat.ConvoGetConvoForMembers_Output, error) {
	resp, err := chat.ConvoGetConvoForMembers(context.TODO(), atpClient.PdsClient, []string{targetDid})
	if err != nil {
		return nil, fmt.Errorf("error getting convo for members %s : %w", targetDid, err)
	}

	return resp, nil
}

func (atpClient *ATPClient) ListConvos(cursor string, limit int64) (*chat.ConvoListConvos_Output, error) {
	resp, err := chat.ConvoListConvos(
		context.TODO(), atpClient.PdsClient, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting chat list: %w", err)
	}

	return resp, nil
}

func (atpClient *ATPClient) GetLog(cursor string) (*chat.ConvoGetLog_Output, error) {
	resp, err := chat.ConvoGetLog(context.TODO(), atpClient.PdsClient, cursor)
	if err != nil {
		return nil, fmt.Errorf("error getting chat log: %w", err)
	}

	return resp, nil
}

func (atpClient *ATPClient) SendMessage(msgInput *chat.ConvoSendMessage_Input) (*chat.ConvoDefs_MessageView, error) {
	resp, err := chat.ConvoSendMessage(context.TODO(), atpClient.PdsClient, msgInput)
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return resp, nil
}

func (atpClient *ATPClient) SendMessageBatch(msgInputs *chat.ConvoSendMessageBatch_Input) (*chat.ConvoSendMessageBatch_Output, error) {
	resp, err := chat.ConvoSendMessageBatch(context.TODO(), atpClient.PdsClient, msgInputs)
	if err != nil {
		return nil, fmt.Errorf("error sending message batch: %w", err)
	}

	return resp, nil
}
