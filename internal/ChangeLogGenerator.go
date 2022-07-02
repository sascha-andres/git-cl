package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"text/template"
)

var (
	subjectRegex = regexp.MustCompile("(?P<type>feat|fix|doc|chore|refactor|test|style|perf|other)(\\((?P<issue>[^)]+)\\))?: (?P<message>.*)")
	l            = log.New(os.Stdout, "[ChangeLogGenerator] ", log.LstdFlags)
)

type (
	// ChangeLogGenerator allows creates a changelog for a list of commit subjects
	ChangeLogGenerator struct {
		/* input */

		// reader is where to get the subjects lines from
		reader io.Reader

		/* configuration */

		// groups maps a type name to a group name
		groups map[string]string
		// Header text that will be added to the beginning of the changelog
		Header string
		// Footer text that will be added to the end of the changelog
		Footer string
		// BodyTemplate, represents a single release in the changelog
		BodyTemplate string

		/* runtime stuff */

		// lines passed to the tool
		lines []string
		// commits groups parsed lines
		commits map[string][]*Commit
		// bodyTemplate is the instantiated template
		bodyTemplate *template.Template
		// changelogTemplate is the overall document
		changelogTemplate *template.Template
	}

	// ChangeLogData is passed to changelogTemplate as data
	ChangeLogData struct {
		// Header text that will be added to the beginning of the changelog
		Header string
		// Footer text that will be added to the end of the changelog
		Footer string
		// Body is the list of commits processed with the body template
		Body string
	}

	// Commit represents one commit subject line
	Commit struct {
		// Type of commit (eg feat)
		Type string
		// Scope of commit or issue
		Scope string
		// Subject og commit
		Subject string
	}

	// ChangeLogGeneratorOption can be used to change options
	ChangeLogGeneratorOption func(generator *ChangeLogGenerator)
)

// NewChangeLogGenerator returns a changelog generator
func NewChangeLogGenerator(reader io.Reader, options ...ChangeLogGeneratorOption) (*ChangeLogGenerator, error) {
	clg := &ChangeLogGenerator{
		reader: reader,
		lines:  make([]string, 0),
		groups: map[string]string{
			"feat":     "Feature",
			"fix":      "Bugfix",
			"doc":      "Other",
			"chore":    "Other",
			"refactor": "Optimization",
			"test":     "Optimization",
			"style":    "Usability",
			"perf":     "Optimization",
			"other":    "Other",
		},
		commits: make(map[string][]*Commit),
		BodyTemplate: `{{ range $index, $element := . }}## {{ $index }}{{ range $val := $element }}
- {{ $val.Subject }}{{ if $val.Scope }} ({{ $val.Scope }}){{ end }}{{ end }}

{{ end }}`,
		Header: "# Changelog",
		Footer: "generated by git-cl",
	}
	var err error
	clg.changelogTemplate, err = template.New("clt").Parse(`{{ .Header }}

{{ .Body }}{{ .Footer }}`)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		if nil == opt { // quirk when using functional options, it is possible to pass nil
			continue
		}
		opt(clg)
	}
	return clg, clg.readSubjectLines()
}

// readSubjectLines gets all lines from reader, called from NewChangeLogGenerator
func (clg *ChangeLogGenerator) readSubjectLines() error {
	rd := bufio.NewReader(clg.reader)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			l.Printf("read file line error: %v", err)
			return err
		}
		clg.lines = append(clg.lines, line)
	}

	return nil
}

// OverrideGroupForType allows changing the group for a type
func (clg *ChangeLogGenerator) OverrideGroupForType(typeName, group string) {
	if _, ok := clg.groups[typeName]; ok {
		clg.groups[typeName] = group
	}
}

// Build runs through the list of provided commit messages and creates the MarkDown output
func (clg *ChangeLogGenerator) Build() (result string, err error) {
	result = ""

	err = clg.loadChangeLogEntryTemplate()
	clg.parseAndGroupCommits()
	wr := bytes.Buffer{}
	err = clg.bodyTemplate.Execute(&wr, clg.commits)

	cld := &ChangeLogData{
		Header: clg.Header,
		Footer: clg.Footer,
		Body:   string(wr.Bytes()),
	}

	err = clg.changelogTemplate.Execute(os.Stdout, cld)
	return
}

// parseAndGroupCommits
func (clg *ChangeLogGenerator) parseAndGroupCommits() {
	for _, line := range clg.lines {
		if !subjectRegex.MatchString(line) {
			continue
		}
		matches := subjectRegex.FindStringSubmatch(line)
		groupNames := subjectRegex.SubexpNames()
		commit := &Commit{}
		for i, match := range matches {
			if groupNames[i] == "type" {
				commit.Type = match
			}
			if groupNames[i] == "message" {
				commit.Subject = match
			}
			if groupNames[i] == "issue" {
				commit.Scope = match
			}
		}
		gn := clg.getGroup(commit)
		if gn == "" {
			l.Printf("omitting commit (%s), does not match", commit)
			continue
		}
		if _, ok := clg.commits[gn]; !ok {
			clg.commits[gn] = make([]*Commit, 0)
		}
		clg.commits[gn] = append(clg.commits[gn], commit)
	}
}

// getGroup looks up group for commit
func (clg *ChangeLogGenerator) getGroup(c *Commit) string {
	if gn, ok := clg.groups[c.Type]; ok {
		return gn
	}
	return ""
}

// loadChangeLogEntryTemplate loads the textually provided template in a template instance
func (clg *ChangeLogGenerator) loadChangeLogEntryTemplate() (err error) {
	clg.bodyTemplate, err = template.New("cle").Parse(clg.BodyTemplate)
	if err != nil {
		l.Printf("error constructing template: %s", err)
		return
	}
	return
}

// String shows a visual representation of a commit
func (c Commit) String() string {
	return fmt.Sprintf("type := %s scope := %s subject := %s", c.Type, c.Scope, c.Subject)
}
