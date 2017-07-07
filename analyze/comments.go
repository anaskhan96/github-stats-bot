package analyze

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/deckarep/golang-set"
	"github.com/jzelinskie/geddit"
)

func AnalyzeComments(comments []*geddit.Comment) error {
	r, err := regexp.Compile(`github\.com\/\w+\/\w+`)
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
			go postReply(comment, links)
		}
	}
	return nil
}

func postReply(comment *geddit.Comment, links []string) {
	for _, link := range links {
		var data map[string]interface{}
		if err := getStats(link, &data); err != nil {
			return
		}
		if data["message"] == "Not Found" {
			return
		}
	}
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
