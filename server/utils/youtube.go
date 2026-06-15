package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/GauravJain98/influencer-coupon/server/models"
)

const youtubeChannelsURL = "https://www.googleapis.com/youtube/v3/channels"
const youtubeSearchURL = "https://www.googleapis.com/youtube/v3/search"

type youtubeChannelsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

type youtubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			PublishedAt string `json:"publishedAt"`
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

func GetLatestChannelVideos(ctx context.Context, apiKey string, channelID string) ([]models.Video, error) {
	return GetChannelVideos(ctx, apiKey, channelID, nil)
}

func GetChannelVideos(ctx context.Context, apiKey string, channelID string, publishedBefore *time.Time) ([]models.Video, error) {
	MAX_RESULTS := "20"
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return nil, fmt.Errorf("youtube channel ID is required")
	}

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("channelId", channelID)
	params.Set("order", "date")
	params.Set("type", "video")
	params.Set("maxResults", MAX_RESULTS)
	params.Set("key", apiKey)
	if publishedBefore != nil {
		params.Set("publishedBefore", publishedBefore.UTC().Format(time.RFC3339))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeSearchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("youtube search request failed with status %s", resp.Status)
	}

	var data youtubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	videos := make([]models.Video, 0, len(data.Items))
	for _, item := range data.Items {
		if item.ID.VideoID == "" {
			continue
		}

		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("parse youtube publishedAt for video %s: %w", item.ID.VideoID, err)
		}

		title := item.Snippet.Title
		description := item.Snippet.Description
		videos = append(videos, models.Video{
			Link:        "https://www.youtube.com/watch?v=" + item.ID.VideoID,
			ChannelID:   channelID,
			Title:       &title,
			Description: &description,
			PublishedAt: &publishedAt,
		})
	}

	return videos, nil
}
