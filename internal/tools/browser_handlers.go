package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/datumbridge/bright-data-mcp/internal/brightdata"
	"github.com/datumbridge/bright-data-mcp/internal/mcp"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var (
	browserMu     sync.Mutex
	browserPage    *rod.Page
	browserRefs    = map[string]string{}
	browserRefSeq  int
)

func registerBrowserTools(cfg ServerConfig, add toolAdder) {
	tools := []struct {
		name, desc string
		props      map[string]interface{}
		required   []string
		h          mcp.ToolHandler
	}{
		{"scraping_browser_navigate", "Navigate scraping browser to URL", baseProps(map[string]interface{}{
			"url": map[string]interface{}{"type": "string"}, "country": map[string]interface{}{"type": "string"},
		}), []string{"url"}, handleBrowserNavigate},
		{"scraping_browser_go_back", "Go back in browser session", baseProps(nil), nil, handleBrowserGoBack},
		{"scraping_browser_go_forward", "Go forward in browser session", baseProps(nil), nil, handleBrowserGoForward},
		{"scraping_browser_snapshot", "ARIA-style snapshot of interactive elements with refs", baseProps(nil), nil, handleBrowserSnapshot},
		{"scraping_browser_click_ref", "Click element by ref from snapshot", baseProps(map[string]interface{}{
			"ref": map[string]interface{}{"type": "string"}, "element": map[string]interface{}{"type": "string"},
		}), []string{"ref", "element"}, handleBrowserClickRef},
		{"scraping_browser_type_ref", "Type into element by ref", baseProps(map[string]interface{}{
			"ref": map[string]interface{}{"type": "string"}, "text": map[string]interface{}{"type": "string"},
			"submit": map[string]interface{}{"type": "boolean"},
		}), []string{"ref", "text"}, handleBrowserTypeRef},
		{"scraping_browser_screenshot", "Screenshot current page (base64 PNG)", baseProps(map[string]interface{}{
			"full_page": map[string]interface{}{"type": "boolean"},
		}), nil, handleBrowserScreenshot},
		{"scraping_browser_network_requests", "List network requests (stub: enable before navigate in future)", baseProps(nil), nil, handleBrowserNetworkRequests},
		{"scraping_browser_wait_for_ref", "Wait for ref element visible", baseProps(map[string]interface{}{
			"ref": map[string]interface{}{"type": "string"}, "timeout_ms": map[string]interface{}{"type": "integer"},
		}), []string{"ref"}, handleBrowserWaitForRef},
		{"scraping_browser_get_text", "Get body text", baseProps(nil), nil, handleBrowserGetText},
		{"scraping_browser_get_html", "Get page HTML", baseProps(map[string]interface{}{
			"full_page": map[string]interface{}{"type": "boolean"},
		}), nil, handleBrowserGetHTML},
		{"scraping_browser_scroll", "Scroll to bottom", baseProps(nil), nil, handleBrowserScroll},
		{"scraping_browser_scroll_to_ref", "Scroll element into view", baseProps(map[string]interface{}{
			"ref": map[string]interface{}{"type": "string"},
		}), []string{"ref"}, handleBrowserScrollToRef},
	}
	for _, t := range tools {
		if !cfg.toolEnabled(t.name) {
			continue
		}
		name := t.name
		h := t.h
		add(name, t.desc, t.props, t.required, func(raw json.RawMessage) map[string]interface{} {
			recordToolCall(name)
			return h(raw)
		})
	}
}

func browserPageOrError(client *brightdata.Client, m map[string]interface{}) (*rod.Page, error) {
	zone := client.BrowserZone()
	if zone == "" {
		return nil, fmt.Errorf("browser_zone required (credentials_json.browser_zone or BRIGHTDATA_BROWSER_ZONE)")
	}
	country := strings.ToLower(strArg(m, "country"))
	browserMu.Lock()
	defer browserMu.Unlock()
	if browserPage != nil {
		return browserPage, nil
	}
	ws, err := brightdata.BrowserCDPEndpoint(client, zone, country)
	if err != nil {
		return nil, err
	}
	u := launcher.MustResolveURL(ws)
	b := rod.New().ControlURL(u).MustConnect()
	browserPage = b.MustPage()
	return browserPage, nil
}

