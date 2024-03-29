package caramelbotrmp

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly"
	"github.com/shawnyu5/debate_dragon_2.0/commands/rmp"
)

type ProfessorQueryResult struct {
	Professors         []Professor
	SearchResultsTotal int
	Remaining          int
	Type               string
}

type Professor struct {
	TDept            string
	TSid             string
	Institution_name string
	TFname           string
	TMiddleName      string
	TLname           string
	TId              int
	TNumRatings      int
	Rating_class     string
	ContentType      string
	CategoryType     string
	Overall_Rating   string
}

type RatingQueryResult struct {
	Ratings   []Rating
	Remaining int
}

type Rating struct {
	Attendence        string
	ClarityColor      string
	EasyColor         string
	HelpColor         string
	HelpCount         int
	Id                int
	NotHelpCount      int
	OnlineClass       string
	Quality           string
	RClarity          int
	RClass            string
	RComments         string
	RDate             string
	REasy             float64
	REasyString       string
	RErrorMsg         string
	RHelpful          int
	RInterest         string
	ROverall          float64
	ROverallString    string
	RStatus           int
	RTextBookUse      string
	RTimestamp        int
	RWouldTakeAgain   string
	SId               int
	TakenForCredit    string
	Teacher           string
	TeacherGrade      string
	TeacherRatingRags []string
	UnUsefulGrouping  string
	UsefulGrouping    string
}

type RMPResult struct {
	professorName       string
	totalRating         string
	numRatings          string
	courses             []string
	totalRatingByCourse map[string]float64
	numRatingsByCourse  map[string]int
	wouldTakeAgain      string
	levelOfDifficulty   string
	ratingDistribution  map[int]int
	topTags             []string
	rmpURL              string
}

func QueryProfessor(professor string) (RMPResult, error) {
	// var result RMPResult

	// // professorQuery := fmt.Sprint("https://www.ratemyprofessors.com/filter/professor/?&page=1&filter=teacherlastname_sort_s+asc&query=", url.QueryEscape(professor), "&queryoption=TEACHER&queryBy=schoolId&sid=1497")
	// professorQuery := fmt.Sprint("https://www.ratemyprofessors.com/search/professors?q=", url.QueryEscape(professor))
	// fmt.Println("Checking query: ", professorQuery)

	// // Query the API for the professor name given
	// professorResponse, err := http.Get(professorQuery)
	// if err != nil {
	//    fmt.Println("Error contacting RateMyProfessors API server")
	//    fmt.Println(err)
	//    return result, errors.New("error contacting RateMyProfessors API server")
	// }
	// defer professorResponse.Body.Close()

	// var professorBody ProfessorQueryResult

	// // Decode the JSON response into a struct
	// err = json.NewDecoder(professorResponse.Body).Decode(&professorBody)
	// if err != nil {
	//    fmt.Println("Error decoding JSON")
	//    return result, errors.New("error decoding JSON")
	// }

	// // If no professors were found, return an error
	// if professorBody.SearchResultsTotal == 0 {
	//    return result, errors.New("professor not found")
	// }

	// fmt.Println("Found ", professorBody.SearchResultsTotal, " results")
	// fmt.Println("First professor found: ", professorBody.Professors[0].TFname, professorBody.Professors[0].TMiddleName, professorBody.Professors[0].TLname)

	// result.professorName = professorBody.Professors[0].TFname + " " + professorBody.Professors[0].TMiddleName + " " + professorBody.Professors[0].TLname

	var result RMPResult
	searchResult := rmp.SearchRmpProfByName(professor)
	// If multiple profs found, use the first prof
	// if len(prof.Data.Search.Teachers.Edges) != 1 {

	// }
	prof := searchResult.Data.Search.Teachers.Edges[0].Node
	log.Debugf("Prof name is %s", prof.FullName())
	scrapeURL := fmt.Sprint("https://www.ratemyprofessors.com/ShowRatings.jsp?tid=", prof.ID)

	result.rmpURL = scrapeURL

	// Create web scraper
	c := colly.NewCollector(
		colly.AllowedDomains("www.ratemyprofessors.com"),
	)

	// Scrape the overall rating, % who would take again, and level of difficulty
	c.OnHTML("div[class]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "RatingValue__Numerator") {
			result.totalRating = e.Text
		}
		if strings.Contains(e.Attr("class"), "FeedbackItem__FeedbackNumber") {
			if strings.Contains(e.Text, "%") {
				result.wouldTakeAgain = e.Text
			} else {
				result.levelOfDifficulty = e.Text
			}
		}
		if strings.Contains(e.Attr("class"), "TeacherTags__TagsContainer") {
			e.ForEach("span", func(_ int, elem *colly.HTMLElement) {
				result.topTags = append(result.topTags, elem.Text)
			})
		}
	})

	// Scrape the number o fratings
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("href") == "#ratingsList" {
			result.numRatings = strings.Split(e.Text, " ")[0]
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(scrapeURL)

	if (result.topTags == nil) || (result.topTags[0] == "") || len(result.topTags) == 0 {
		result.topTags = []string{"No tags"}
	}

	if result.wouldTakeAgain == "" {
		result.wouldTakeAgain = "N/A"
	}

	if result.levelOfDifficulty == "" {
		result.levelOfDifficulty = "N/A"
	}

	// Pagination variables
	page := 1
	remaining := 1
	ratingQuery := fmt.Sprint("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=", prof.ID, "&filter=&courseCode=&page=", page)

	// Start building result
	result.professorName = prof.FullName()
	result.numRatingsByCourse = make(map[string]int)
	result.totalRatingByCourse = make(map[string]float64)

	// To try to make up for RMP's artificial skew
	result.ratingDistribution = make(map[int]int)

	// Loop until no ratings are left
	for ok := true; ok; ok = remaining > 0 {
		fmt.Println("Checking query: ", ratingQuery)
		ratingResponse, err := http.Get(ratingQuery)

		// Query the API for ratings
		if err != nil {
			// fmt.Println("Error contacting RateMyProfessors API server")
			// fmt.Println(err)
			return result, errors.New("error contacting RateMyProfessors API server")
		}
		defer ratingResponse.Body.Close()

		// Decode the JSON response into a struct
		var ratingBody RatingQueryResult
		err = json.NewDecoder(ratingResponse.Body).Decode(&ratingBody)
		if err != nil {
			// fmt.Println("Error decoding JSON")
			// fmt.Println(err)
			return result, errors.New("error decoding JSON")
		}

		// Check remaining ratings
		remaining = ratingBody.Remaining
		fmt.Println(ratingBody.Remaining, " results remaining")
		page++
		ratingQuery = fmt.Sprint("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=", prof.ID, "&filter=&courseCode=&page=", page)

		if len(ratingBody.Ratings) > 0 {
			// Incorporate ratings into result
			for _, rating := range ratingBody.Ratings {
				// Ignore rating if not helpful AND a 1.0 rating (not sure why, seems fucked up tbh)
				result.ratingDistribution[int(math.Ceil(rating.ROverall))]++
				if result.numRatingsByCourse[rating.RClass] == 0 {
					result.courses = append(result.courses, rating.RClass)
				}
				result.totalRatingByCourse[rating.RClass] += float64(rating.RClarity) + float64(rating.RHelpful)
				result.numRatingsByCourse[rating.RClass] += 2
			}
		} else {
			return result, errors.New("no ratings found")
		}
	}

	for _, course := range result.courses {
		result.totalRatingByCourse[course] = result.totalRatingByCourse[course] / float64(result.numRatingsByCourse[course])
		result.numRatingsByCourse[course] /= 2
	}

	return result, nil
}

