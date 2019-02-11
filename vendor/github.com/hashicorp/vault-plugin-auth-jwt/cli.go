package jwtauth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/api"
)

const defaultMount = "oidc"
const defaultPort = "8300"

type CLIHandler struct{}

type loginResp struct {
	secret *api.Secret
	err    error
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	// handle ctrl-c while waiting for the callback
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)

	doneCh := make(chan loginResp)

	mount, ok := m["mount"]
	if !ok {
		mount = defaultMount
	}

	port, ok := m["port"]
	if !ok {
		port = defaultPort
	}

	role := m["role"]
	if role == "" {
		return nil, errors.New("a 'role' must be specified")
	}

	secret, err := fetchAuthURL(c, role, mount, port)
	if err != nil {
		return nil, err
	}

	authURL := secret.Data["auth_url"].(string)
	if authURL == "" {
		return nil, errors.New(fmt.Sprintf("Unable to authorize role %q. Check Vault logs for more information.", role))
	}

	fmt.Fprintf(os.Stderr, "Complete the login via your OIDC provider. Launching browser to:\n\n    %s\n\n\n", authURL)
	if err := openURL(authURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error attempting to automatically open browser: '%s'.\nPlease visit the authorization URL manually.", err)
	}

	// Set up callback handler
	http.HandleFunc(fmt.Sprintf("/v1/auth/%s/oidc/callback", mount), func(w http.ResponseWriter, req *http.Request) {
		var response string

		query := req.URL.Query()
		code := query.Get("code")
		state := query.Get("state")
		data := map[string][]string{
			"code":  {code},
			"state": {state},
		}

		secret, err := c.Logical().ReadWithData(fmt.Sprintf("auth/%s/oidc/callback", mount), data)
		if err != nil {
			summary, detail := parseError(err)
			response = errorHTML(summary, detail)
		} else {
			response = successHTML
		}

		w.Write([]byte(response))
		doneCh <- loginResp{secret, err}
	})

	// Start local server
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error listening for callback: %v\n\n", err.Error())
		}
	}()

	// Wait for either the callback to finish or SIGINT to be received
	select {
	case s := <-doneCh:
		return s.secret, s.err
	case <-ch:
		return nil, errors.New("interrupted")
	}
}

func fetchAuthURL(c *api.Client, role, mount, port string) (*api.Secret, error) {
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": fmt.Sprintf("http://localhost:%s/v1/auth/%s/oidc/callback", port, mount),
	}

	return c.Logical().Write(fmt.Sprintf("auth/%s/oidc/auth_url", mount), data)
}

// openURL opens the specified URL in the default browser of the user.
// Source: https://stackoverflow.com/a/39324149/453290
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// parseError converts error from the API into summary and detailed portions.
func parseError(err error) (string, string) {
	headers := []string{errNoResponse, errLoginFailed, errTokenVerification}
	summary := "Login error"
	detail := ""

	re := regexp.MustCompile(`(?s)Errors:.*\* *(.*)`)

	errorParts := re.FindStringSubmatch(err.Error())
	switch len(errorParts) {
	case 0:
		summary = ""
	case 1:
		detail = errorParts[0]
	case 2:
		for _, h := range headers {
			if strings.HasPrefix(errorParts[1], h) {
				summary = h
				detail = strings.TrimSpace(errorParts[1][len(h):])
				break
			}
		}
		if detail == "" {
			detail = errorParts[1]
		}
	}

	return summary, detail

}

