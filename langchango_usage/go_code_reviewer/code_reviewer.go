package gocodereviewer

import (
	"context"
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

type CoderReviewer struct {
	llm      llms.Model
	template *prompts.PromptTemplate
}

func NewCodeReviewer() (*CoderReviewer, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, err
	}

	template := prompts.NewPromptTemplate(`
You are an expert Go code reviewer. Analyze this Go code for:

1. **Bugs and Logic Issues**: Potential runtime errors, nil pointer dereferences, race conditions
2. **Performance**: Inefficient algorithms, unnecessary allocations, string concatenation issues
3. **Style**: Go idioms, naming conventions, error handling patterns
4. **Security**: Input validation, sensitive data handling

Code to review:
'''go 
{{.code}}
'''

File: {{.filename}}
Provide specific, actionable feedback. For each issue:
- Explain WHY it's a problem
- Show HOW to fix it with code examples
- Rate severity: Critical, Warning, Suggestion

Focus on the most important issues first.`,
		[]string{"code", "filename"})

	return &CoderReviewer{
		llm:      llm,
		template: &template,
	}, nil
}

func (c *CoderReviewer) ReviewFile(fileName string) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file err: %w", err)
	}

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, fileName, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing  go files fail, err: %w", err)
	}

	prompt, err := c.template.Format(map[string]any{
		"code":     string(content),
		"filename": fileName,
	})
	if err != nil {
		return fmt.Errorf("formatting prompt fail, err: %w", err)
	}

	ctx := context.Background()
	response, err := c.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return fmt.Errorf("generating review err: %w", err)
	}

	fmt.Printf("\n==== review for %s ==== \n", fileName)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(response.Choices[0].Content)
	fmt.Println(strings.Repeat("=", 80))
	return nil
}

func ReviewDirectory(reviewer *CoderReviewer, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			return reviewer.ReviewFile(path)
		}
		return nil
	})
}

// level: 0 current modify file, 1: preview-1 git commit files
func ReviewGitChanges(reviewer *CoderReviewer, level int) error {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("HEAD~%d", level))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("get git change files err: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range files {
		if strings.HasSuffix(file, ".go") && file != "" {
			if err := reviewer.ReviewFile(file); err != nil {
				fmt.Printf("review file: %v, err: %v", file, err)
			}
		}
	}
	return nil
}

func RunReviewerCode() {
	var (
		file = flag.String("file", "", "Go file to review")
		dir  = flag.String("dir", "", "Directory to review (all .go files)")
		git  = flag.Int("git", -1, "Review files changed in git commit level for working directory")
	)
	flag.Parse()

	reviewer, err := NewCodeReviewer()
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *file != "":
		if err := reviewer.ReviewFile(*file); err != nil {
			log.Fatal(err)
		}
	case *dir != "":
		if err := ReviewDirectory(reviewer, *dir); err != nil {
			log.Fatal(err)
		}
	case *git >= 0:
		if err := ReviewGitChanges(reviewer, *git); err != nil {
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
