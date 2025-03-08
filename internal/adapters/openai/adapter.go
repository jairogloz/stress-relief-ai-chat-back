package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
	"time"
)

type handler struct {
	assistantID string
	client      *openai.Client
	logger      ports.Logger
}

func NewOpenAIAdapter(apiKey, assistantID string, l ports.Logger) ports.ChatHandler {
	if apiKey == "" {
		panic("Cannot create OpenAI handler without an API key")
	}
	h := &handler{
		assistantID: assistantID,
		client:      openai.NewClient(apiKey),
		logger:      l,
	}
	if h.assistantID == "" {
		panic("Cannot create OpenAI handler without an Assistant ID")
	}
	if h.logger == nil {
		panic("Cannot create OpenAI handler without a Logger")
	}
	return h
}

// ProcessMessage
//
// Params:
//   - threadID is an optional parameter that can be used to continue a conversation thread.
func (h *handler) ProcessMessage(ctx context.Context, message *domain.ChatMessage, threadID *string) (response *domain.ChatResponse, err error) {
	if err := message.Validate(); err != nil {
		return nil, fmt.Errorf("invalid message: %s", err.Error())
	}

	var run openai.Run
	if threadID == nil {
		startCreateThread := time.Now().UTC()
		// there's no open thread for the user, so create a new thread and run with the prompt
		run, err = h.client.CreateThreadAndRun(
			ctx,
			openai.CreateThreadAndRunRequest{
				RunRequest: openai.RunRequest{
					AssistantID: h.assistantID,
				},
				Thread: openai.ThreadRequest{
					Messages: []openai.ThreadMessage{
						{
							Role:    openai.ThreadMessageRoleUser,
							Content: message.Content,
						},
					},
				},
			},
		)
		if err != nil {
			h.logger.Error(ctx, "Error creating thread and run", "error", err)
			return nil, fmt.Errorf("could not create thread and run: %w", err)
		}
		h.logger.Debug(ctx, "Thread and run created", "time", time.Since(startCreateThread).String())
		h.logger.Debug(ctx, "Thread and run created", "threadID", run.ThreadID, "runID", run.ID)

		// As a new thread was created, we append it to the user metadata
		//newUserMetadata[MetadataFieldThreadID] = run.ThreadID
	} else {
		h.logger.Debug(ctx, "Thread found for user", "thread_id", threadID)
	}

	startWaitForRunCompletion := time.Now().UTC()
	err = waitForRunCompletion(ctx, h.client, run.ThreadID, run.ID, 500*time.Millisecond)
	if err != nil {
		h.logger.Error(ctx, "Error waiting for run completion", "error", err)
		return nil, fmt.Errorf("could not wait for run completion: %w", err)
	}
	h.logger.Debug(ctx, "Run completed", "time", time.Since(startWaitForRunCompletion).String())

	messageList, err := h.client.ListMessage(ctx, run.ThreadID, domain.IntPtr(1),
		domain.StrPtr("desc"), nil, nil, domain.StrPtr(run.ID))
	if err != nil {
		h.logger.Error(ctx, "Error listing messages", "error", err)
		return nil, fmt.Errorf("could not list messages: %w", err)
	}

	if len(messageList.Messages) == 0 {
		h.logger.Error(ctx, "No messages found")
		return nil, errors.New("no messages found")
	}

	messageResponse := messageList.Messages[0]
	h.logger.Debug(ctx, "Message retrieved", "message_id", messageResponse.ID)
	if msgContent := messageResponse.Content; len(msgContent) > 0 {
		msg := msgContent[0]
		if msgTxt := msg.Text; msgTxt != nil {
			return &domain.ChatResponse{
				Content: (*msgTxt).Value,
			}, nil
		} else {
			return nil, errors.New("no text in message")
		}
	} else {
		return nil, errors.New("no content in message")
	}

}

// waitForRunCompletion waits for the completion of a run in a given thread.
// It periodically checks the status of the run at the specified check interval.
//
// Parameters:
//   - ctx: The context to control cancellation and timeout.
//   - client: The OpenAI client used to retrieve the run status.
//   - threadID: The ID of the thread containing the run.
//   - runID: The ID of the run to wait for completion.
//   - checkInterval: The interval at which to check the run status.
//
// Returns:
//   - error: An error if the context is done or if there is an issue retrieving the run status.
func waitForRunCompletion(ctx context.Context, client *openai.Client, threadID, runID string, checkInterval time.Duration) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			run, err := client.RetrieveRun(ctx, threadID, runID)
			if err != nil {
				return fmt.Errorf("could not retrieve run: %w", err)
			}
			if run.Status == openai.RunStatusCompleted {
				return nil
			}
		}
	}
}
