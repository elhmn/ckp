package fixtures

import (
	"time"

	"github.com/elhmn/ckp/internal/store"
)

//GetListWithMoreThan10Elements returns list of scripts that contains more than 10 elements
func GetListWithMoreThan10Elements() []store.Script {
	creationTimeFix, err := time.Parse(time.RFC1123, "Thu, 13 May 2021 11:36:17 CEST")
	if err != nil {
		return nil
	}

	updateTimeFix, err := time.Parse(time.RFC1123, "Thu, 13 May 2021 11:38:17 CEST")
	if err != nil {
		return nil
	}

	return []store.Script{
		{ID: "ID_1", CreationTime: creationTimeFix, UpdateTime: updateTimeFix, Comment: "comment_1", Code: store.Code{Content: "code_content_1"}},
		{ID: "ID_2", CreationTime: creationTimeFix, UpdateTime: updateTimeFix, Comment: "Comment_2", Code: store.Code{Content: "code_content_2", Alias: "Alias_2"}},
		{ID: "ID_3", CreationTime: creationTimeFix, UpdateTime: updateTimeFix, Code: store.Code{Content: "code_content_3"}},
		{ID: "ID_4", CreationTime: creationTimeFix, UpdateTime: updateTimeFix, Comment: "Comment_4", Solution: store.Solution{Content: "solution_content_4"}},
		{ID: "ID_5", CreationTime: creationTimeFix, UpdateTime: updateTimeFix, Solution: store.Solution{Content: "solution_content_5"}},
	}
}

func GetPrintListWithMoreThan10Elements() string {
	return `ID: ID_1
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Comment: comment_1
  Code: code_content_1

ID: ID_2
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Alias: Alias_2
  Comment: Comment_2
  Code: code_content_2

ID: ID_3
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Code: code_content_3

ID: ID_4
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Solution
  Comment: Comment_4
  Solution: solution_content_4

ID: ID_5
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Solution
  Solution: solution_content_5

`
}

func GetPrintListWithLessThan2Elements() string {
	return `ID: ID_1
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Comment: comment_1
  Code: code_content_1

ID: ID_2
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Alias: Alias_2
  Comment: Comment_2
  Code: code_content_2

`
}

func GetPrintListOnlyCode() string {
	return `ID: ID_1
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Comment: comment_1
  Code: code_content_1

ID: ID_2
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Alias: Alias_2
  Comment: Comment_2
  Code: code_content_2

ID: ID_3
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Code
  Code: code_content_3

`
}

func GetPrintListOnlySolution() string {
	return `ID: ID_4
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Solution
  Comment: Comment_4
  Solution: solution_content_4

ID: ID_5
CreationTime: Thu, 13 May 2021 11:36:17 CEST
UpdateTime: Thu, 13 May 2021 11:38:17 CEST
  Type: Solution
  Solution: solution_content_5

`
}
