package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/viper"

	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/browser"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/server"
	"github.com/Phillezi/kthcloud-cli/pkg/v1/auth/session"
	storageclient "github.com/Phillezi/kthcloud-cli/pkg/v1/auth/storage-client"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
)

type Client struct {
	kcBaseURL     string
	baseURL       string
	clientID      string
	clientSecret  string
	realm         string
	client        *resty.Client
	jar           http.CookieJar
	Session       *session.Session
	StorageClient *storageclient.Client
}

var (
	instance *Client
	once     sync.Once
)

func Get() *Client {
	return GetInstance(
		viper.GetString("api-url"),
		viper.GetString(
			"keycloak-host"),
		viper.GetString("client-id"),
		"",
		viper.GetString("keycloak-realm"),
	)
}

func (c *Client) Client() *resty.Client {
	return c.client
}

func GetInstance(baseURL, kcBaseURL, clientID, clientSecret, realm string) *Client {
	once.Do(func() {
		client := resty.New()
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Error creating cookie jar: %v", err)
		}
		client.SetCookieJar(jar)
		sess, err := session.Load(viper.GetString("session-path"))
		if err != nil || sess.IsExpired() {
			// TODO: try to refresh token here later
			sess = nil
		}
		instance = &Client{
			baseURL:       baseURL,
			kcBaseURL:     kcBaseURL,
			clientID:      clientID,
			clientSecret:  clientSecret,
			realm:         realm,
			client:        client,
			jar:           jar,
			Session:       sess,
			StorageClient: nil,
		}
		if sess != nil {
			instance.client.SetAuthToken(instance.Session.Token.AccessToken)
		}
		instance.client.SetBaseURL(baseURL)
	})
	return instance
}

func (c *Client) HasValidSession() bool {
	return c.Session != nil && !c.Session.IsExpired()
}

func (c *Client) Login() (*session.Session, error) {
	kcURL := c.generateKCUrl()

	sessionChannel := make(chan *session.Session)
	server := server.New(":3000", kcURL, sessionChannel, c.fetchOAuthToken, context.Background())

	server.Start()

	err := browser.Open("http://localhost:3000")
	if err != nil {
		return nil, err
	}

	fmt.Println("Waiting for login...")

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue")
	s.Start()

	session := <-sessionChannel

	s.Stop()

	if session != nil {
		c.Session = session
		c.client.SetAuthToken(c.Session.Token.AccessToken)
	}

	return session, nil
}

func (c *Client) Logout() error {
	c.Session = nil
	return c.Session.Save(viper.GetString("session-path"))
}

func (c *Client) fetchOAuthToken(redirectURI, code string) (*http.Response, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.kcBaseURL, c.realm)

	// Create the URL-encoded form data
	form := url.Values{}
	form.Add("client_id", c.clientID)
	form.Add("redirect_uri", redirectURI)
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)

	// Create the POST request
	req, err := http.NewRequestWithContext(context.Background(), "POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}

func (c *Client) generateKCUrl() string {
	redirectURI := fmt.Sprintf("%s/auth/callback", "http://localhost:3000")
	state := generateRandomState()
	nonce := generateRandomState()

	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&response_mode=query&scope=openid&nonce=%s&state=%s",
		c.kcBaseURL, c.realm, c.clientID, url.QueryEscape(redirectURI), nonce, state)
}

func (c *Client) RedirectToKeycloak(w http.ResponseWriter) {
	http.Redirect(w, nil, c.generateKCUrl(), http.StatusFound)
}

// HandleCallback processes the authorization code returned by Keycloak and exchanges it for tokens.
func (c *Client) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.kcBaseURL, c.realm)

	resp, err := c.client.R().
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  "http://localhost:3000/auth/callback",
		}).
		Post(tokenURL)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error requesting token: %v", err), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to exchange code for token: %s", resp.String()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Authentication successful! Token response: %s", resp.String())
}

func generateRandomState() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	const charset = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	randomize := func(c rune) rune {
		r := rand.Intn(16)
		if c == 'x' {
			return rune(fmt.Sprintf("%x", r)[0])
		}
		return rune(fmt.Sprintf("%x", (r&0x3)|0x8)[0])
	}
	return strings.Map(randomize, charset)
}
