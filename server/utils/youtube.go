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
const youtubeVideosURL = "https://www.googleapis.com/youtube/v3/videos"

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

type youtubeVideosResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			ChannelID   string `json:"channelId"`
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

func NOGetVideoDetails(ctx context.Context, apiKey string, video *models.Video) error {
	if video == nil {
		return fmt.Errorf("video is required")
	}

	videoID := strings.TrimSpace(video.Link)
	if videoID == "" {
		return fmt.Errorf("youtube video link is required")
	}

	if parsedURL, err := url.Parse(videoID); err == nil && parsedURL.Host != "" {
		if parsedURL.Host == "youtu.be" {
			videoID = strings.Trim(parsedURL.Path, "/")
		} else {
			videoID = parsedURL.Query().Get("v")
		}
	}

	videoID = strings.TrimSpace(videoID)
	if videoID == "" {
		return fmt.Errorf("youtube video ID is required")
	}

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("id", videoID)
	params.Set("key", apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeVideosURL+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("youtube videos request failed with status %s", resp.Status)
	}

	var data youtubeVideosResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if len(data.Items) == 0 || data.Items[0].ID == "" {
		return fmt.Errorf("youtube video not found for ID %q", videoID)
	}

	publishedAt, err := time.Parse(time.RFC3339, data.Items[0].Snippet.PublishedAt)
	if err != nil {
		return fmt.Errorf("parse youtube publishedAt for video %s: %w", videoID, err)
	}

	title := data.Items[0].Snippet.Title
	description := data.Items[0].Snippet.Description
	video.Link = "https://www.youtube.com/watch?v=" + data.Items[0].ID
	video.ChannelID = data.Items[0].Snippet.ChannelID
	video.Title = &title
	video.Description = &description
	video.PublishedAt = &publishedAt
	return nil
}

func ListVideoDetails(ctx context.Context, apiKey string, videos []models.Video) (time.Time, []models.Video, error) {
	lastPublishedTime := time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC)
	if len(videos) == 0 {
		return lastPublishedTime, videos, nil
	}

	videoIndexes := make(map[string][]int, len(videos))
	videoIDs := make([]string, 0, len(videos))
	for i := range videos {
		videoID := strings.TrimSpace(videos[i].Link)
		if videoID == "" {
			return lastPublishedTime, videos, fmt.Errorf("youtube video link is required")
		}

		if parsedURL, err := url.Parse(videoID); err == nil && parsedURL.Host != "" {
			if parsedURL.Host == "youtu.be" {
				videoID = strings.Trim(parsedURL.Path, "/")
			} else {
				videoID = parsedURL.Query().Get("v")
			}
		}

		videoID = strings.TrimSpace(videoID)
		if videoID == "" {
			return lastPublishedTime, videos, fmt.Errorf("youtube video ID is required")
		}

		if _, ok := videoIndexes[videoID]; !ok {
			videoIDs = append(videoIDs, videoID)
		}
		videoIndexes[videoID] = append(videoIndexes[videoID], i)
	}

	const maxYoutubeVideoIDs = 50

	for start := 0; start < len(videoIDs); start += maxYoutubeVideoIDs {
		end := start + maxYoutubeVideoIDs
		if end > len(videoIDs) {
			end = len(videoIDs)
		}

		params := url.Values{}
		params.Set("part", "snippet")
		params.Set("id", strings.Join(videoIDs[start:end], ","))
		params.Set("key", apiKey)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeVideosURL+"?"+params.Encode(), nil)
		if err != nil {
			return lastPublishedTime, videos, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return lastPublishedTime, videos, err
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			resp.Body.Close()
			return lastPublishedTime, videos, fmt.Errorf("youtube videos request failed with status %s", resp.Status)
		}

		var data youtubeVideosResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			return lastPublishedTime, videos, err
		}
		resp.Body.Close()

		for _, item := range data.Items {
			indexes := videoIndexes[item.ID]
			if len(indexes) == 0 {
				continue
			}

			publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
			if err != nil {
				return lastPublishedTime, videos, fmt.Errorf("parse youtube publishedAt for video %s: %w", item.ID, err)
			}

			for _, index := range indexes {
				title := item.Snippet.Title
				description := item.Snippet.Description
				videos[index].Link = "https://www.youtube.com/watch?v=" + item.ID
				videos[index].ChannelID = item.Snippet.ChannelID
				videos[index].Title = &title
				videos[index].Description = &description
				videos[index].PublishedAt = &publishedAt
				if publishedAt.After(lastPublishedTime) {
					lastPublishedTime = publishedAt
				}
			}
			delete(videoIndexes, item.ID)
		}
	}

	if len(videoIndexes) > 0 {
		return lastPublishedTime, videos, fmt.Errorf("youtube video not found for IDs %s", fmt.Sprint(videoIndexes))
	}

	return lastPublishedTime, videos, nil
}

func GetChannelVideos(ctx context.Context, apiKey string, channelID string, publishedBefore time.Time) ([]models.Video, error) {
	MAX_RESULTS := "50"
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
	params.Set("publishedBefore", publishedBefore.UTC().Format(time.RFC3339))

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
