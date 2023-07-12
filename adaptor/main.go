package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/romana/rlog"
	"github.com/valyala/fasthttp"
)

type BitBucketUser struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Links        struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type BitBucketReviewers []struct {
	User     BitBucketUser `json:"user"`
	Role     string        `json:"role"`
	Approved bool          `json:"approved"`
	Status   string        `json:"status"`
}

type BitBucketPREvent struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		ID           int    `json:"id"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"actor"`
	PullRequest struct {
		ID          int    `json:"id"`
		Version     int    `json:"version"`
		Title       string `json:"title"`
		State       string `json:"state"`
		Open        bool   `json:"open"`
		Closed      bool   `json:"closed"`
		CreatedDate int64  `json:"createdDate"`
		UpdatedDate int64  `json:"updatedDate"`
		FromRef     struct {
			ID           string `json:"id"`
			DisplayID    string `json:"displayId"`
			LatestCommit string `json:"latestCommit"`
			Type         string `json:"type"`
			Repository   struct {
				Slug          string `json:"slug"`
				ID            int    `json:"id"`
				Name          string `json:"name"`
				HierarchyID   string `json:"hierarchyId"`
				ScmID         string `json:"scmId"`
				State         string `json:"state"`
				StatusMessage string `json:"statusMessage"`
				Forkable      bool   `json:"forkable"`
				Project       struct {
					Key         string `json:"key"`
					ID          int    `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
					Public      bool   `json:"public"`
					Type        string `json:"type"`
					Links       struct {
						Self []struct {
							Href string `json:"href"`
						} `json:"self"`
					} `json:"links"`
				} `json:"project"`
				Public bool `json:"public"`
				Links  struct {
					Clone []struct {
						Href string `json:"href"`
						Name string `json:"name"`
					} `json:"clone"`
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"repository"`
		} `json:"fromRef"`
		ToRef struct {
			ID           string `json:"id"`
			DisplayID    string `json:"displayId"`
			LatestCommit string `json:"latestCommit"`
			Type         string `json:"type"`
			Repository   struct {
				Slug          string `json:"slug"`
				ID            int    `json:"id"`
				Name          string `json:"name"`
				HierarchyID   string `json:"hierarchyId"`
				ScmID         string `json:"scmId"`
				State         string `json:"state"`
				StatusMessage string `json:"statusMessage"`
				Forkable      bool   `json:"forkable"`
				Project       struct {
					Key         string `json:"key"`
					ID          int    `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
					Public      bool   `json:"public"`
					Type        string `json:"type"`
					Links       struct {
						Self []struct {
							Href string `json:"href"`
						} `json:"self"`
					} `json:"links"`
				} `json:"project"`
				Public bool `json:"public"`
				Links  struct {
					Clone []struct {
						Href string `json:"href"`
						Name string `json:"name"`
					} `json:"clone"`
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"repository"`
		} `json:"toRef"`
		Locked bool `json:"locked"`
		Author struct {
			User struct {
				Name         string `json:"name"`
				EmailAddress string `json:"emailAddress"`
				ID           int    `json:"id"`
				DisplayName  string `json:"displayName"`
				Active       bool   `json:"active"`
				Slug         string `json:"slug"`
				Type         string `json:"type"`
				Links        struct {
					Self []struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"links"`
			} `json:"user"`
			Role     string `json:"role"`
			Approved bool   `json:"approved"`
			Status   string `json:"status"`
		} `json:"author"`
		Reviewers    BitBucketReviewers `json:"reviewers"`
		Participants []any              `json:"participants"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"pullRequest"`
}

type ReviewerEntity struct {
	Type      string `default:"mention" json:"type"`
	Text      string `json:"text"`
	Mentioned struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"mentioned"`
}

type ReviewerEntitiesList []ReviewerEntity

type TeamsMsgBody struct {
	Type   string `default:"TextBlock" json:"type"`
	Size   string `default:"Medium" json:"size,omitempty"`
	Weight string `default:"Bolder" json:"weight,omitempty"`
	Text   string `default:"Webhook Connector" json:"text"`
	Wrap   bool   `default:"true" json:"wrap"`
}

type TeamsMsgAttachement struct {
	ContentType string `default:"application/vnd.microsoft.card.adaptive" json:"contentType"`
	Content     struct {
		Type    string         `default:"AdaptiveCard" json:"type"`
		Body    []TeamsMsgBody `json:"body"`
		Schema  string         `default:"http://adaptivecards.io/schemas/adaptive-card.json" json:"$schema"`
		Version string         `default:"1.0" json:"version"`
		Msteams struct {
			Width    string               `default:"Full" json:"width"`
			Entities ReviewerEntitiesList `json:"entities"`
		} `json:"msteams"`
	} `json:"content"`
}

type TeamsMsg struct {
	Type        string                `default:"message" json:"type"`
	Attachments []TeamsMsgAttachement `json:"attachments"`
}

// Produce JSON with <> not escaped as unicode
func (t *TeamsMsg) NonEscapedJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Parse BitBucket PR event json payload, maps data to build Teams notification webhook json
func ParsePR(eventJson []byte) ([]byte, error) {
	var inventory BitBucketPREvent
	if err := json.Unmarshal([]byte(eventJson), &inventory); err != nil {
		errMsg := fmt.Sprintf("Error Unmarshalling payload JSON : %s", err.Error())
		rlog.Error(errMsg)
		return []byte(""), errors.New(errMsg)
	}
	rlog.Tracef(0, "inventory : %+v\n", inventory)
	var reviewersList string = ""
	var reviewersEntity ReviewerEntity
	var reviewersEntityList ReviewerEntitiesList
	reviewersEntity.Type = "mention"
	if len(inventory.PullRequest.Reviewers) == 0 {
		rlog.Errorf("Reviewers count is 0")
	} else {
		for _, val := range inventory.PullRequest.Reviewers {
			text := "<at>" + val.User.Name + " UPN</at>"
			reviewersList += text + ", "
			reviewersEntity.Text = text
			reviewersEntity.Mentioned.ID = val.User.EmailAddress
			reviewersEntity.Mentioned.Name = val.User.DisplayName
			reviewersEntityList = append(reviewersEntityList, reviewersEntity)
		}
	}
	reviewersEntity.Mentioned.ID = inventory.PullRequest.Author.User.EmailAddress
	reviewersEntity.Mentioned.Name = inventory.PullRequest.Author.User.DisplayName
	reviewersEntity.Text = "<at>" + inventory.PullRequest.Author.User.Name + " UPN</at>"
	reviewersEntityList = append(reviewersEntityList, reviewersEntity) // add PR author to mentions format
	reviewersList = strings.TrimRight(reviewersList, ", ")
	rlog.Tracef(0, "reviewersEntityList : %+v\n", reviewersEntityList)

	var prAction string
	switch strings.TrimLeft(inventory.EventKey, "pr:") {
	case "opened":
		prAction = "opened a PR"
	case "from_ref_updated":
		prAction = "updated source branch in PR"
	default:
		prAction = "(PR EventKey: " + inventory.EventKey + ")"
	}

	bodyText := fmt.Sprintf("Hi Team, %s %s, please review: [%s](%s) \n\n", reviewersEntity.Text, prAction, inventory.PullRequest.Title, inventory.PullRequest.Links.Self[0].Href)
	bodyText += fmt.Sprintf("CC: %s", reviewersList)
	rlog.Tracef(0, "bodyText : %s \n", bodyText)

	var msg TeamsMsg
	msg.Type = "message"

	var msgAttachement TeamsMsgAttachement
	msgAttachement.ContentType = "application/vnd.microsoft.card.adaptive"
	msgAttachement.Content.Type = "AdaptiveCard"

	var msgBodyList []TeamsMsgBody
	var msgBody TeamsMsgBody
	msgBody.Type = "TextBlock"
	msgBody.Text = bodyText
	msgBody.Wrap = true
	msgBodyList = append(msgBodyList, msgBody)

	msgAttachement.Content.Body = msgBodyList
	msgAttachement.Content.Schema = "http://adaptivecards.io/schemas/adaptive-card.json"
	msgAttachement.Content.Version = "1.0"
	msgAttachement.Content.Msteams.Width = "Full"
	msgAttachement.Content.Msteams.Entities = reviewersEntityList

	msg.Attachments = append(msg.Attachments, msgAttachement)

	b, err := msg.NonEscapedJSON()
	if err != nil {
		rlog.Errorf("NonEscapedJSON error: %s", err.Error())
	}
	return b, nil
}

func isTraceLevel(tLevel int64) bool {
	return tLevel >= 0
}

const (
	levelNone = iota
	levelCrit
	levelErr
	levelWarn
	levelInfo
	levelDebug
	levelTrace
)

func main() {
	// override fiber encoder/decoder with one provided by goccy/go-json
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	os.Setenv("RLOG_LOG_STREAM", "stdout")
	rlog.UpdateEnv()
	var logLevel string = os.Getenv("RLOG_LOG_LEVEL")
	var traceLevelEnv string = os.Getenv("RLOG_TRACE_LEVEL")
	var traceLevel int64
	var httpScheme string = os.Getenv("HTTP_SCHEME") // http (for local development) or https, default is https
	if httpScheme == "" {
		httpScheme = "https"
	}
	var envTlsInsecureSkipVerify = os.Getenv("TLS_INSECURE_SKIP_VERIFY") // for testing we can ignore check for self-signed CA cert
	var tlsInsecureSkipVerify bool = false
	if envTlsInsecureSkipVerify != "" {
		var parseBoolErr error
		tlsInsecureSkipVerify, parseBoolErr = strconv.ParseBool(envTlsInsecureSkipVerify)
		if parseBoolErr != nil {
			rlog.Criticalf("Not a Boolean value in envvar TLS_INSECURE_SKIP_VERIFY: %s ; Error: %s", envTlsInsecureSkipVerify, parseBoolErr)
			os.Exit(1)
		}
	}
	var teamsHost string = os.Getenv("TEAMS_HOSTNAME") // somecorp.webhook.office.com
	if teamsHost == "" {
		rlog.Critical("Mandatory environment variable TEAMS_HOSTNAME (FQDN from webhook) is not set. You can set it as localhost for development, exiting")
		os.Exit(1)
	} else {
		rlog.Info("TEAMS_HOSTNAME: ", teamsHost)
	}
	if logLevel == "" {
		logLevel = "INFO"
	}
	// If this variable is undefined, or set to -1 then no Trace messages are printed :
	if traceLevelEnv == "" {
		traceLevel = -1
	} else {
		x, err := strconv.ParseInt(traceLevelEnv, 10, 64)
		if err != nil {
			rlog.Criticalf("RLOG_TRACE_LEVEL value provided is not int64 type. Error : %s", err.Error())
			os.Exit(1)
		} else {
			traceLevel = x
		}
	}
	rlog.Infof("RLOG_LOG_LEVEL: %s; RLOG_TRACE_LEVEL: %d; TLS_INSECURE_SKIP_VERIFY: %t", logLevel, traceLevel, tlsInsecureSkipVerify)

	app.Use(requestid.New(requestid.Config{
		Next:       nil,
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUIDv4,
		ContextKey: "requestid",
	}))

	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		Format:     "${time} ACCESS   : [${ip}]:${port} ${locals:requestid} ${status} - ${latency} ${bytesReceived} ${method} ${path}\n",
	}))

	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/healthz"
		},
		Level: compress.LevelBestSpeed, // 1
	}))

	type SomeStruct struct {
		RequestID string
	}

	// GET /healthz
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// POST /webhookb2/uid1@uid2/IncomingWebhook/uid3/uid4
	app.Post("/webhookb2/:id1/IncomingWebhook/:id2/:id3", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"
		c.AcceptsEncodings("compress", "br")
		data := SomeStruct{
			RequestID: c.GetRespHeader("X-Request-Id"),
		}
		rlog.Debugf("X-Request-Id : %s", data.RequestID)
		pathid1 := c.Params("id1")
		pathid2 := c.Params("id2")
		pathid3 := c.Params("id3")
		rlog.Debugf("hook ids: %s, %s, %s ; body: %s \n", pathid1, pathid2, pathid3, c.Body())

		var newPath string = ""
		if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
			// https://docs.gofiber.io/api/ctx#path :
			// override Path with sha256 encoded webhook credentials
			id1 := fmt.Sprintf("%x", sha256.Sum256([]byte(pathid1)))
			id2 := fmt.Sprintf("%x", sha256.Sum256([]byte(pathid2)))
			id3 := fmt.Sprintf("%x", sha256.Sum256([]byte(pathid3)))
			newPath = fmt.Sprintf("/webhookb2/%s/IncomingWebhook/%s/%s", id1[0:7], id2[0:7], id3[0:7])
		}

		// send request to teams , curl -v -X POST -H 'Content-Type: application/json' 'https://somecorp.webhook.office.com/webhookb2/
		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{},
			InsecureSkipVerify: tlsInsecureSkipVerify,
		}
		teamsURI := fmt.Sprintf("%s://%s/webhookb2/%s/IncomingWebhook/%s/%s", httpScheme, teamsHost, pathid1, pathid2, pathid3)
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(teamsURI)
		// don't parse further if body don't exist or empty string : without -d or curl -d ''
		if (c.Body()) == nil {
			errMsg := "Request Body is nil"
			rlog.Debug(errMsg)
			c.Set("Content-Type", "text/plain; charset=utf-8")
			if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
				c.Path(newPath) // override to not log sensitive webhook parts
			}
			return c.Status(400).SendString("Error: " + errMsg)
		} else if bytes.Equal(c.Body(), []byte("")) {
			errMsg := "Request Body is empty"
			rlog.Debug(errMsg)
			c.Set("Content-Type", "text/plain; charset=utf-8")
			if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
				c.Path(newPath) // override to not log sensitive webhook parts
			}
			return c.Status(400).SendString("Error: " + errMsg)
		}

		var notificationBody []byte
		// TODO : parse JSON and check jsonpath : .test=true
		if bytes.Equal(c.Body(), []byte("{\"test\": true}")) {
			rlog.Debug("Request was Test ping ")
			notificationBody = c.Body()
			c.Set("Content-Type", "text/plain; charset=utf-8")
			if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
				c.Path(newPath) // override to not log sensitive webhook parts
			}
			return c.Status(200).SendString("ok")
		} else {
			var parseErr error
			notificationBody, parseErr = ParsePR(c.Body())
			if parseErr == nil {
				rlog.Debugf("notificationBody : %s", notificationBody)
				c.Set("Content-Type", "application/json")
			} else {
				errMsg := fmt.Sprintf("JSON parsing error was: %s", parseErr.Error())
				rlog.Error(errMsg)
				c.Set("Content-Type", "text/plain; charset=utf-8")
				if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
					c.Path(newPath) // override to not log sensitive webhook parts
				}
				return c.Status(400).SendString("Error: " + errMsg)
			}
		}

		req.Header.SetMethod("POST")
		req.Header.Add("X-Request-Id", data.RequestID)
		req.Header.Set("Content-Type", "application/json")
		req.SetBody(notificationBody)
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)
		client := &fasthttp.Client{
			TLSConfig: tlsConfig,
		}
		errs := client.Do(req, resp) // sending request to Teams host
		code := resp.StatusCode()
		body := resp.Body()
		var respHeader fasthttp.ResponseHeader
		resp.Header.CopyTo(&respHeader)

		// moved after RequestURI evaluated and sent because pathid1 was changing after changing Path :
		if (logLevel != "DEBUG") && !(isTraceLevel(traceLevel)) {
			// https://docs.gofiber.io/api/ctx#path :
			// override Path with sha256 encoded webhook credentials
			// newPath concatenation moved earlier to meet early exit path logging adjustments
			c.Path(newPath) // override to not log sensitive webhook parts
		}

		rlog.Infof("Notification sent to Teams, request Id: %s ; result code:%d", data.RequestID, code)
		if code >= 400 {
			errMsg := fmt.Sprintf("Teams API request (%s) failed with HTTP code: %d", data.RequestID, code)
			rlog.Error(errMsg)
			c.Set("Content-Type", "text/plain; charset=utf-8")
			return c.Status(code).SendString("Error: " + errMsg)
		}
		rlog.Debugf("Notification response body: %s", body)
		respContType := respHeader.ContentType()
		if respContType != nil {
			rlog.Debugf("Notification response header contentType: %s", respContType)
		}
		respEnc := respHeader.ContentEncoding()
		if respEnc != nil {
			rlog.Debugf("Notification response header contentEncoding: %s", respEnc)
		}
		if errs != nil {
			errMsg := fmt.Sprintf("Teams API request (%s) reported error: %s \n", data.RequestID, errs.Error())
			rlog.Error(errMsg)
			c.Set("Content-Type", "text/plain; charset=utf-8")
			return c.Status(504).SendString("Error: " + errMsg)
		}
		return c.Send(body)
	})

	go func() {
		err := app.Listen(":8080")
		if err != nil {
			rlog.Criticalf("Listener on port 8080 error: %s", err.Error())
			os.Exit(1)
		}
	}()

	appHealth := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	// GET /healthz
	appHealth.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})
	errHealthz := appHealth.Listen(":9000")
	if errHealthz != nil {
		rlog.Criticalf("Listener on port 9000 error: %s", errHealthz.Error())
		os.Exit(1)
	}

}
