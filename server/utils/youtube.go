package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/GauravJain98/influencer-coupon/server/models"
)

const youtubeChannelsURL = "https://www.googleapis.com/youtube/v3/channels"

type youtubeChannelsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

func GetChannelIDAndNameFromHandle(ctx context.Context, apiKey string, channel *models.Channel) error {
	if channel.Handle == nil {
		return fmt.Errorf("youtube handle is required")
	}

	handle := strings.TrimSpace(*channel.Handle)
	if handle == "" {
		return fmt.Errorf("youtube handle is required")
	}

	params := url.Values{}
	params.Set("part", "id,snippet")
	params.Set("forHandle", handle)
	params.Set("key", apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeChannelsURL+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("youtube channels request failed with status %s", resp.Status)
	}

	var data youtubeChannelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if len(data.Items) == 0 || data.Items[0].ID == "" {
		return fmt.Errorf("youtube channel not found for handle %q", handle)
	}

	channel.ChannelID = data.Items[0].ID
	channel.Name = &data.Items[0].Snippet.Title
	return nil
}
