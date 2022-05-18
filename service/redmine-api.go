package service

import (
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type RedmineConfiguration struct {
	Source      string
	Project     string
	DryRun      bool
	APIEndpoint string
	APIKey      string
}

type Project struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Tracker struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Status struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsClosed bool   `json:"is_closed,omitempty"`
}

type Priority struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Author struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AssignedTo struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type FixedVersion struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CustomField struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Issue struct {
	ID                  int           `json:"id,omitempty"`
	Project             Project       `json:"project,omitempty"`
	Tracker             Tracker       `json:"tracker,omitempty"`
	Status              Status        `json:"status,omitempty"`
	Priority            Priority      `json:"priority,omitempty"`
	Author              Author        `json:"author,omitempty"`
	AssignedTo          AssignedTo    `json:"assigned_to,omitempty"`
	FixedVersion        FixedVersion  `json:"fixed_version,omitempty"`
	Subject             string        `json:"subject,omitempty"`
	Description         string        `json:"description,omitempty"`
	StartDate           string        `json:"start_date,omitempty"`
	DueDate             interface{}   `json:"due_date,omitempty"`
	DoneRatio           int           `json:"done_ratio,omitempty"`
	IsPrivate           bool          `json:"is_private,omitempty"`
	EstimatedHours      interface{}   `json:"estimated_hours,omitempty"`
	TotalEstimatedHours interface{}   `json:"total_estimated_hours,omitempty"`
	SpentHours          float64       `json:"spent_hours,omitempty"`
	TotalSpentHours     float64       `json:"total_spent_hours,omitempty"`
	CustomFields        []CustomField `json:"custom_fields,omitempty"`
	CreatedOn           string        `json:"created_on,omitempty"`
	UpdatedOn           string        `json:"updated_on,omitempty"`
	ClosedOn            interface{}   `json:"closed_on,omitempty"`
}

type RedmineIssues struct {
	Issues     []Issue `json:"issues,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
	Offset     int     `json:"offset,omitempty"`
	Limit      int     `json:"limit,omitempty"`
}

type NewIssue struct {
	ProjectId    int           `json:"project_id,omitempty"`
	PriorityId   int           `json:"priority_id,omitempty"`
	StatusId     int           `json:"status_id,omitempty"`
	Subject      string        `json:"subject,omitempty"`
	Description  string        `json:"description,omitempty"`
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type RedmineIssue struct {
	Issue NewIssue `json:"issue,omitempty"`
}

func CreateIssue(conf *RedmineConfiguration, vulnerability services.Vulnerability) {
	var a = &RedmineIssue{
		Issue: NewIssue{
			ProjectId:   1,
			PriorityId:  1,
			StatusId:    1,
			Subject:     "X-Ray vulnerability " + vulnerability.IssueId,
			Description: vulnerability.Summary,
			CustomFields: []CustomField{
				{
					ID:    1,
					Value: vulnerability.IssueId,
				},
				{
					ID:    2,
					Value: conf.Project,
				},
			},
		},
	}

	issueJson, err := json.Marshal(a)
	if err != nil {
		log.Error(err)
		return
	}

	if conf.DryRun {
		log.Info("DryRun: create issue: ", string(issueJson))
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", conf.APIEndpoint+"/issues.json", strings.NewReader(string(issueJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Redmine-API-Key", conf.APIKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return
	}
	url, _ := resp.Location()
	log.Info("Issue created: ", url)
}

func CloseIssue(conf *RedmineConfiguration, issue Issue) {
	log.Info("Closing redmine issue ", issue.ID)

	var a = &RedmineIssue{
		Issue: NewIssue{
			StatusId: 2,
		},
	}

	issueJson, err := json.Marshal(a)
	if err != nil {
		log.Error(err)
		return
	}

	if conf.DryRun {
		log.Info("DryRun: close issue: ", string(issueJson))
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", conf.APIEndpoint+"/issues/"+strconv.Itoa(issue.ID)+".json", strings.NewReader(string(issueJson)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Redmine-API-Key", conf.APIKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
	}

	if resp.StatusCode == 204 {
		log.Info("Closed ", issue.ID)
	}

}

func GetIssues(conf *RedmineConfiguration) RedmineIssues {
	log.Info("Getting redmine issues for project: ", conf.Project)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", conf.APIEndpoint+"/issues.json?cf_2="+conf.Project+"&limit=100", nil)
	req.Header.Set("X-Redmine-API-Key", conf.APIKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	responseBody := string(body)

	var redmineIssues RedmineIssues
	json.Unmarshal([]byte(responseBody), &redmineIssues)

	log.Info("Open issues found: ", len(redmineIssues.Issues))

	return redmineIssues
}

func MergeIssues(conf *RedmineConfiguration, results []services.ScanResponse, issues RedmineIssues) {

	redmine := map[string]Issue{}
	xray := map[string]services.Vulnerability{}

	//Map redmine issues
	for _, issue := range issues.Issues {
		var issueID string
		for _, customField := range issue.CustomFields {
			if customField.ID == 1 {
				issueID = customField.Value
			}
		}
		redmine[issueID] = issue
	}

	//Map xray vulnerabilities
	for _, result := range results {
		for _, vulnerability := range result.Vulnerabilities {
			xray[vulnerability.IssueId] = vulnerability
			_, hasKey := redmine[vulnerability.IssueId]
			if !hasKey {
				log.Info("Creating a new redmine issue for x-ray vulnerability: ", vulnerability.IssueId)
				CreateIssue(conf, vulnerability)
				redmine[vulnerability.IssueId] = Issue{}
			}
		}
	}

	//delete fixed issues
	for key, value := range redmine {
		_, hasKey := xray[key]
		if !hasKey {
			log.Info("XRay Vulnerability fixed ", value.ID)
			CloseIssue(conf, value)
		}
	}

}
