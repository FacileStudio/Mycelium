package memory

import "strings"

// stampPage sets `updated:` to today and `created:` to the old page's, or
// today for a new page or an old one that had none.
func stampPage(content, old, day string) (string, error) {
	stamped, err := bumpUpdated(content, day)
	if err != nil {
		return "", err
	}
	created := scalarOf(old, "created")
	if created == "" {
		created = day
	}
	return setScalar(stamped, "created", created), nil
}

// scalarOf reads one frontmatter scalar of a page, or "" without one.
func scalarOf(content, key string) string {
	if !strings.HasPrefix(content, "---\n") {
		return ""
	}
	end := strings.Index(content[4:], "\n---")
	if end < 0 {
		return ""
	}
	return scalar(content[4:4+end], key)
}

// setScalar writes one frontmatter scalar, adding it when absent. The page
// is known to have closed frontmatter.
func setScalar(content, key, value string) string {
	end := strings.Index(content[4:], "\n---")
	head := strings.Split(content[4:4+end], "\n")
	for i, line := range head {
		if strings.HasPrefix(line, key+":") {
			head[i] = key + ": " + value
			return "---\n" + strings.Join(head, "\n") + content[4+end:]
		}
	}
	return "---\n" + strings.Join(append(head, key+": "+value), "\n") + content[4+end:]
}

// dropPointer removes every index line linking to the page.
func dropPointer(index, page string) string {
	lines := strings.Split(index, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if !strings.Contains(line, "("+page+")") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}