func errorHTML(summary, detail string) string {
	const html = `
<!DOCTYPE html>
<html lang="en" >

<head>

  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
<title>HashiCorp Vault</title>
      <style>
      body {
  font-size: 14px;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
    "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans",
    "Helvetica Neue", sans-serif;
}
hr {
  border-color: #fdfdfe;
  margin: 24px 0;
}
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 70vh;
}
#logo {
  display: block;
  fill: #6f7682;
  margin-bottom: 16px;
}
.message {
  display: flex;
  min-width: 40vw;
  background: #fafdfa;
  border: 1px solid #c6e9c9;
  margin-bottom: 12px;
  padding: 12px 16px 16px 12px;
  position: relative;
  border-radius: 2px;
  font-size: 14px;
}
.message.is-danger {
  background: #fdfafb;
  border-color: #f9ecee;
}
.message-content {
  margin-left: 4px;
}
.message svg {
  fill: #2eb039;
}

.message.is-danger svg {
  fill: #c73445;
}
.message .message-title {
  color: #1e7125;
  font-size: 16px;
  font-weight: 700;
  line-height: 1.25;
}
.message.is-danger .message-title {
  color: #7f222c;
  
}
.message .message-body {
  border: 0;
  margin-top: 4px;
}
.message p {
  font-size: 12px;
  margin: 0;
  padding: 0;
  color: #17421b;
}
.message.is-danger p {
  color: #1f2124;
}
a {
  display: block;
  margin: 8px 0;
  color: #1563ff;
  text-decoration: none;
  font-weight: 600;
}
a:hover {
  color: black;
}
a svg {
  fill: currentcolor;
}
.icon {
  align-items: center;
  display: inline-flex;
  justify-content: center;
  height: 21px;
  width: 21px;
  vertical-align: middle;
}

h1 {
  font-size: 17.5px;
  font-weight: 700;
  margin-bottom: 0;
}

h1 + p {
  margin: 8px 0 16px 0;
}
    </style>
</head>
<body translate="no" >
  <div class="container">
    <div>
      <svg id="logo" width="146" height="51" viewBox="0 0 146 51" xmlns="http://www.w3.org/2000/svg">
  <g id="vault-logo-v" fill-rule="nonzero">
    <path d="M0,0 L25.4070312,51 L51,0 L0,0 Z M28.5,10.5 L31.5,10.5 L31.5,13.5 L28.5,13.5 L28.5,10.5 Z M22.5,22.5 L19.5,22.5 L19.5,19.5 L22.5,19.5 L22.5,22.5 Z M22.5,18 L19.5,18 L19.5,15 L22.5,15 L22.5,18 Z M22.5,13.5 L19.5,13.5 L19.5,10.5 L22.5,10.5 L22.5,13.5 Z M26.991018,27 L24,27 L24,24 L27,24 L26.991018,27 Z M26.991018,22.5 L24,22.5 L24,19.5 L27,19.5 L26.991018,22.5 Z M26.991018,18 L24,18 L24,15 L27,15 L26.991018,18 Z M26.991018,13.5 L24,13.5 L24,10.5 L27,10.5 L26.991018,13.5 Z M28.5,15 L31.5,15 L31.5,18 L28.5089552,18 L28.5,15 Z M28.5,22.5 L28.5,19.5 L31.5,19.5 L31.5,22.4601182 L28.5,22.5 Z"></path>
  </g>
  <path id="vault-logo-name" d="M69.7218638,30.2482468 L63.2587814,8.45301543 L58,8.45301543 L65.9885305,34.6072931 L73.4551971,34.6072931 L81.4437276,8.45301543 L76.1849462,8.45301543 L69.7218638,30.2482468 Z M97.6329749,22.0014025 C97.6329749,17.2103787 95.8265233,15.0897616 89.6845878,15.0897616 C87.5168459,15.0897616 84.8272401,15.4431978 82.9806452,15.9929874 L83.5827957,19.6451613 C85.3089606,19.2917251 87.2358423,19.056101 89.0021505,19.056101 C92.1333333,19.056101 92.7354839,19.802244 92.7354839,21.9228612 L92.7354839,23.9256662 L88.0387097,23.9256662 C84.0645161,23.9256662 82.3383513,25.4179523 82.3383513,29.3057504 C82.3383513,32.6044881 83.8637993,35 87.4365591,35 C89.4035842,35 91.4910394,34.4502104 93.2573477,33.3113604 L93.618638,34.6072931 L97.6329749,34.6072931 L97.6329749,22.0014025 Z M92.7354839,30.2089762 C91.8121864,30.7194951 90.4874552,31.1907433 89.0422939,31.1907433 C87.5168459,31.1907433 87.0752688,30.601683 87.0752688,29.2664797 C87.0752688,27.8134642 87.5168459,27.3814867 89.1225806,27.3814867 L92.7354839,27.3814867 L92.7354839,30.2089762 Z M102.421505,15.4824684 L102.421505,29.345021 C102.421505,32.7615708 103.585663,35 106.837276,35 C109.125448,35 112.216487,34.1753156 114.665233,32.997195 L115.146953,34.6072931 L118.880287,34.6072931 L118.880287,15.4824684 L113.982796,15.4824684 L113.982796,28.7559607 C112.216487,29.6591865 110.088889,30.3660589 108.884588,30.3660589 C107.760573,30.3660589 107.318996,29.85554 107.318996,28.8345021 L107.318996,15.4824684 L102.421505,15.4824684 Z M129.168459,34.6072931 L129.168459,7 L124.270968,7.66760168 L124.270968,34.6072931 L129.168459,34.6072931 Z M144.394265,30.601683 C143.551254,30.8373072 142.6681,30.9943899 141.94552,30.9943899 C140.660932,30.9943899 140.179211,30.3267882 140.179211,29.3057504 L140.179211,19.2917251 L144.875986,19.2917251 L145.197133,15.4824684 L140.179211,15.4824684 L140.179211,10.0631136 L135.28172,10.7307153 L135.28172,15.4824684 L132.351254,15.4824684 L132.351254,19.2917251 L135.28172,19.2917251 L135.28172,29.9340813 C135.28172,33.3506311 137.088172,35 140.660932,35 C141.905376,35 143.912545,34.6858345 144.956272,34.2538569 L144.394265,30.601683 Z"></path>
</svg>
      <div class="message is-danger">
       <svg width="20" height="20" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <path d="M19 3c1.1 0 2 .9 2 2v14c0 1.1-.9 2-2 2H5c-1.1 0-2-.9-2-2V5c0-1.1.9-2 2-2h14zm-2 12.59L13.41 12 17 8.41 15.59 7 12 10.59 8.41 7 7 8.41 10.59 12 7 15.59 8.41 17 12 13.41 15.59 17 17 15.59z"></path>
</svg> 
        <div class="message-content">
          <div class="message-title">
            %s
          </div>
          <p class="message-body">
            %s
          </p>
        </div>
      </div>
      <hr />

      <h1>Not sure how to get started?</h1>
      <p class="learn">
        Check out beginner and advanced guides on HashiCorp Vault at the HashiCorp Learn site or read more in the official documentation.
      </p>
      <a href="https://learn.hashicorp.com/vault" rel="noreferrer noopener">
       <span class="icon">
      <svg width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
  <path d="M8.338 2.255a.79.79 0 0 0-.645 0L.657 5.378c-.363.162-.534.538-.534.875 0 .337.171.713.534.875l1.436.637c-.332.495-.638 1.18-.744 2.106a.887.887 0 0 0-.26 1.559c.02.081.03.215.013.392-.02.205-.074.43-.162.636-.186.431-.45.64-.741.64v.98c.651 0 1.108-.365 1.403-.797l.06.073c.32.372.826.763 1.455.763v-.98c-.215 0-.474-.145-.71-.42-.111-.13-.2-.27-.259-.393a1.014 1.014 0 0 1-.06-.155c-.01-.036-.013-.055-.013-.058h-.022a2.544 2.544 0 0 0 .031-.641.886.886 0 0 0-.006-1.51c.1-.868.398-1.477.699-1.891l.332.147-.023.746v2.228c0 .115.04.22.105.304.124.276.343.5.587.677.297.217.675.396 1.097.54.846.288 1.943.456 3.127.456 1.185 0 2.281-.168 3.128-.456.422-.144.8-.323 1.097-.54.244-.177.462-.401.586-.677a.488.488 0 0 0 .106-.304V8.218l2.455-1.09c.363-.162.534-.538.534-.875 0-.337-.17-.713-.534-.875L8.338 2.255zm-.34 2.955L3.64 7.38l4.375 1.942 6.912-3.069-6.912-3.07-6.912 3.07 1.665.74 4.901-2.44.328.657zM14.307 1H12.5a.5.5 0 1 1 0-1h3a.499.499 0 0 1 .5.65V3.5a.5.5 0 1 1-1 0V1.72l-1.793 1.774a.5.5 0 0 1-.713-.701L14.307 1zm-2.368 7.653v2.383a.436.436 0 0 0-.007.021c-.017.063-.084.178-.282.322-.193.14-.473.28-.836.404-.724.247-1.71.404-2.812.404-1.1 0-2.087-.157-2.811-.404a3.188 3.188 0 0 1-.836-.404c-.198-.144-.265-.26-.282-.322a.437.437 0 0 0-.007-.02V8.983l.01-.338 3.617 1.605a.791.791 0 0 0 .645 0l3.6-1.598z" fill-rule="evenodd"></path>
</svg>
        </span>
  Get started with Vault
      </a>

      
      <a href="https://vaultproject.io/docs" rel="noreferrer noopener">
       <span class="icon">
        <svg width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
  <path d="M13.307 1H11.5a.5.5 0 1 1 0-1h3a.499.499 0 0 1 .5.65V3.5a.5.5 0 1 1-1 0V1.72l-1.793 1.774a.5.5 0 0 1-.713-.701L13.307 1zM12 14V8a.5.5 0 1 1 1 0v6.5a.5.5 0 0 1-.5.5H.563a.5.5 0 0 1-.5-.5v-13a.5.5 0 0 1 .5-.5H8a.5.5 0 0 1 0 1H1v12h11zM4 6a.5.5 0 0 1 0-1h3a.5.5 0 0 1 0 1H4zm0 2.5a.5.5 0 0 1 0-1h5a.5.5 0 0 1 0 1H4zM4 11a.5.5 0 1 1 0-1h5a.5.5 0 1 1 0 1H4z"/>
</svg> 
        </span>
        View the official Vault documentation
      </a>
    </div>
  </div>
</body>

</html>
`
	return fmt.Sprintf(html, summary, detail)
}

