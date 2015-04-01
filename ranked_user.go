package main

type RankedUser struct {
	Name     string
	PRs      []string
	TotalPRs int
}

type RankedUsers []*RankedUser

func (a RankedUsers) Len() int           { return len(a) }
func (a RankedUsers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RankedUsers) Less(i, j int) bool { return a[i].TotalPRs < a[j].TotalPRs }