func handleBrowserNavigate(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	target := strArg(m, "url")
	if err := page.Navigate(target); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := page.WaitLoad(); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	browserRefs = map[string]string{}
	browserRefSeq = 0
	title, _ := page.Info()
	return mcp.ToolResultText(fmt.Sprintf("Successfully navigated to %s\nTitle: %s\nURL: %s", target, title.Title, page.MustInfo().URL))
}

func handleBrowserGoBack(raw json.RawMessage) map[string]interface{} {
	return browserNavAction(raw, func(p *rod.Page) error { p.MustNavigateBack(); return nil })
}

func handleBrowserGoForward(raw json.RawMessage) map[string]interface{} {
	return browserNavAction(raw, func(p *rod.Page) error { p.MustNavigateForward(); return nil })
}

func browserNavAction(raw json.RawMessage, fn func(*rod.Page) error) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if err := fn(page); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	info, _ := page.Info()
	return mcp.ToolResultText(fmt.Sprintf("OK\nTitle: %s\nURL: %s", info.Title, info.URL))
}

func handleBrowserSnapshot(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	els, err := page.Elements("a, button, input, textarea, select, [role=button]")
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	browserRefs = map[string]string{}
	browserRefSeq = 0
	var lines []string
	for _, el := range els {
		if browserRefSeq >= 200 {
			break
		}
		browserRefSeq++
		ref := fmt.Sprintf("ref-%d", browserRefSeq)
		sel, _ := el.Attribute("id")
		tag := el.MustEval(`() => this.tagName.toLowerCase()`).String()
		text, _ := el.Text()
		text = strings.TrimSpace(text)
		if len(text) > 80 {
			text = text[:80] + "..."
		}
		if sel != nil && *sel != "" {
			browserRefs[ref] = "#" + *sel
		} else {
			browserRefs[ref] = tag
		}
		lines = append(lines, fmt.Sprintf("[%s] <%s> %s", ref, tag, text))
	}
	return mcp.ToolResultText(strings.Join(lines, "\n"))
}

func handleBrowserClickRef(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ref := strArg(m, "ref")
	sel, ok := browserRefs[ref]
	if !ok {
		return mcp.ToolResultError("unknown ref: " + ref)
	}
	page.MustElement(sel).MustClick()
	return mcp.ToolResultText("clicked " + ref)
}

func handleBrowserTypeRef(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ref := strArg(m, "ref")
	sel, ok := browserRefs[ref]
	if !ok {
		return mcp.ToolResultError("unknown ref: " + ref)
	}
	el := page.MustElement(sel)
	if err := el.Input(strArg(m, "text")); err != nil {
		return mcp.ToolResultError(err.Error())
	}
	if boolArg(m, "submit", false) {
		page.Keyboard.MustType('\n')
	}
	return mcp.ToolResultText("typed into " + ref)
}

func handleBrowserScreenshot(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	full := boolArg(m, "full_page", false)
	b, err := page.Screenshot(full, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return mcp.ToolResultText("data:image/png;base64," + base64.StdEncoding.EncodeToString(b))
}

func handleBrowserNetworkRequests(json.RawMessage) map[string]interface{} {
	return mcp.ToolResultText("[]")
}

func handleBrowserWaitForRef(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ref := strArg(m, "ref")
	sel, ok := browserRefs[ref]
	if !ok {
		return mcp.ToolResultError("unknown ref: " + ref)
	}
	timeout := intArg(m, "timeout_ms", 30000)
	_ = page.Timeout(timeDuration(timeout)).MustElement(sel).MustWaitVisible()
	return mcp.ToolResultText("visible: " + ref)
}

func timeDuration(ms int) time.Duration {
	if ms <= 0 {
		ms = 30000
	}
	return time.Duration(ms) * time.Millisecond
}

func handleBrowserGetText(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	text, err := page.MustElement("body").Text()
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(text)
}

func handleBrowserGetHTML(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	html, err := page.HTML()
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return untrustedTextResult(html)
}

func handleBrowserScroll(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page.MustEval(`() => window.scrollTo(0, document.body.scrollHeight)`)
	return mcp.ToolResultText("scrolled to bottom")
}

func handleBrowserScrollToRef(raw json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(raw)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	page, err := browserPageOrError(client, m)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	ref := strArg(m, "ref")
	sel, ok := browserRefs[ref]
	if !ok {
		return mcp.ToolResultError("unknown ref: " + ref)
	}
	page.MustElement(sel).MustScrollIntoView()
	return mcp.ToolResultText("scrolled to " + ref)
}
