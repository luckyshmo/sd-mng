package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/ratelimit"
)

var (
	limitGlobal = ratelimit.New(5, ratelimit.Per(time.Second))
	limitAtHome = ratelimit.New(40, ratelimit.Per(time.Minute))
)

var APIBaseURL, _ = url.Parse(`https://api.mangadex.org/`)

type Client struct {
	http    *http.Client
	baseURL url.URL
}

func NewClient() *Client {
	return &Client{
		http:    http.DefaultClient,
		baseURL: *APIBaseURL,
	}
}

func (c *Client) WithBaseURL(url url.URL) *Client {
	c.baseURL = url
	return c
}

func (c *Client) WithHTTPClient(http *http.Client) *Client {
	c.http = http
	return c
}

func (c *Client) GetManga(ctx context.Context, mangaID string) (*Manga, error) {
	v := new(Manga)
	err := c.doJSON(ctx, "GET", "/manga/"+mangaID, v, nil)
	return v, err
}

func (c *Client) GetFeed(ctx context.Context, mangaID string, args QueryArgs) (*ChapterList, error) {
	v := new(ChapterList)
	url := fmt.Sprintf("/manga/%v/feed?%v", mangaID, args.Values().Encode())
	err := c.doJSON(ctx, "GET", url, v, nil)
	return v, err
}

func (c *Client) GetCovers(ctx context.Context, args QueryArgs) (*CoverList, error) {
	v := new(CoverList)
	err := c.doJSON(ctx, "GET", "/cover?"+args.Values().Encode(), v, nil)
	return v, err
}

func (c *Client) GetAuthors(ctx context.Context, args QueryArgs) (*AuthorList, error) {
	v := new(AuthorList)
	err := c.doJSON(ctx, "GET", "/author?"+args.Values().Encode(), v, nil)
	return v, err
}

func (c *Client) GetGroups(ctx context.Context, args QueryArgs) (*GroupList, error) {
	v := new(GroupList)
	err := c.doJSON(ctx, "GET", "/group?"+args.Values().Encode(), v, nil)
	return v, err
}

func (c *Client) GetAtHome(ctx context.Context, chapterID string) (*AtHome, error) {
	v := new(AtHome)
	limitAtHome.Take()
	err := c.doJSON(ctx, "GET", "/at-home/server/"+chapterID, v, nil)
	return v, err
}

func (c *Client) PostIDMapping(ctx context.Context, tp string, legacyIDs ...int) (*IDMappingList, error) {
	v := new(IDMappingList)
	err := c.doJSON(ctx, "POST", "/legacy/mapping", &v, map[string]interface{}{
		"ids":  legacyIDs,
		"type": tp,
	})

	return v, err
}

func (c *Client) doJSON(ctx context.Context, method, ref string, result, body interface{}) error {
	url, err := c.baseURL.Parse(ref)
	if err != nil {
		return fmt.Errorf("url: %w", err)
	}

	rw := io.ReadWriter(nil)
	if body != nil {
		rw = bytes.NewBuffer(nil)
		if err := json.NewEncoder(rw).Encode(body); err != nil {
			return fmt.Errorf("encode: %w", err)
		}
	}

	fmt.Println(url.String()) //! del
	req, err := http.NewRequestWithContext(ctx, method, url.String(), rw)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	limitGlobal.Take()
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errs := new(Errors)
		if err := dec.Decode(errs); err != nil {
			return fmt.Errorf("error decode: %w", err)
		} else if len(errs.Errors) != 0 {
			return fmt.Errorf("detail: %v", errs.Errors[0].Detail)
		} else {
			return fmt.Errorf("status: %v", resp.Status)
		}
	} else if err := dec.Decode(result); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
