# YouTube Channel Backfill Plan

## Goal

Create a recoverable backfill worker that stores channel videos in batches of 20. For a new channel, the worker first saves the latest 20 videos. After that, each run continues from the oldest saved video and fetches the next 20 older videos.

## YouTube API Constraint

YouTube does not provide a direct channel listing sorted oldest-first. The Search API can return newest-first videos and supports `publishedBefore`, so backfill should use the oldest saved video's publish time as the cursor for the next older batch.

## Data Model Changes

Add `PublishedAt` to `models.Video`:

```go
PublishedAt *time.Time `gorm:"column:published_at;index"`
```

Create a tracking model, for example `models.ChannelVideoBackfill`:

```go
type ChannelVideoBackfill struct {
	gorm.Model
	ChannelID   string     `gorm:"type:text;not null;uniqueIndex"`
	LastRunAt   *time.Time
	CompletedAt *time.Time
	LastError   *string    `gorm:"type:text"`

	Channel Channel `gorm:"foreignKey:ChannelID;references:ChannelID"`
}
```

Use the saved videos themselves as the main progress cursor. The tracking row records completion and last errors so completed channels are not repeatedly queried.

Include `ChannelVideoBackfill` in `app.Migrate`.

## YouTube Utilities

Update `utils/youtube.go` to support video search.

Use:

```text
GET https://www.googleapis.com/youtube/v3/search
```

Common params:

```text
part=snippet
channelId=<channelID>
order=date
type=video
maxResults=20
key=<apiKey>
```

For older batches, add:

```text
publishedBefore=<oldestPublishedAt RFC3339>
```

Implement or update:

```go
func GetLatestChannelVideos(ctx context.Context, apiKey string, channelID string) ([]models.Video, error)
```

Behavior:

- Fetch the latest 20 videos for a channel.
- Parse title, description, video ID, and `snippet.publishedAt`.
- Store links as `https://www.youtube.com/watch?v=<videoID>`.

Implement or update:

```go
func GetChannelVideos(ctx context.Context, apiKey string, channelID string, publishedBefore *time.Time) ([]models.Video, error)
```

Behavior:

- Fetch one batch of up to 20 videos.
- If `publishedBefore` is set, fetch videos older than that timestamp.
- Return only one page so the worker is recoverable.
- Do not use a `fetchAll` loop for backfill.

## New Worker

Create:

```go
func NewChannelBackfillWorker(config config.Config, db *gorm.DB)
```

Flow per channel:

1. Load channels from the database.
2. Skip channels with an empty `ChannelID`.
3. Find or create a `ChannelVideoBackfill` row.
4. Skip the channel if `CompletedAt` is already set.
5. Query the oldest saved video for that channel where `published_at IS NOT NULL`:

```go
db.Where("channel_id = ? AND published_at IS NOT NULL", channel.ChannelID).
	Order("published_at ASC").
	First(&oldestVideo)
```

6. If no saved video exists, call `GetLatestChannelVideos` and save the latest 20 videos.
7. If an oldest saved video exists, call `GetChannelVideos` with `oldestVideo.PublishedAt` to fetch the next 20 older videos.
8. Upsert videos by `link`.
9. If no videos are returned, set `CompletedAt`.
10. Set `LastRunAt` on each attempted channel.
11. On error, set `LastError` and continue to the next channel.

## Recovery Behavior

The worker should not depend on persisted page tokens.

If it crashes before saving a batch, the next run retries the same batch.

If it crashes after saving a batch, the next run finds the new oldest saved video and continues from there.

If a channel has no saved videos, it starts with the latest 20 first. Future runs move older from the oldest saved video.

## Existing Worker Compatibility

Keep `VideoListFetcherWorker` for regular latest-video refreshes.

Update it to call the revised utility signatures. If it only needs latest videos, it should call `GetLatestChannelVideos`.

## Verification

After implementation:

```sh
gofmt -w .
go test ./...
```
