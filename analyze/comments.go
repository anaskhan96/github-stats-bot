package analyze

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/deckarep/golang-set"
	"github.com/jzelinskie/geddit"
)

func AnalyzeComments(s *geddit.OAuthSession, comments []*geddit.Comment) error {
	r, err := regexp.Compile(`github\.com\/[\w\-]+\/[\w\-]+`)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		if comment.Author == "github-stats-bot" {
			continue
		}
		if r.MatchString(comment.Body) {
			linkSet := mapset.NewSet()
			for _, link := range r.FindAllString(comment.Body, -1) {
				linkSet.Add(link)
			}
			var links []string
			for _, link := range linkSet.ToSlice() {
				links = append(links, link.(string))
			}
			if err = postReply(s, comment, links); err != nil {
				return err
			}
		}
	}
	return nil
}

func postReply(s *geddit.OAuthSession, comment *geddit.Comment, links []string) error {
	footer := "***\n^(This is Earth radio, and now here's human music ♫)\n\n^[Source](https://github.com/anaskhan96/github-stats-bot) ^| ^[PMme](https://np.reddit.com/message/compose?to=github-stats-bot)"
	var reply string
	for _, link := range links {
		var data map[string]interface{}
		if err := getStats(link, &data); err != nil {
			return err
		}
		if data["message"] == "Not Found" {
			return errors.New("Wrong GitHub API endpoint")
		}
		description := data["description"]
		stargazers := data["stargazers_count"]
		forks := data["forks_count"]
		issuesURL := "https://" + link + "/issues"
		pullsURL := "https://" + link + "/pulls"
		reply += fmt.Sprintf("\n[%s](https://%s)\n\n> *Description*: %v\n\n> *Stars*: %v\n\n> *Forks*: %v\n\n> [Issues](%s) | [Pull Requests](%s)\n\n",
			link[11:], link, description, stargazers, forks, issuesURL, pullsURL)
	}
	reply += footer
	if _, err := s.Reply(comment, reply); err != nil {
		return err
	}
	return nil
}

func getStats(link string, data *map[string]interface{}) error {
	res, err := http.Get("https://api.github.com/repos/" + link[11:])
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}
	return nil
}
