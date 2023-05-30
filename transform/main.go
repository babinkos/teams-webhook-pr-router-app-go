package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"
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

// type BitBucketReviewers []struct {
// 	User BitBucketUser
// }

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

func (t *TeamsMsg) NonEscapedJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func ParsePR(eventJson []byte) TeamsMsg {
	// value := "stub"
	var inventory BitBucketPREvent
	if err := json.Unmarshal([]byte(eventJson), &inventory); err != nil {
		log.Fatal(err)
	}
	var reviewersList string = ""
	var reviewersEntity ReviewerEntity
	var reviewersEntityList ReviewerEntitiesList
	reviewersEntity.Type = "mention"
	for _, val := range inventory.PullRequest.Reviewers {
		text := "<at>" + val.User.Name + " UPN</at>"
		reviewersList += text + ", "
		reviewersEntity.Text = text
		reviewersEntity.Mentioned.ID = val.User.EmailAddress
		reviewersEntity.Mentioned.Name = val.User.DisplayName
		reviewersEntityList = append(reviewersEntityList, reviewersEntity)
	}
	reviewersEntity.Mentioned.ID = inventory.PullRequest.Author.User.EmailAddress
	reviewersEntity.Mentioned.Name = inventory.PullRequest.Author.User.DisplayName
	reviewersEntity.Text = "<at>" + inventory.PullRequest.Author.User.Name + " UPN</at>"
	reviewersEntityList = append(reviewersEntityList, reviewersEntity) // add PR author to mentions format
	reviewersList = strings.TrimRight(reviewersList, ", ")
	// fmt.Printf("%+v\n", reviewersEntityList)
	bodyText := fmt.Sprintf("Hi Team, %s %s a PR, please review: <a href=\"%s\">%s</a> \n", reviewersEntity.Text, strings.TrimLeft(inventory.EventKey, "pr:"), inventory.PullRequest.Links.Self[0].Href, inventory.PullRequest.Title)
	bodyText += fmt.Sprintf("CC: %s", reviewersList)
	fmt.Printf("%s\n", bodyText)

	var msg TeamsMsg
	msg.Type = "message"

	var msgAttachement TeamsMsgAttachement
	msgAttachement.ContentType = "application/vnd.microsoft.card.adaptive"
	msgAttachement.Content.Type = "AdaptiveCard"

	var msgBodyList []TeamsMsgBody
	var msgBody TeamsMsgBody
	msgBody.Type = "TextBlock"
	msgBody.Size = "Medium"
	msgBody.Weight = "Bolder"
	msgBody.Text = bodyText
	msgBodyList = append(msgBodyList, msgBody)

	msgAttachement.Content.Body = msgBodyList
	msgAttachement.Content.Schema = "http://adaptivecards.io/schemas/adaptive-card.json"
	msgAttachement.Content.Version = "1.0"
	msgAttachement.Content.Msteams.Width = "Full"
	msgAttachement.Content.Msteams.Entities = reviewersEntityList

	msg.Attachments = append(msg.Attachments, msgAttachement)

	b, err := msg.NonEscapedJSON()
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	return msg
}

func main() {

	ParsePR([]byte(`{"eventKey":"pr:opened","date":"2023-05-29T09:19:13-0400","actor":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"pullRequest":{"id":269,"version":0,"title":"PMP-20074 Dataseeding request for azfdemo org - test notification","state":"OPEN","open":true,"closed":false,"createdDate":1685366353203,"updatedDate":1685366353203,"fromRef":{"id":"refs/heads/feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","displayId":"feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","latestCommit":"9e4da367c66845f212567e839743694f6e8439f6","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"toRef":{"id":"refs/heads/develop","displayId":"develop","latestCommit":"b496d7bc5618fb20faeed7f9f4559fdde4d42adb","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"locked":false,"author":{"user":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"role":"AUTHOR","approved":false,"status":"UNAPPROVED"},"reviewers":[{"user":{"name":"svaladez","emailAddress":"svaladez@lenovo.com","id":1454,"displayName":"Sergio Valadez Cruz","active":true,"slug":"svaladez","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/svaladez"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"},{"user":{"name":"rthoomu1","emailAddress":"rthoomu1@lenovo.com","id":8851,"displayName":"Ramesh Babu Thoomu1","active":true,"slug":"rthoomu1","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/rthoomu1"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"}],"participants":[],"links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/pull-requests/269"}]}}}`))
	ParsePR([]byte(`{"eventKey":"pr:modified","date":"2023-05-29T09:19:13-0401","actor":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"pullRequest":{"id":269,"version":0,"title":"DEV-9080 K8s update in Pen - test notification","state":"OPEN","open":true,"closed":false,"createdDate":1685366353203,"updatedDate":1685366353203,"fromRef":{"id":"refs/heads/feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","displayId":"feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","latestCommit":"9e4da367c66845f212567e839743694f6e8439f6","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"toRef":{"id":"refs/heads/develop","displayId":"develop","latestCommit":"b496d7bc5618fb20faeed7f9f4559fdde4d42adb","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"locked":false,"author":{"user":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"role":"AUTHOR","approved":false,"status":"UNAPPROVED"},"reviewers":[{"user":{"name":"svaladez","emailAddress":"svaladez@lenovo.com","id":1454,"displayName":"Sergio Valadez Cruz","active":true,"slug":"svaladez","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/svaladez"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"},{"user":{"name":"rthoomu1","emailAddress":"rthoomu1@lenovo.com","id":8851,"displayName":"Ramesh Babu Thoomu1","active":true,"slug":"rthoomu1","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/rthoomu1"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"}],"participants":[],"links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/pull-requests/270"}]}}}`))
}
