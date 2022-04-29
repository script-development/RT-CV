package matcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParsedLine struct {
	Name        string
	ParrentName string
	Group1      string
	Group2      string
	Group3      string
	Group4      string
	Group5      string
	Group6      string
}

type B struct {
	Is       TitleKind
	Children map[string]*B
}

func (b *B) toMatcherBranch(name string, parents []primitive.ObjectID, onParsedBranch func(child Branch)) Branch {
	m := db.NewM()

	branches := []primitive.ObjectID{}
	parsedBranches := []Branch{}

	childParents := make([]primitive.ObjectID, len(parents)+1)
	copy(childParents, append(parents, m.ID))

	for title, b := range b.Children {
		child := b.toMatcherBranch(title, childParents, onParsedBranch)
		branches = append(branches, child.ID)
		parsedBranches = append(parsedBranches, child)
		onParsedBranch(child)
	}

	return Branch{
		M:              m,
		Titles:         []string{name},
		Branches:       branches,
		ParsedBranches: parsedBranches,
		Parents:        parents,
	}
}

type PathPart struct {
	Is   TitleKind
	Name string
}

func ImportUWVJobs(dbConn db.Connection) {
	csvBytes, err := ioutil.ReadFile("/Users/mark/Downloads/uwv-beroepen.csv")
	if err != nil {
		log.Fatal(err.Error())
	}

	lines := strings.Split(string(csvBytes), "\r\n")
	jobs := []ParsedLine{}
	for _, line := range lines[2:] {
		lineSplit := parseCSVLine(line)
		if len(lineSplit) < 9 {
			continue
		}

		job := ParsedLine{
			Name:        strings.TrimSpace(lineSplit[1]),
			ParrentName: strings.TrimSpace(lineSplit[3]),
			Group1:      strings.TrimSpace(lineSplit[4]),
			Group2:      strings.TrimSpace(lineSplit[5]),
			Group3:      strings.TrimSpace(lineSplit[6]),
			Group4:      strings.TrimSpace(lineSplit[7]),
			Group5:      strings.TrimSpace(lineSplit[8]),
			Group6:      strings.TrimSpace(lineSplit[9]),
		}
		if job.Name == "" || job.Group1 == "" {
			continue
		}
		jobs = append(jobs, job)
	}

	tree := &B{
		Children: map[string]*B{},
	}

	for _, job := range jobs {
		path := []PathPart{{Sector, job.Group1}}
		if len(job.Group2) != 0 {
			path = append(path, PathPart{Sector, job.Group2})
			if len(job.Group3) != 0 {
				path = append(path, PathPart{Sector, job.Group3})
				if len(job.Group4) != 0 {
					path = append(path, PathPart{Sector, job.Group4})
					if len(job.Group5) != 0 {
						path = append(path, PathPart{Sector, job.Group5})
						if len(job.Group6) != 0 {
							path = append(path, PathPart{Sector, job.Group6})
						}
					}
				}
			}
		}
		if len(job.ParrentName) != 0 {
			path = append(path, PathPart{Job, job.ParrentName})
		}
		path = append(path, PathPart{Job, job.Name})

		currentTreePath := tree
		for _, part := range path {
			newBranch, ok := currentTreePath.Children[part.Name]
			if ok {
				currentTreePath = newBranch
			} else {
				currentTreePath.Children[part.Name] = &B{
					Is:       part.Is,
					Children: map[string]*B{},
				}
				currentTreePath = currentTreePath.Children[part.Name]
			}
		}
	}

	branches := []Branch{}
	addBranch := func(branch Branch) {
		branches = append(branches, Branch{
			M:              db.M{ID: branch.ID},
			Titles:         branch.Titles,
			TitleKind:      branch.TitleKind,
			Branches:       branch.Branches,
			ParsedBranches: branch.ParsedBranches,
			Parents:        branch.Parents,
		})
	}
	for title, child := range tree.Children {
		addBranch(child.toMatcherBranch(title, []primitive.ObjectID{}, addBranch))
	}

	insertionData := []db.Entry{}
	for i := range branches {
		insertionData = append(insertionData, &branches[i])
	}

	fmt.Println(dbConn.Insert(insertionData...))
}

func parseCSVLine(line string) []string {
	result := []string{}
	quoteSplits := strings.Split(line, "\"")
	for idx, quote := range quoteSplits {
		canSplit := idx%2 == 0 // All even pars of the split can be formatted, the other once are escaped using quotes
		if canSplit {
			splits := strings.Split(quote, ",")
			removePrefixComma := idx > 0
			removeSuffixComma := idx < len(quoteSplits)-1
			if removePrefixComma {
				splits = splits[1:]
			}
			if removeSuffixComma {
				splits = splits[:len(splits)-1]
			}
			result = append(result, splits...)
		} else {
			result = append(result, quote)
		}
	}
	return result
}