func CompareOverallRating(rmp1, rmp2 RMPResult) string {
	if rmp1.totalRating > rmp2.totalRating {
		return rmp1.professorName + " has an overall rating of " + rmp1.totalRating + "/5" + " with " + rmp1.numRatings + " ratings"
	} else if rmp1.totalRating < rmp2.totalRating {
		return rmp2.professorName + " has an overall rating of " + rmp2.totalRating + "/5" + " with " + rmp2.numRatings + " ratings"
	} else {
		return "Both professors have an overall rating of " + rmp1.totalRating + "/5"
	}
}

func CompareWouldTakeAgain(rmp1, rmp2 RMPResult) string {
	if (rmp1.wouldTakeAgain == "N/A") && (rmp2.wouldTakeAgain != "N/A") {
		return rmp2.professorName + " has " + rmp2.wouldTakeAgain + " who would take again"
	} else if (rmp1.wouldTakeAgain != "N/A") && (rmp2.wouldTakeAgain == "N/A") {
		return rmp1.professorName + " has " + rmp1.wouldTakeAgain + " who would take again"
	}
	if rmp1.wouldTakeAgain > rmp2.wouldTakeAgain {
		return rmp1.professorName + " has " + rmp1.wouldTakeAgain + " who would take again"
	} else if rmp1.wouldTakeAgain < rmp2.wouldTakeAgain {
		return rmp2.professorName + " has " + rmp2.wouldTakeAgain + " who would take again"
	} else {
		return "Both professors have " + rmp1.wouldTakeAgain + " who would take again"
	}
}

func CompareDifficulty(rmp1, rmp2 RMPResult) string {
	if (rmp1.levelOfDifficulty == "N/A") && (rmp2.levelOfDifficulty != "N/A") {
		return rmp2.professorName + " has " + rmp2.levelOfDifficulty + " who would take again"
	} else if (rmp1.levelOfDifficulty != "N/A") && (rmp2.levelOfDifficulty == "N/A") {
		return rmp1.professorName + " has " + rmp1.levelOfDifficulty + " who would take again"
	}
	if rmp1.levelOfDifficulty < rmp2.levelOfDifficulty {
		return rmp1.professorName + " has a level of difficulty of " + rmp1.levelOfDifficulty
	} else if rmp1.levelOfDifficulty > rmp2.levelOfDifficulty {
		return rmp2.professorName + " has a level of difficulty of " + rmp2.levelOfDifficulty
	} else {
		return "Both professors have a level of difficulty of " + rmp1.levelOfDifficulty
	}
}

func CompareBestByCourse(rmp1, rmp2 RMPResult) string {
	bestByCourse := ""
	for _, course := range rmp1.courses {
		// If the course is in both professors' lists, compare ratings
		if rmp2.numRatingsByCourse[course] > 0 {
			bestByCourse += "**" + course + "**" + ": "
			if rmp1.totalRatingByCourse[course] > rmp2.totalRatingByCourse[course] {
				bestByCourse += fmt.Sprintf("%s%s%.2f%s%s%d%s", rmp1.professorName, " has a rating of ", rmp1.totalRatingByCourse[course], "/5\n", "with ", rmp1.numRatingsByCourse[course], " ratings\n")
			} else if rmp1.totalRatingByCourse[course] < rmp2.totalRatingByCourse[course] {
				bestByCourse += fmt.Sprintf("%s%s%.2f%s%s%d%s", rmp2.professorName, " has a rating of ", rmp2.totalRatingByCourse[course], "/5\n", "with ", rmp2.numRatingsByCourse[course], " ratings\n")
			} else {
				bestByCourse += fmt.Sprintf("%s%.2f%s", "Both professors have a rating of ", rmp1.totalRatingByCourse[course], "/5\n")
			}
		}
	}

	if bestByCourse == "" {
		bestByCourse = "No courses found in common"
	}

	return bestByCourse
}

func stringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
