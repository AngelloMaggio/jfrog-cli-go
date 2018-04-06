package utils

import (
	"bufio"
	"errors"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"os"
	"path/filepath"
	"strings"
)

type gitManager struct {
	path     string
	err      error
	revision string
	url      string
}

func NewGitManager(path string) *gitManager {
	dotGitPath := filepath.Join(path, ".git")
	return &gitManager{path: dotGitPath}
}

func (m *gitManager) ReadGitConfig() error {
	if m.path == "" {
		return errorutils.CheckError(errors.New(".git path must be defined."))
	}
	m.readRevision()
	m.readUrl()
	return m.err
}

func (m *gitManager) GetUrl() string {
	return m.url
}

func (m *gitManager) GetRevision() string {
	return m.revision
}

func (m *gitManager) readUrl() {
	if m.err != nil {
		return
	}
	dotGitPath := filepath.Join(m.path, "config")
	file, err := os.Open(dotGitPath)
	if errorutils.CheckError(err) != nil {
		m.err = err
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var IsNextLineUrl bool
	var originUrl string
	for scanner.Scan() {
		if IsNextLineUrl {
			text := scanner.Text()
			strings.HasPrefix(text, "url")
			originUrl = strings.TrimSpace(strings.SplitAfter(text, "=")[1])
			break
		}
		if scanner.Text() == "[remote \"origin\"]" {
			IsNextLineUrl = true
		}
	}
	if err := scanner.Err(); err != nil {
		errorutils.CheckError(err)
		m.err = err
		return
	}
	m.url = originUrl
}

func (m *gitManager) getRevisionOrBranchPath() (revision, refUrl string, err error) {
	dotGitPath := filepath.Join(m.path, "HEAD")
	file, e := os.Open(dotGitPath)
	if errorutils.CheckError(e) != nil {
		err = e
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "ref") {
			refUrl = strings.TrimSpace(strings.SplitAfter(text, ":")[1])
			break
		}
		revision = text
	}
	if err = scanner.Err(); err != nil {
		errorutils.CheckError(err)
	}
	return
}

func (m *gitManager) readRevision() {
	if m.err != nil {
		return
	}
	// This function will either return the revision or the branch ref:
	revision, ref, err := m.getRevisionOrBranchPath()
	if err != nil {
		m.err = err
		return
	}
	// If the revision was returned, then we're done:
	if revision != "" {
		m.revision = revision
		return
	}

	// Since the revision was not returned, then we'll fetch it, by using the ref:
	dotGitPath := filepath.Join(m.path, ref)
	file, err := os.Open(dotGitPath)
	if err != nil {
		m.err = err
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		revision = strings.TrimSpace(text)
		break
	}
	if err := scanner.Err(); err != nil {
		errorutils.CheckError(err)
		m.err = err
		return
	}

	m.revision = revision
	return
}
