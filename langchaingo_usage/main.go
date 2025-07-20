package main

import (
	"flag"

	gcr "github.com/achilsh/go-dev-tools-collection/langchaingo_usage/go_code_reviewer"
)

func main() {
	gcr.InitCodeReview()
	gcr.InitCodereviewStructOutput()

	flag.Parse()

	gcr.RunReviewerCode()
	gcr.RunReviewercodeStructOutput()

}
