package main

import (
	"fmt"
	"log"

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

type BitBucketReviewers []struct {
	User BitBucketUser
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
		Reviewers []struct {
			User BitBucketUser
			} `json:"user"`
			Role     string `json:"role"`
			Approved bool   `json:"approved"`
			Status   string `json:"status"`
		} `json:"reviewers"`
		Participants []any `json:"participants"`
		Links        struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
	} `json:"pullRequest"`
}

func ParsePR(eventJson []byte) string {
	// value := "stub"
	var inventory BitBucketPREvent
	if err := json.Unmarshal([]byte(eventJson), &inventory); err != nil {
		log.Fatal(err)
	}
	var reviewersList 
	return string(fmt.Sprintf("Author: %s ; PR: %s; Reviewers Count: %d", inventory.Actor.DisplayName, inventory.PullRequest.Links.Self[0].Href, len(inventory.PullRequest.Reviewers)))
}

func main() {

	fmt.Println(ParsePR([]byte(`{"eventKey":"pr:opened","date":"2023-05-29T09:19:13-0400","actor":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"pullRequest":{"id":269,"version":0,"title":"PMP-20074 Dataseeding request for azfdemo org - test notification","state":"OPEN","open":true,"closed":false,"createdDate":1685366353203,"updatedDate":1685366353203,"fromRef":{"id":"refs/heads/feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","displayId":"feature/PMP-20074-create-scripts-to-initialize-and-update-the-demo-azfdemo-2023-05-17-na","latestCommit":"9e4da367c66845f212567e839743694f6e8439f6","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"toRef":{"id":"refs/heads/develop","displayId":"develop","latestCommit":"b496d7bc5618fb20faeed7f9f4559fdde4d42adb","type":"BRANCH","repository":{"slug":"pas-data-generation","id":1006,"name":"pas-data-generation","hierarchyId":"bbd53dac7de24a3aee69","scmId":"git","state":"AVAILABLE","statusMessage":"Available","forkable":true,"project":{"key":"PMP","id":212,"name":"CSW LDI (Lenovo Device Intelligence)","description":"(Formerly PMP Predictive Maintenance)","public":false,"type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP"}]}},"public":false,"links":{"clone":[{"href":"https://bitbucket.tc.lenovo.com/scm/pmp/pas-data-generation.git","name":"http"},{"href":"ssh://git@bitbucket.tc.lenovo.com/pmp/pas-data-generation.git","name":"ssh"}],"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/browse"}]}}},"locked":false,"author":{"user":{"name":"kbabin","emailAddress":"kbabin@lenovo.com","id":8902,"displayName":"Konstantin Babin","active":true,"slug":"kbabin","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/kbabin"}]}},"role":"AUTHOR","approved":false,"status":"UNAPPROVED"},"reviewers":[{"user":{"name":"svaladez","emailAddress":"svaladez@lenovo.com","id":1454,"displayName":"Sergio Valadez Cruz","active":true,"slug":"svaladez","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/svaladez"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"},{"user":{"name":"rthoomu1","emailAddress":"rthoomu1@lenovo.com","id":8851,"displayName":"Ramesh Babu Thoomu1","active":true,"slug":"rthoomu1","type":"NORMAL","links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/users/rthoomu1"}]}},"role":"REVIEWER","approved":false,"status":"UNAPPROVED"}],"participants":[],"links":{"self":[{"href":"https://bitbucket.tc.lenovo.com/projects/PMP/repos/pas-data-generation/pull-requests/269"}]}}}`)))
}
