package rmp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	googleQuery "github.com/google/go-querystring/query"
)

// a search query for rmp
type SearchQuery struct {
	Query     string `json:"query"`
	Variables struct {
		query    `json:"query"`
		SchoolID string `json:"schoolID"`
	} `json:"variables"`
}

// query component within `SearchQuery`
type query struct {
	DepartmentID interface{} `json:"departmentID"`
	Fallback     bool        `json:"fallback"`
	SchoolID     string      `json:"schoolID"`
	Text         string      `json:"text"`
}

// a single prof from rmp
type ProfNode struct {
	Cursor string `json:"cursor"`
	Node   struct {
		Typename      string  `json:"__typename"`
		AvgDifficulty float64 `json:"avgDifficulty"`
		AvgRating     float64 `json:"avgRating"`
		Department    string  `json:"department"`
		FirstName     string  `json:"firstName"`
		ID            string  `json:"id"`
		IsSaved       bool    `json:"isSaved"`
		LastName      string  `json:"lastName"`
		LegacyID      int64   `json:"legacyId"`
		NumRatings    int64   `json:"numRatings"`
		School        struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"school"`
		WouldTakeAgainPercent float64 `json:"wouldTakeAgainPercent"`
	} `json:"node"`
}

// profDescription generate a description the professor
func (p *ProfNode) profDescription() string {
	return fmt.Sprintf(`- **Average rating**: %f
- **Average difficulty**: %f
- **Would take again**: %f%%`, p.Node.AvgRating, p.Node.AvgDifficulty, p.Node.WouldTakeAgainPercent)
}

// fullName return a profs first and last name
func (p *ProfNode) fullName() string {
	return fmt.Sprintf("%s %s", p.Node.FirstName, p.Node.LastName)
}

// rmpURL generate a url for the current professor to their RMP page
func (p *ProfNode) rmpURL() string {
	type Options struct {
		Tid int64 `url:"tid"`
	}
	opt := Options{Tid: p.Node.LegacyID}
	v, _ := googleQuery.Values(opt)
	return fmt.Sprintf("https://www.ratemyprofessors.com/professor?%s", v.Encode())
}

// api response when searching for a prof
type searchResponse struct {
	Data struct {
		Search struct {
			Teachers struct {
				DidFallback bool       `json:"didFallback"`
				Edges       []ProfNode `json:"edges"`
				Filters     []struct {
					Field   string `json:"field"`
					Options []struct {
						ID    string `json:"id"`
						Value string `json:"value"`
					} `json:"options"`
				} `json:"filters"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
				ResultCount int64 `json:"resultCount"`
			} `json:"teachers"`
		} `json:"search"`
	} `json:"data"`
	Errors []struct {
		Locations []struct {
			Column int64 `json:"column"`
			Line   int64 `json:"line"`
		} `json:"locations"`
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

// generateQuery generate a search query for rmp to search for a prof by name
// profName: the name of the prof to search for
// returns a `SearchQuery` struct
func generateQuery(profName string) SearchQuery {
	return SearchQuery{
		Query: "query TeacherSearchResultsPageQuery(\n  $query: TeacherSearchQuery!\n  $schoolID: ID\n) {\n  search: newSearch {\n    ...TeacherSearchPagination_search_1ZLmLD\n  }\n  school: node(id: $schoolID) {\n    __typename\n    ... on School {\n      name\n    }\n    id\n  }\n}\n\nfragment TeacherSearchPagination_search_1ZLmLD on newSearch {\n  teachers(query: $query, first: 8, after: \"\") {\n    didFallback\n    edges {\n      cursor\n      node {\n        ...TeacherCard_teacher\n        id\n        __typename\n      }\n    }\n    pageInfo {\n      hasNextPage\n      endCursor\n    }\n    resultCount\n    filters {\n      field\n      options {\n        value\n        id\n      }\n    }\n  }\n}\n\nfragment TeacherCard_teacher on Teacher {\n  id\n  legacyId\n  avgRating\n  numRatings\n  ...CardFeedback_teacher\n  ...CardSchool_teacher\n  ...CardName_teacher\n  ...TeacherBookmark_teacher\n}\n\nfragment CardFeedback_teacher on Teacher {\n  wouldTakeAgainPercent\n  avgDifficulty\n}\n\nfragment CardSchool_teacher on Teacher {\n  department\n  school {\n    name\n    id\n  }\n}\n\nfragment CardName_teacher on Teacher {\n  firstName\n  lastName\n}\n\nfragment TeacherBookmark_teacher on Teacher {\n  id\n  isSaved\n}\n",
		Variables: struct {
			query    `json:"query"`
			SchoolID string `json:"schoolID"`
		}{
			query: query{
				DepartmentID: nil,
				Fallback:     true,
				SchoolID:     "U2Nob29sLTE0OTc=",
				Text:         profName,
			},
			SchoolID: "U2Nob29sLTE0OTc=",
		},
	}
}

// SearchRmpProfByName search for a professor seneca professor by name on RMP, using the graphql api
// name: the name of the professor to search for
// Returns a `searchResponse` struct
func SearchRmpProfByName(name string) searchResponse {
	// generate a query for the prof name
	query := generateQuery(name)
	j, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	// make request with above query
	req, err := http.NewRequest(http.MethodPost, "https://www.ratemyprofessors.com/graphql", bytes.NewBuffer([]byte(j)))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", " Basic dGVzdDp0ZXN0")
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var result searchResponse
	json.Unmarshal(body, &result)
	return result
}

// FilterProfNodes get prof nodes from search result.
// return: list of prof nodes
func FilterProfNodes(profs searchResponse) []ProfNode {
	senecaProfs := []ProfNode{}
	for _, prof := range profs.Data.Search.Teachers.Edges {
		senecaProfs = append(senecaProfs, prof)
	}
	return senecaProfs
}
