package web

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/chromedp/chromedp"
)

type Client struct {
	ctx     *context.Context
	actions []chromedp.Action
	cancels []context.CancelFunc
}

func NewClient(ctx *context.Context) *Client {
	return &Client{ctx: ctx}
}

func (c *Client) Context() context.Context {
	return *c.ctx
}

func (c *Client) Add(v chromedp.Action) *Client {
	c.actions = append(c.actions, v)
	return c
}

// Navigation

func (c *Client) Go(url string) *Client {
	c.actions = append(c.actions, chromedp.Navigate(url))
	return c
}

func (c *Client) Back() *Client {
	c.actions = append(c.actions, chromedp.NavigateBack())
	return c
}

func (c *Client) Forward() *Client {
	c.actions = append(c.actions, chromedp.NavigateForward())
	return c
}

func (c *Client) Reload() *Client {
	c.actions = append(c.actions, chromedp.Reload())
	return c
}

// Wait

func (c *Client) WaitVisible(selector string) *Client {
	c.actions = append(c.actions, chromedp.WaitVisible(selector, chromedp.ByQuery))
	return c
}

func (c *Client) WaitNotVisible(selector string) *Client {
	c.actions = append(c.actions, chromedp.WaitNotVisible(selector, chromedp.ByQuery))
	return c
}

func (c *Client) WaitEnabled(selector string) *Client {
	c.actions = append(c.actions, chromedp.WaitEnabled(selector, chromedp.ByQuery))
	return c
}

func (c *Client) WaitReady(selector string) *Client {
	c.actions = append(c.actions, chromedp.WaitReady(selector, chromedp.ByQuery))
	return c
}

func (c *Client) WaitSelected(selector string) *Client {
	c.actions = append(c.actions, chromedp.WaitSelected(selector, chromedp.ByQuery))
	return c
}

func (c *Client) Sleep(d time.Duration) *Client {
	c.actions = append(c.actions, chromedp.Sleep(d))
	return c
}

// User Actions

func (c *Client) Click(selector string) *Client {
	c.actions = append(c.actions, chromedp.Click(selector, chromedp.NodeVisible, chromedp.ByQuery))
	return c
}

func (c *Client) SendKeys(selector string, v string) *Client {
	c.actions = append(c.actions, chromedp.SendKeys(selector, v, chromedp.ByQuery))
	return c
}

// Get Information

func (c *Client) Value(selector string, value *string) *Client {
	c.actions = append(c.actions, chromedp.Value(selector, value, chromedp.ByQuery))
	return c
}

func (c *Client) Text(selector string, value *string) *Client {
	c.actions = append(c.actions, chromedp.Text(selector, value, chromedp.ByQuery))
	return c
}

// Screenshot

type ScreenshotAction struct {
	picPath  string
	selector string
}

func (s *ScreenshotAction) Do(ctx context.Context) error {
	var buf []byte

	if err := chromedp.Run(ctx, chromedp.Screenshot(s.selector, &buf, chromedp.ByQuery)); err != nil {
		return err
	}

	if err := ioutil.WriteFile(s.picPath, buf, 0644); err != nil {
		return err
	}
	return nil
}

func (c *Client) Screenshot(selector string, picPath string) *Client {
	c.actions = append(c.actions, &ScreenshotAction{
		picPath:  picPath,
		selector: selector,
	})
	return c
}

// Browser Config

func (c *Client) Timeout(t time.Duration) *Client {
	ctx, cancel := context.WithTimeout(*c.ctx, t)
	c.ctx = &ctx
	c.cancels = append(c.cancels, cancel)
	return c
}

func (c *Client) Do() {
	defer func() {
		if c.cancels == nil {
			return
		}

		for _, cancel := range c.cancels {
			cancel()
		}
	}()

	err := chromedp.Run(*c.ctx, c.actions...)
	if err != nil {
		panic(err)
	}
}
