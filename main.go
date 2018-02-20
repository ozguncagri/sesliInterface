package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sesliInterface",
		Short: "SesliSözlük Dictionary Interface for Alfred",
		Long:  "SesliSözlük Dictionary Interface for Alfred. It gets all arguments for search input and returns back all results as alfred compatible JSON format.",
		Run: func(cmd *cobra.Command, args []string) {
			searchURL := fmt.Sprintf("http://m.seslisozluk.net/index6.php?word=%v&dN=iPhone", strings.Join(args, "%20"))
			results := extractResults(searchURL)
			fmt.Println(alfredResultsToJSON(results))
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func extractResults(pageURL string) []AlfredResult {
	doc, err := goquery.NewDocument(pageURL)

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	allResults := []AlfredResult{}

	doc.Find("body .resultset .dict_result ol").Each(func(i int, s *goquery.Selection) {
		s.Find("li").Each(func(i int, li *goquery.Selection) {
			result := li.Text()
			result = strings.Replace(result, "  ", " : ", 0)
			result = strings.Replace(result, "  ", " ", -1)
			result = strings.Trim(result, " ")

			tmp := AlfredResult{
				Title:        result,
				QuickLookURL: pageURL,
				Text: AlfredResultItemText{
					Copy:      result,
					LargeType: result,
				},
			}

			allResults = append(allResults, tmp)
		})
	})

	if len(allResults) == 0 {
		doc.Find("body #didumean .resultset .dict_result ul").Each(func(i int, s *goquery.Selection) {
			s.Find("li").Each(func(i int, li *goquery.Selection) {
				result := li.Text()
				result = strings.Replace(result, "  ", " : ", 0)
				result = strings.Replace(result, "  ", " ", -1)
				result = strings.Trim(result, " ")

				tmp := AlfredResult{
					Title:        result,
					Subtitle:     "Did you mean " + result + "?",
					AutoComplete: result,
					QuickLookURL: pageURL,
					Text: AlfredResultItemText{
						Copy:      result,
						LargeType: result,
					},
				}

				allResults = append(allResults, tmp)
			})
		})
	}

	return allResults
}

func alfredResultsToJSON(results []AlfredResult) string {
	data := make(map[string][]AlfredResult)

	data["items"] = results

	out, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	return string(out)
}

//AlfredResult is general result item for searches
type AlfredResult struct {
	Title        string               `json:"title"`
	Subtitle     string               `json:"subtitle,omitempty"`
	AutoComplete string               `json:"autocomplete,omitempty"`
	QuickLookURL string               `json:"quicklookurl"`
	Text         AlfredResultItemText `json:"text"`
}

//AlfredResultItemText contains datas for alfred's copy and large type features for results
type AlfredResultItemText struct {
	Copy      string `json:"copy"`
	LargeType string `json:"largetype"`
}
