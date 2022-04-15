package ctfd

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"path"

	"github.com/dghubble/sling"
	"golang.org/x/xerrors"
)

type Credential struct {
	Username string
	Password string
}

type Client struct {
	sling *sling.Sling
}

func NewClient(url string, cred *Credential) (*Client, error) {
	c, err := httpClientWithCredentials(url, cred)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if strings.HasSuffix(url, "/") {
		url += "/"
	}
	s := sling.New().Base(url).Path("api/v1/").Client(c)
	return &Client{
		sling: s,
	}, nil
}

func httpClientWithCredentials(base string, cred *Credential) (*http.Client, error) {
	u, err := joinPath(base, "login")
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	// get CSRF token for login
	html, cookies, err := getLoginHTML(u.String())
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	token := extractCSRFToken(html)
	if token == "" {
		return nil, xerrors.Errorf("Unable to extract CSRF token")
	}

	// login
	loginClient := httpClientWithCookies(u, cookies)
	resp, err := loginClient.PostForm(u.String(), url.Values{
		"name":     {cred.Username},
		"password": {cred.Password},
		"_submit":  {"Submit"},
		"nonce":    {token},
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, xerrors.Errorf("Bad response %d", resp.StatusCode)
	}

	// return http.Client with logged in cookies
	return httpClientWithCookies(u, resp.Cookies()), nil
}

func getLoginHTML(url string) (string, []*http.Cookie, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", nil, xerrors.Errorf("Bad response %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, xerrors.Errorf(": %w", err)
	}
	return string(html), resp.Cookies(), nil
}

func extractCSRFToken(html string) string {
	r := regexp.MustCompile(`[0-9a-fA-F]{64}`)
	return r.FindString(html)
}

func httpClientWithCookies(u *url.URL, cookies []*http.Cookie) *http.Client {
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)

	c := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return c
}

func joinPath(base string, elem ...string) (*url.URL, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if len(elem) > 0 {
		elem = append([]string{u.Path}, elem...)
		u.Path = path.Join(elem...)
	}
	return u, nil
}
