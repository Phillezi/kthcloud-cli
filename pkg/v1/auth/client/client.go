package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
	"golang.org/x/term"
)

type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	realm        string
	client       *resty.Client
	jar          http.CookieJar
}

var (
	instance *Client
	once     sync.Once
)

func GetInstance(baseURL, clientID, clientSecret, realm string) *Client {
	once.Do(func() {
		client := resty.New()
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Error creating cookie jar: %v", err)
		}
		client.SetCookieJar(jar)
		instance = &Client{
			baseURL:      baseURL,
			clientID:     clientID,
			clientSecret: clientSecret,
			realm:        realm,
			client:       resty.New(),
			jar:          jar,
		}
	})
	return instance
}

// Authenticate handles the initial authentication and returns an access token.
// It detects if MFA is required and handles it.
func (c *Client) Authenticate(username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, c.realm)
	// Step 1: Request the token with username and password
	resp, err := c.client.R().
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
			"username":      username,
			"password":      password,
			"grant_type":    "password",
		}).
		Post(tokenURL)

	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	if resp.StatusCode() == 401 && resp.Header().Get("WWW-Authenticate") == "mfa_required" {
		fmt.Println("MFA required. Please enter your MFA code:")
		var mfaCode string
		fmt.Scanln(&mfaCode)

		// Step 2: Send MFA code
		resp, err = c.client.R().
			SetFormData(map[string]string{
				"client_id":     c.clientID,
				"client_secret": c.clientSecret,
				"username":      username,
				"password":      password,
				"grant_type":    "password",
				"totp":          mfaCode, // MFA code from user
			}).
			Post(tokenURL)

		if err != nil {
			return "", fmt.Errorf("error sending MFA request: %v", err)
		}
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("authentication failed: %s", resp.String())
	}

	// Extract the token from the response
	tokenResp, err := util.ProcessResponse[map[string]interface{}](resp.String())
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	accessToken := (*tokenResp)["access_token"].(string)
	return accessToken, nil
}

func (c *Client) Authv2() {

	c.DoAuth()
}

func (c *Client) login(url string, cookies []*http.Cookie) ([]*http.Cookie, error) {

	newURL, err := replaceRedirectURI(url, "http://localhost:3000/auth/callback")
	if err != nil {
		return nil, err
	}

	cookieChan := make(chan []*http.Cookie)

	go func() {
		server := &http.Server{Addr: ":3000"}

		// Serve the / endpoint
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			for _, cookie := range cookies {
				http.SetCookie(w, cookie)
			}

			http.Redirect(w, r, newURL, http.StatusFound)
		})

		http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				//http.Redirect(w, r, newURL, http.StatusFound)
				fmt.Fprintln(w, "no code provided")
				http.Redirect(w, r, newURL, http.StatusFound)
				return
			}

			for _, cookie := range r.Cookies() {
				fmt.Println(cookie)
				http.SetCookie(w, cookie)
			}

			resp, err := c.fetchOAuthToken("http://localhost:3000/auth/callback", code)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %v\n", err)
				return
			}
			fmt.Println(string(body))

			//fmt.Fprintln(w, "code: ", code)
			http.Redirect(w, r, "http://localhost:3000/auth/done", http.StatusFound)
		})

		// Handle /auth/callback
		http.HandleFunc("/auth/done", func(w http.ResponseWriter, r *http.Request) {
			collectedCookies := r.Cookies()
			cookieChan <- collectedCookies
			fmt.Fprintln(w, "Callback received. Server will now shut down.")

			// Close the server
			go func() {
				time.Sleep(500 * time.Millisecond)
				if err := server.Shutdown(context.Background()); err != nil {
					log.Fatalf("Server Shutdown Failed:%+v", err)
				}
				fmt.Println("Server stopped after serving the callback request")
			}()
		})

		// Start the server on localhost:3000
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	// Open the browser with the authURL
	err = OpenBrowser("http://localhost:3000")
	if err != nil {
		return nil, err
	}
	return <-cookieChan, nil
}

func (c *Client) fetchOAuthToken(redirectURI, code string) (*http.Response, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, c.realm)

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

func extractParams(rawURL string) (map[string]string, map[string]string, error) {
	// Split the URL into main part and fragment
	urlParts := strings.Split(rawURL, "#")
	if len(urlParts) < 2 {
		return nil, nil, fmt.Errorf("URL does not contain a fragment")
	}

	baseURL := urlParts[0]
	fragment := urlParts[1]

	// Parse the base URL query parameters
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
	}
	baseParams := base.Query()

	// Parse the fragment parameters
	fragmentParams, err := url.ParseQuery(fragment)
	if err != nil {
		return nil, nil, err
	}

	// Convert the query params from url.Values to a map
	baseParamsMap := make(map[string]string)
	for key, values := range baseParams {
		if len(values) > 0 {
			baseParamsMap[key] = values[0]
		}
	}

	fragmentParamsMap := make(map[string]string)
	for key, values := range fragmentParams {
		if len(values) > 0 {
			fragmentParamsMap[key] = values[0]
		}
	}

	return baseParamsMap, fragmentParamsMap, nil
}

func replaceRedirectURI(authURL, newRedirectURI string) (string, error) {
	// Parse the authURL
	parsedURL, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("error parsing authURL: %w", err)
	}

	// Parse the query parameters
	queryParams := parsedURL.Query()

	// Update the redirect_uri parameter
	queryParams.Set("redirect_uri", newRedirectURI)

	// Reconstruct the URL with the updated query parameters
	parsedURL.RawQuery = queryParams.Encode()
	return parsedURL.String(), nil
}

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch {
	case runtime.GOOS == "linux":
		cmd = exec.Command("xdg-open", url)
	case runtime.GOOS == "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case runtime.GOOS == "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	fmt.Printf("Trying to open: %s in web browser\n\n", url)
	return cmd.Start()
}