// Help method for OIDC cli
func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=oidc [CONFIG K=V...]

  The OIDC auth method allows users to authenticate using an OIDC provider.
  The provider must be configured as part of a role by the operator.

  Authenticate using role "engineering":

      $ vault login -method=oidc role=engineering
      Complete the login via your OIDC provider. Launching browser to:

          https://accounts.google.com/o/oauth2/v2/...
   
  The default browser will be opened for the user to complete the login. Alternatively,
  the user may visit the provided URL directly.

Configuration:

  role=<string>
      Vault role of type "OIDC" to use for authentication.

  port=<string>
      Optional localhost port to use for OIDC callback (default: 8300).
`

	return strings.TrimSpace(help)
}

const successHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Vault Authentication Succeeded</title>
    <style>
      body {
        font-size: 14px;
        font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
          "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans",
          "Helvetica Neue", sans-serif;
      }
      hr {
        border-color: #fdfdfe;
        margin: 24px 0;
      }
      .container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 70vh;
      }
      #logo {
        display: block;
        fill: #6f7682;
        margin-bottom: 16px;
      }
      .message {
        display: flex;
        min-width: 40vw;
        background: #fafdfa;
        border: 1px solid #c6e9c9;
        margin-bottom: 12px;
        padding: 12px 16px 16px 12px;
        position: relative;
        border-radius: 2px;
        font-size: 14px;
      }
      .message-content {
        margin-left: 4px;
      }
      .message #checkbox {
        fill: #2eb039;
      }
      .message .message-title {
        color: #1e7125;
        font-size: 16px;
        font-weight: 700;
        line-height: 1.25;
      }
      .message .message-body {
        border: 0;
        margin-top: 4px;
      }
      .message p {
        font-size: 12px;
        margin: 0;
        padding: 0;
        color: #17421b;
      }
      a {
        display: block;
        margin: 8px 0;
        color: #1563ff;
        text-decoration: none;
        font-weight: 600;
      }
      a:hover {
        color: black;
      }
      a svg {
        fill: currentcolor;
      }
      .icon {
        align-items: center;
        display: inline-flex;
        justify-content: center;
        height: 21px;
        width: 21px;
        vertical-align: middle;
      }
      h1 {
        font-size: 17.5px;
        font-weight: 700;
        margin-bottom: 0;
      }
      h1 + p {
        margin: 8px 0 16px 0;
      }
    </style>
  </head>
  <body translate="no" >
    <div class="container">
      <div>
        <svg id="logo" width="146" height="51" viewBox="0 0 146 51" xmlns="http://www.w3.org/2000/svg">
          <g id="vault-logo-v" fill-rule="nonzero">
            <path d="M0,0 L25.4070312,51 L51,0 L0,0 Z M28.5,10.5 L31.5,10.5 L31.5,13.5 L28.5,13.5 L28.5,10.5 Z M22.5,22.5 L19.5,22.5 L19.5,19.5 L22.5,19.5 L22.5,22.5 Z M22.5,18 L19.5,18 L19.5,15 L22.5,15 L22.5,18 Z M22.5,13.5 L19.5,13.5 L19.5,10.5 L22.5,10.5 L22.5,13.5 Z M26.991018,27 L24,27 L24,24 L27,24 L26.991018,27 Z M26.991018,22.5 L24,22.5 L24,19.5 L27,19.5 L26.991018,22.5 Z M26.991018,18 L24,18 L24,15 L27,15 L26.991018,18 Z M26.991018,13.5 L24,13.5 L24,10.5 L27,10.5 L26.991018,13.5 Z M28.5,15 L31.5,15 L31.5,18 L28.5089552,18 L28.5,15 Z M28.5,22.5 L28.5,19.5 L31.5,19.5 L31.5,22.4601182 L28.5,22.5 Z"></path>
          </g>
          <path id="vault-logo-name" d="M69.7218638,30.2482468 L63.2587814,8.45301543 L58,8.45301543 L65.9885305,34.6072931 L73.4551971,34.6072931 L81.4437276,8.45301543 L76.1849462,8.45301543 L69.7218638,30.2482468 Z M97.6329749,22.0014025 C97.6329749,17.2103787 95.8265233,15.0897616 89.6845878,15.0897616 C87.5168459,15.0897616 84.8272401,15.4431978 82.9806452,15.9929874 L83.5827957,19.6451613 C85.3089606,19.2917251 87.2358423,19.056101 89.0021505,19.056101 C92.1333333,19.056101 92.7354839,19.802244 92.7354839,21.9228612 L92.7354839,23.9256662 L88.0387097,23.9256662 C84.0645161,23.9256662 82.3383513,25.4179523 82.3383513,29.3057504 C82.3383513,32.6044881 83.8637993,35 87.4365591,35 C89.4035842,35 91.4910394,34.4502104 93.2573477,33.3113604 L93.618638,34.6072931 L97.6329749,34.6072931 L97.6329749,22.0014025 Z M92.7354839,30.2089762 C91.8121864,30.7194951 90.4874552,31.1907433 89.0422939,31.1907433 C87.5168459,31.1907433 87.0752688,30.601683 87.0752688,29.2664797 C87.0752688,27.8134642 87.5168459,27.3814867 89.1225806,27.3814867 L92.7354839,27.3814867 L92.7354839,30.2089762 Z M102.421505,15.4824684 L102.421505,29.345021 C102.421505,32.7615708 103.585663,35 106.837276,35 C109.125448,35 112.216487,34.1753156 114.665233,32.997195 L115.146953,34.6072931 L118.880287,34.6072931 L118.880287,15.4824684 L113.982796,15.4824684 L113.982796,28.7559607 C112.216487,29.6591865 110.088889,30.3660589 108.884588,30.3660589 C107.760573,30.3660589 107.318996,29.85554 107.318996,28.8345021 L107.318996,15.4824684 L102.421505,15.4824684 Z M129.168459,34.6072931 L129.168459,7 L124.270968,7.66760168 L124.270968,34.6072931 L129.168459,34.6072931 Z M144.394265,30.601683 C143.551254,30.8373072 142.6681,30.9943899 141.94552,30.9943899 C140.660932,30.9943899 140.179211,30.3267882 140.179211,29.3057504 L140.179211,19.2917251 L144.875986,19.2917251 L145.197133,15.4824684 L140.179211,15.4824684 L140.179211,10.0631136 L135.28172,10.7307153 L135.28172,15.4824684 L132.351254,15.4824684 L132.351254,19.2917251 L135.28172,19.2917251 L135.28172,29.9340813 C135.28172,33.3506311 137.088172,35 140.660932,35 C141.905376,35 143.912545,34.6858345 144.956272,34.2538569 L144.394265,30.601683 Z"></path>
        </svg>
        <div class="message is-success">
          <svg id="checkbox" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 512 512">
            <path d="M256 32C132.3 32 32 132.3 32 256s100.3 224 224 224 224-100.3 224-224S379.7 32 256 32zm114.9 149.1L231.8 359.6c-1.1 1.1-2.9 3.5-5.1 3.5-2.3 0-3.8-1.6-5.1-2.9-1.3-1.3-78.9-75.9-78.9-75.9l-1.5-1.5c-.6-.9-1.1-2-1.1-3.2 0-1.2.5-2.3 1.1-3.2.4-.4.7-.7 1.1-1.2 7.7-8.1 23.3-24.5 24.3-25.5 1.3-1.3 2.4-3 4.8-3 2.5 0 4.1 2.1 5.3 3.3 1.2 1.2 45 43.3 45 43.3l111.3-143c1-.8 2.2-1.4 3.5-1.4 1.3 0 2.5.5 3.5 1.3l30.6 24.1c.8 1 1.3 2.2 1.3 3.5.1 1.3-.4 2.4-1 3.3z"></path>
        </svg>
          <div class="message-content">
            <div class="message-title">
              Signed in via your OIDC provider
            </div>
            <p class="message-body">
              You can now close this window and start using Vault.
            </p>
          </div>
        </div>
        <hr />
        <h1>Not sure how to get started?</h1>
        <p class="learn">
          Check out beginner and advanced guides on HashiCorp Vault at the HashiCorp Learn site or read more in the official documentation.
        </p>
        <a href="https://learn.hashicorp.com/vault" rel="noreferrer noopener">
         <span class="icon">
            <svg width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
              <path d="M8.338 2.255a.79.79 0 0 0-.645 0L.657 5.378c-.363.162-.534.538-.534.875 0 .337.171.713.534.875l1.436.637c-.332.495-.638 1.18-.744 2.106a.887.887 0 0 0-.26 1.559c.02.081.03.215.013.392-.02.205-.074.43-.162.636-.186.431-.45.64-.741.64v.98c.651 0 1.108-.365 1.403-.797l.06.073c.32.372.826.763 1.455.763v-.98c-.215 0-.474-.145-.71-.42-.111-.13-.2-.27-.259-.393a1.014 1.014 0 0 1-.06-.155c-.01-.036-.013-.055-.013-.058h-.022a2.544 2.544 0 0 0 .031-.641.886.886 0 0 0-.006-1.51c.1-.868.398-1.477.699-1.891l.332.147-.023.746v2.228c0 .115.04.22.105.304.124.276.343.5.587.677.297.217.675.396 1.097.54.846.288 1.943.456 3.127.456 1.185 0 2.281-.168 3.128-.456.422-.144.8-.323 1.097-.54.244-.177.462-.401.586-.677a.488.488 0 0 0 .106-.304V8.218l2.455-1.09c.363-.162.534-.538.534-.875 0-.337-.17-.713-.534-.875L8.338 2.255zm-.34 2.955L3.64 7.38l4.375 1.942 6.912-3.069-6.912-3.07-6.912 3.07 1.665.74 4.901-2.44.328.657zM14.307 1H12.5a.5.5 0 1 1 0-1h3a.499.499 0 0 1 .5.65V3.5a.5.5 0 1 1-1 0V1.72l-1.793 1.774a.5.5 0 0 1-.713-.701L14.307 1zm-2.368 7.653v2.383a.436.436 0 0 0-.007.021c-.017.063-.084.178-.282.322-.193.14-.473.28-.836.404-.724.247-1.71.404-2.812.404-1.1 0-2.087-.157-2.811-.404a3.188 3.188 0 0 1-.836-.404c-.198-.144-.265-.26-.282-.322a.437.437 0 0 0-.007-.02V8.983l.01-.338 3.617 1.605a.791.791 0 0 0 .645 0l3.6-1.598z" fill-rule="evenodd"></path>
            </svg>
          </span>
          Get started with Vault
        </a>
        <a href="https://vaultproject.io/docs" rel="noreferrer noopener">
         <span class="icon">
          <svg width="16" height="16" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
    <path d="M13.307 1H11.5a.5.5 0 1 1 0-1h3a.499.499 0 0 1 .5.65V3.5a.5.5 0 1 1-1 0V1.72l-1.793 1.774a.5.5 0 0 1-.713-.701L13.307 1zM12 14V8a.5.5 0 1 1 1 0v6.5a.5.5 0 0 1-.5.5H.563a.5.5 0 0 1-.5-.5v-13a.5.5 0 0 1 .5-.5H8a.5.5 0 0 1 0 1H1v12h11zM4 6a.5.5 0 0 1 0-1h3a.5.5 0 0 1 0 1H4zm0 2.5a.5.5 0 0 1 0-1h5a.5.5 0 0 1 0 1H4zM4 11a.5.5 0 1 1 0-1h5a.5.5 0 1 1 0 1H4z"/>
  </svg> 
          </span>
          View the official Vault documentation
        </a>
      </div>
    </div>
  </body>
</html>
`
