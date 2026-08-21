package tenets_test

import (
	"os"
	"strings"
	"testing"
)

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func section(md, name string) string {
	needle := "## " + name
	i := strings.Index(md, needle)
	if i < 0 {
		return ""
	}
	rest := md[i:]
	nl := strings.Index(rest, "\n")
	if nl < 0 {
		return ""
	}
	body := rest[nl+1:]
	if j := strings.Index(body, "\n## "); j >= 0 {
		body = body[:j]
	}
	return body
}

func TestPickerAffirmativeButtonIsAnAnswer(t *testing.T) {
	first := section(mustRead(t, "prompt.md"), "First turn")
	if first == "" {
		t.Fatal("missing First turn section")
	}
	if !strings.Contains(first, "Shall I proceed?") {
		t.Error("keep Shall I proceed? as the picker prompt")
	}
	if strings.Contains(first, `buttons are "Shall I proceed?"`) ||
		strings.Contains(first, "buttons are Shall I proceed?") {
		t.Error("the affirmative picker button must be an answer, not the question; a host that records the label would then fail the yes-gate")
	}
	hasYes := strings.Contains(first, `buttons are Yes`) ||
		strings.Contains(first, `buttons are "Yes"`) ||
		strings.Contains(first, `The buttons are Yes`) ||
		strings.Contains(first, `The buttons are "Yes"`)
	hasProceed := strings.Contains(first, `buttons are Proceed`) ||
		strings.Contains(first, `buttons are "Proceed"`)
	if !hasYes && !hasProceed {
		t.Error(`affirmative picker button must be labeled Yes or Proceed`)
	}
	if !strings.Contains(first, "Not yet") {
		t.Error("negative picker button must stay Not yet")
	}
}

func TestFirstTurnHasNoInterviewQuestion(t *testing.T) {
	prompt := mustRead(t, "prompt.md")
	if strings.Contains(prompt, "Ask one question") {
		t.Error("Ask one question reintroduces an interview on the gated turn; missing pieces belong in the summary, which still ends on Shall I proceed?")
	}
	first := section(prompt, "First turn")
	if first == "" {
		t.Fatal("missing First turn section")
	}
	if !strings.Contains(first, "Shall I proceed?") {
		t.Error("first turn must end on Shall I proceed?")
	}
	if strings.Contains(first, "write the set") && !strings.Contains(first, "never a working set") {
		t.Error("first turn must not write the set")
	}
}

func TestAspirationalMeansTheyDoNotLiveItYet(t *testing.T) {
	prompt := mustRead(t, "prompt.md")
	if strings.Contains(prompt, "Aspirational is reality") {
		t.Error(`"Aspirational is reality" fights "even if it does not yet"; Aspirational is that they do not live it yet`)
	}
	if !strings.Contains(prompt, "they do not live it yet") &&
		!strings.Contains(prompt, "do not live it yet") {
		t.Error("Aspirational must mean they do not live it yet, not intended operating behavior described as reality")
	}
	if !strings.Contains(prompt, "Foundational is role") {
		t.Error("Foundational is role, a separate axis from Aspirational")
	}
}

func TestReadmeNamesGuessedTenetAndRejectedTenets(t *testing.T) {
	gotBack := section(mustRead(t, "README.md"), "What you should get back")
	if gotBack == "" {
		t.Fatal("missing What you should get back section")
	}
	if !strings.Contains(gotBack, "guessed tenet") {
		t.Error("first-turn summary in README must include any guessed tenet")
	}
	if !strings.Contains(gotBack, "tenets considered and rejected") {
		t.Error(`README must keep "tenets considered and rejected"`)
	}
}
