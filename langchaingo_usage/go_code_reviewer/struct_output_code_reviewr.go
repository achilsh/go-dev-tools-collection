package gocodereviewer

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

type Issue struct {
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

type ReviewResult struct {
	Filename string  `json:"filename"`
	Issues   []Issue `json:"issues"`
	Score    int     `json:"score"` // 0-100
}

type CodeReviewerStructOutput struct {
	llm llms.Model
}

func NewCodeReviewerStructOutput() (*CodeReviewerStructOutput, error) {
	llm, err := openai.New(openai.WithModel("gpt-4o-mini-2024-07-18"))
	if err != nil {
		return nil, err
	}

	return &CodeReviewerStructOutput{
		llm: llm,
	}, nil
}

func (cr *CodeReviewerStructOutput) ReviewFile(filename string) (*ReviewResult, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// Parse for line numbers
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing Go file: %w", err)
	}

	template := prompts.NewPromptTemplate(`
Analyze this Go code and return a JSON response with this exact structure:

{
  "filename": "{{.filename}}",
  "issues": [
    {
      "severity": "critical|warning|suggestion",
      "type": "bug|performance|style|security",
      "line": 42,
      "description": "Detailed issue description",
      "suggestion": "How to fix this issue"
    }
  ],
  "score": 85
}

Code to analyze:
'''go
{{.code}}
'''

Focus on real issues. Score: 100 = perfect, 0 = many serious issues.`,
		[]string{"code", "filename"})

	prompt, err := template.Format(map[string]any{
		"code":     string(content),
		"filename": filename,
	})
	if err != nil {
		return nil, fmt.Errorf("formatting prompt: %w", err)
	}

	ctx := context.Background()
	response, err := cr.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}, llms.WithJSONMode())
	if err != nil {
		return nil, fmt.Errorf("generating review: %w", err)
	}

	var result ReviewResult
	if err := json.Unmarshal([]byte(response.Choices[0].Content), &result); err != nil {
		return nil, fmt.Errorf("parsing JSON response: %w", err)
	}

	return &result, nil
}

func ReviewDirectoryStructOutput(reviewer *CodeReviewerStructOutput, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			ret, err := reviewer.ReviewFile(path)
			if err != nil {
				return err
			}
			retStr, _ := json.MarshalIndent(ret, "", "\t")
			fmt.Printf("%+v\n", string(retStr))
		}
		return nil
	})
}

// level: 0 current modify file, 1: preview-1 git commit files
func ReviewGitChangesStructOutput(reviewer *CodeReviewerStructOutput, level int) error {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("HEAD~%d", level))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("get git change files err: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range files {
		if strings.HasSuffix(file, ".go") && file != "" {
			ret, err := reviewer.ReviewFile(file)
			if err != nil {
				fmt.Printf("review file: %v, err: %v", file, err)
			} else {
				retStr, _ := json.MarshalIndent(ret, "", "\t")
				fmt.Printf("%+v\n", string(retStr))
			}

		}
	}
	return nil
}
func RunReviewercodeStructOutput() {
	var (
		file = flag.String("file", "", "Go file to review")
		dir  = flag.String("dir", "", "Directory to review (all .go files)")
		git  = flag.Int("git", -1, "Review files changed in git commit level for working directory")
	)
	flag.Parse()

	reviewer, err := NewCodeReviewerStructOutput()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *file != "":
		if ret, err := reviewer.ReviewFile(*file); err != nil {
			log.Fatal(err)
		} else {
			retStr, _ := json.MarshalIndent(ret, "", "\t")
			fmt.Printf("%+v\n", string(retStr))
		}
	case *dir != "":
		if err := ReviewDirectoryStructOutput(reviewer, *dir); err != nil {
			log.Fatal(err)
		}
	case *git >= 0:
		if err := ReviewGitChangesStructOutput(reviewer, *git); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Usage:")
		fmt.Println("  code-reviewer -file=main.go")
		fmt.Println("  code-reviewer -dir=./pkg")
		fmt.Println("  code-reviewer -git=0")
		os.Exit(1)
	}
}
