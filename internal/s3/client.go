package s3

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/shirwahersi/go-s3-browser/internal/config"
)

type Folder struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Type string `json:"type"`
}

type File struct {
	Name         string    `json:"name"`
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Type         string    `json:"type"`
}

type ListResult struct {
	Folders []Folder `json:"folders"`
	Files   []File   `json:"files"`
}

type Client struct {
	s3Client       *s3.Client
	presignClient  *s3.PresignClient
	bucket         string
}

// NewClient creates a new S3 client configured for the given settings.
func NewClient(cfg *config.Config) (*Client, error) {
	s3Cfg := aws.Config{
		Region: cfg.S3.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
			"",
		),
	}

	s3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &Client{
		s3Client:      s3Client,
		presignClient: presignClient,
		bucket:        cfg.S3.Bucket,
	}, nil
}

// ListObjects lists objects at the given prefix, returning folders and files separately.
func (c *Client) ListObjects(ctx context.Context, prefix string) (*ListResult, error) {
	delimiter := "/"
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(c.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	}

	output, err := c.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	result := &ListResult{
		Folders: make([]Folder, 0),
		Files:   make([]File, 0),
	}

	// Extract folders from CommonPrefixes
	for _, cp := range output.CommonPrefixes {
		if cp.Prefix != nil {
			folderName := getFolderName(*cp.Prefix, prefix)
			result.Folders = append(result.Folders, Folder{
				Name: folderName,
				Key:  *cp.Prefix,
				Type: "folder",
			})
		}
	}

	// Extract files from Contents (excluding the prefix itself)
	for _, obj := range output.Contents {
		if obj.Key != nil && *obj.Key != prefix {
			fileName := getFileName(*obj.Key, prefix)
			file := File{
				Name: fileName,
				Key:  *obj.Key,
				Type: "file",
			}
			if obj.Size != nil {
				file.Size = *obj.Size
			}
			if obj.LastModified != nil {
				file.LastModified = *obj.LastModified
			}
			result.Files = append(result.Files, file)
		}
	}

	return result, nil
}

// GetPresignedURL generates a pre-signed URL for downloading the given key.
func (c *Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	presignedReq, err := c.presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

// getFolderName extracts the folder name from a full prefix.
func getFolderName(fullPrefix, currentPrefix string) string {
	relative := strings.TrimPrefix(fullPrefix, currentPrefix)
	return strings.TrimSuffix(relative, "/")
}

// getFileName extracts the file name from a full key.
func getFileName(fullKey, currentPrefix string) string {
	return strings.TrimPrefix(fullKey, currentPrefix)
}