func (c *Client) DoAuth() {
	/*initialResponse, err := c.client.R().
		Get(c.generateKCUrl())
	if err != nil {
		log.Fatalf("Error initiating request: %v", err)
	}

	redirectToKthURL, err := extractURL(initialResponse.String())
	if err != nil {
		log.Fatal(err)
	}

	c.client.SetRedirectPolicy(resty.NoRedirectPolicy())

	cookies := initialResponse.Cookies()

	kthResp, err := c.client.R().
		SetCookies(cookies). // Set collected cookies for the request
		Get(viper.GetString("keycloak-host") + redirectToKthURL)
	if err != nil && !strings.Contains(err.Error(), "auto redirect is disabled") {
		log.Fatalf("Error initiating request after redirect: %v", err)
	}

	c.client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))

	cookies = append(cookies, kthResp.Cookies()...)

	kthLoginURL := kthResp.Header().Get("Location")

	kthResp, err = c.client.R().
		SetCookies(cookies).
		Get(kthLoginURL)
	if err != nil {
		log.Fatal(err)
	}
	cookies = append(cookies, kthResp.Cookies()...)

	clientReqID, err := extractClientRequestID(kthResp.String())
	if err != nil {
		log.Fatal(err)
	}*/

	cookies, err := c.login(c.generateKCUrl(), []*http.Cookie{})
	if err != nil {
		log.Fatal(err)
	}

	/*username, password, err := getUsernameAndPassword()
	if err != nil {
		log.Fatal(err)
	}
	formData := map[string]string{
		"UserName":   username,
		"Password":   password,
		"AuthMethod": "FormsAuthentication",
	}

	kthResp, err = c.client.R().
		SetCookies(cookies).
		SetFormData(formData).
		Post(kthLoginURL + "&client-request-id=" + clientReqID)
	if err != nil {
		log.Fatal(err)
	}
	cookies = append(cookies, kthResp.Cookies()...)

	mfaNumFound := false
	mfaNumPrinted := false
	printedMfaNUM := ""

	mfaNum, err := extractElementByID(kthResp.String(), "validEntropyNumber")
	if err == nil {
		mfaNumFound = true
	}

	for {
		if mfaNumFound && !mfaNumPrinted || ((printedMfaNUM != mfaNum) && mfaNumPrinted) {
			fmt.Println("\n\n" + "MFA NUMBER: " + mfaNum + "\n")
			mfaNumPrinted = true
		}
		kthResp, err = c.client.R().
			SetCookies(cookies).
			Post(kthLoginURL + "&client-request-id=" + clientReqID)
		if err != nil {
			log.Fatal(err)
		}
		mfaNum, err = extractElementByID(kthResp.String(), "validEntropyNumber")
		if err == nil {
			mfaNumFound = true
		} else {
			fmt.Println(err)
		}
		if kthResp.StatusCode() == 302 {
			cookies = append(cookies, kthResp.Cookies()...)
			break
		}
		time.Sleep(1 * time.Second)
	}

	kthResp, err = c.client.R().
		SetCookies(cookies).
		Get(kthLoginURL + "&client-request-id=" + clientReqID)
	if err != nil {
		log.Fatal(err)
	}
	cookies = append(cookies, kthResp.Cookies()...)*/
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}

	fmt.Println("Authentication flow completed.")
}

func getUsernameAndPassword() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for username
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("error reading username: %v", err)
	}
	username = username[:len(username)-1] // Remove trailing newline

	// Prompt for password
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", fmt.Errorf("error reading password: %v", err)
	}
	password := string(passwordBytes)

	fmt.Println()

	return username, password, nil
}

func extractElementByID(htmlContent, id string) (string, error) {
	// Load the HTML content into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to load HTML: %v", err)
	}

	// Find the element with the specific id
	selection := doc.Find(fmt.Sprintf("#%s", id))
	if selection.Length() == 0 {
		return "", fmt.Errorf("element with id %s not found", id)
	}

	// Extract and return the HTML content of the element
	return selection.Text(), nil
}

func extractClientRequestID(html string) (string, error) {
	// Define the regular expression pattern
	pattern := `client-request-id=([^&"\s]+)`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %v", err)
	}

	// Find the match
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		return "", fmt.Errorf("client-request-id not found in the HTML")
	}

	// Return the captured value
	return match[1], nil
}

func extractURL(html string) (string, error) {
	// Define the regular expression pattern
	// This pattern looks for href attributes starting with the specified prefix and captures the URL
	pattern := `href="/realms/cloud/broker/oidc/login\?client_id=landing[^"]*"`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %v", err)
	}

	// Find the first match
	match := re.FindString(html)
	if match == "" {
		return "", fmt.Errorf("no matching URL found")
	}

	// Extract the URL part from the href attribute
	url := match[len(`href="`) : len(match)-1]
	return url, nil
}

func (c *Client) generateKCUrl() string {
	redirectURI := fmt.Sprintf("%s/auth/callback", "http://localhost:3000")
	state := generateRandomState()
	nonce := generateRandomState()

	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&response_mode=query&scope=openid&nonce=%s&state=%s",
		c.baseURL, c.realm, c.clientID, url.QueryEscape(redirectURI), nonce, state)
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

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseURL, c.realm)

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
