package request

import (
	"github.com/fuyoumingyan/gofinger/pkg/utils"
	"github.com/projectdiscovery/gologger"
	"regexp"
	"strings"
)

func GetJSRedirectURL(baseUrl, body string) string {
	rules := []string{
		`window\.location\.href\("([^"]+)"\)`,
		`window\.location\.replace\("([^"]+)"\)`,
		`location\.replace\("([^"]+)"\)`,
		`<meta http-equiv="refresh".*url=(.*)"`,
		`(?im)\s*content=['"]\d;url=['"](.*?)['"]`,
	}
	var redirectURL string
	for _, rule := range rules {
		re := regexp.MustCompile(rule)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 1 {
			redirectURL = matches[1]
		}
	}
	if redirectURL != "" && !strings.HasPrefix(redirectURL, "http") {
		uri, err := utils.JoinURL(baseUrl, redirectURL)
		if err != nil {
			gologger.Error().Msg(err.Error())
			return ""
		}
		redirectURL = uri
	}
	return redirectURL
}
