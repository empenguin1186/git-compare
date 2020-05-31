package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	// コミットログ取得
	out, _ := exec.Command("git", "log", "-n", "10", "--oneline").Output()
	logs := strings.Split(string(out), "\n")

	// 比較するコミットハッシュ取得
	activity := NewActivity(logs)
	before := activity.ChooseCommit("choose before commit")
	after := activity.ChooseCommit("choose after commit")

	// コミットハッシュのみの配列を作成
	regex := regexp.MustCompile(`^[a-z0-9]*`)
	commitHashs := make([]string, len(logs))
	for i, e := range logs {
		commitHashs[i] = regex.FindStringSubmatch(e)[0]
	}

	// URL生成
	out, _ = exec.Command("git", "config", "--get", "remote.origin.url").Output()
	url := string(out)
	url = strings.Replace(url, ".git\n", "", -1)
	var b strings.Builder
	b.WriteString(url)
	b.WriteString("/compare/")
	b.WriteString(commitHashs[before])
	b.WriteString("..")
	b.WriteString(commitHashs[after])
	result := b.String()
	fmt.Println(result)
}
