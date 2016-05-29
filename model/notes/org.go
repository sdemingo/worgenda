package notes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// These functions parser a plain text content formatted with org-mode
// and return Notes

const (
	ORGDATEHOURFORMAT = "2006-01-02 Mon 15:04"
	ORGDATEFORMAT     = "2006-01-02"
)

var noteTitleReg = regexp.MustCompile("(?m)^(\\*{1,3} .+\\n)")
var separator = "@@@@\n"
var separatorReg = regexp.MustCompile("(?m)^@@@@\\n")

var stampReg = regexp.MustCompile("\\<[^\\>]+\\>")
var dateReg = regexp.MustCompile("\\<\\d{4}-\\d{2}-\\d{2} .{3}\\>")
var hourDateReg = regexp.MustCompile("\\<\\d{4}-\\d{2}-\\d{2} .{3} \\d{2}:\\d{2}\\>")

var anniversaryReg = regexp.MustCompile("\\%\\%\\(diary-anniversary \\d{1,2} \\d{1,2} \\d{4}\\)(.*)")
var deadlineReg = regexp.MustCompile("(DEADLINE:)?( )*\\<\\d{4}-\\d{2}-\\d{2} .{3}( \\d{2}:\\d{2})?( [\\+\\-]+\\d+[a-z])+\\>")
var repetitionStatusReg = regexp.MustCompile("(.)*- State \"DONE\"(.)*from \"TODO\"(.)*")
var propertiesGroupReg = regexp.MustCompile("(?s):PROPERTIES:(.)*:END:")
var warningReg = regexp.MustCompile(" \\-\\d+[a-z]")

var yearMonthDayReg = regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
var monthDayAnniversaryReg = regexp.MustCompile("\\d{1,2} \\d{1,2} \\d{4}")
var hourMinReg = regexp.MustCompile("\\d{2}:\\d{2}")

var dayDuration = time.Duration(24) * time.Hour
var weekDuration = time.Duration(7) * dayDuration
var monthDuration = time.Duration(4) * weekDuration
var yearDuration = time.Duration(52) * weekDuration

// Parse string content in org mode and recover notes from it
func Parse(content string) []Note {
	notes := make([]Note, 0)

	content = noteTitleReg.ReplaceAllString(content, separator+"$1")
	rawNotes := separatorReg.Split(content, -1)

	for _, rnote := range rawNotes {
		note := NewNote()
		note.Status = parseStatus(rnote)
		note.Title = parseTitle(rnote)
		note.Body = parseBody(rnote)
		note.Stamps = parseStamps(rnote)
		note.Deadline = parseDeadlines(rnote)
		note.Warning = parseDeadlineWarning(rnote)

		notes = append(notes, *note)
	}

	return notes
}

func parseDeadlines(orgnote string) time.Time {
	deadline := nullTime
	if deadlineReg.FindString(orgnote) != "" {
		dl := yearMonthDayReg.FindAllString(orgnote, 1)
		deadline, _ = time.Parse(ORGDATEFORMAT, dl[0])
		return deadline
	}
	return deadline
}

func parseDeadlineWarning(orgnote string) time.Duration {

	if dl := deadlineReg.FindString(orgnote); dl != "" {
		w := warningReg.FindString(dl)
		numberReg := regexp.MustCompile("\\d+")
		durReg := regexp.MustCompile("[a-z]")
		numberStr := numberReg.FindString(w)
		durStr := durReg.FindString(w)
		number, _ := strconv.ParseInt(numberStr, 10, 64)

		switch durStr {
		case "h":
			return time.Duration(number) * time.Hour
		case "d":
			return time.Duration(number) * dayDuration
		case "m":
			return time.Duration(number) * monthDuration
		case "w":
			return time.Duration(number) * weekDuration
		case "y":
			return time.Duration(number) * yearDuration
		}
	}

	return time.Duration(0)
}

func parseStatus(orgnote string) string {
	title := noteTitleReg.FindString(orgnote)
	prefix := regexp.MustCompile(" TODO|DONE ")
	return strings.TrimSpace(prefix.FindString(title))
}

// Clean the title of the note
func parseTitle(orgnote string) string {
	title := noteTitleReg.FindString(orgnote)
	prefix := regexp.MustCompile("(?m)^\\*+( TODO| DONE)?( \\[#(A|B|C)\\])?")
	return strings.Trim(prefix.ReplaceAllString(title, ""), " \n\t")
}

// Extract and clean the body note. If the body content a
// %%(diary-anniversary d m year) function it calculates the
// difference and print the number in the body
func parseBody(orgnote string) string {

	body := orgnote

	if anniversaryReg.MatchString(body) {
		mdy := strings.Fields(monthDayAnniversaryReg.FindString(body))
		if len(mdy) < 3 {
			return body
		}
		ymd := fmt.Sprintf("%s-%02s-%02s", mdy[2], mdy[0], mdy[1])
		ymd2 := fmt.Sprintf("%d-%02s-%02s", time.Now().Year(), mdy[0], mdy[1])
		t, err := time.Parse("2006-01-02", ymd)
		t2, err := time.Parse("2006-01-02", ymd2)
		if err != nil {
			return body
		}

		body = anniversaryReg.ReplaceAllString(orgnote, "$1")
		diffYears := int(t2.Sub(t).Hours() / 8760)
		body = strings.Replace(body, "%d", fmt.Sprintf("%d", diffYears), 1)
		body = fmt.Sprintf("%s\n", body)
	}

	body = deadlineReg.ReplaceAllString(body, "")
	//body = repetitionStatusReg.ReplaceAllString(body, "")
	body = propertiesGroupReg.ReplaceAllString(body, "")
	body = noteTitleReg.ReplaceAllString(body, "")
	//body = dateReg.ReplaceAllString(body, "")
	//body = hourDateReg.ReplaceAllString(body, "")

	return strings.Trim(body, " \n\t")
}

// Extract the stamps of the note in several formats: <day-month-year>,
// <day-month-year weekday hour:min> and <day-month-year weekday hour:min repetition>.
// If the stamp has not hour is dangerous for
// datesInSameDay() function that we save it as 00:00. For that reason we
// add a second after the midnight

func parseStamps(orgnote string) []time.Time {

	times := make([]time.Time, 0)

	// extract normal stamp with hour or not or with period
	rawTimes := stampReg.FindAllString(orgnote, -1)
	for _, rt := range rawTimes {
		if dateReg.MatchString(rt) {
			ymd := yearMonthDayReg.FindString(rt)
			t, err := time.Parse("2006-01-02", ymd)
			t = t.Add(time.Second)
			if err == nil {
				times = append(times, t)
			}
			continue
		}

		if hourDateReg.MatchString(rt) || deadlineReg.MatchString(rt) { //repetitionReg.MatchString(rt) {
			ymd := yearMonthDayReg.FindString(rt)
			hm := hourMinReg.FindString(rt)
			var t time.Time
			var err error
			if hm != "" {
				t, err = time.Parse("2006-01-02 15:04", ymd+" "+hm)
			} else {
				t, err = time.Parse("2006-01-02", ymd)
				t = t.Add(time.Second)
			}
			if err == nil {
				times = append(times, t)
			}
			continue
		}

	}

	// extract anniversary dates format
	rawTimes = anniversaryReg.FindAllString(orgnote, -1)
	for _, rt := range rawTimes {
		mdy := strings.Fields(monthDayAnniversaryReg.FindString(rt))
		if len(mdy) < 3 {
			continue
		}
		year := time.Now().Year()
		ymd := fmt.Sprintf("%d-%02s-%02s", year, mdy[0], mdy[1])
		t, err := time.Parse("2006-01-02", ymd)
		t = t.Add(time.Second)
		if err == nil {
			times = append(times, t)
		}
	}

	return times
}

/*
 HTML conversion
*/

var head1Reg = regexp.MustCompile("(?m)^\\*(?P<status> TODO| DONE)? (?P<head>.+)$")
var head2Reg = regexp.MustCompile("(?m)^\\*\\*(?P<status> TODO| DONE)? (?P<head>.+)$")
var head3Reg = regexp.MustCompile("(?m)^\\*\\*\\* (?P<head>.+)$")
var head4Reg = regexp.MustCompile("(?m)^\\*\\*\\*\\* (?P<head>.+)$")
var tagsReg = regexp.MustCompile("(?m)^(?P<head>\\*+ .+)\\s*:(?P<tags>.+):$")

var linkReg = regexp.MustCompile("\\[\\[(?P<url>[^\\]]+)\\]\\[(?P<text>[^\\]]+)\\]\\]")
var imgLinkReg = regexp.MustCompile("\\[\\[file:\\.\\./img/(?P<img>[^\\]]+)\\]\\[file:\\.\\./img/(?P<thumb>[^\\]]+)\\]\\]")
var imgReg = regexp.MustCompile("\\[\\[\\.\\./img/(?P<src>[^\\]]+)\\]\\]")

var codeReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_SRC \\w*\\n(?P<code>(?s)[^\\#]+)^\\#\\+END_SRC\\n")
var codeHeaderReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_SRC \\w*\\n")
var codeFooterReg = regexp.MustCompile("(?m)^\\#\\+END_SRC\\n")

var quoteReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_QUOTE\\s*\\n(?P<cite>(?s)[^\\#]+)^\\#\\+END_QUOTE\\n")
var quoteHeaderReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_QUOTE\\s*\\n")
var quoteFooterReg = regexp.MustCompile("(?m)^\\#\\+END_QUOTE\\n")

var centerReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_CENTER\\s*\\n(?P<cite>(?s)[^\\#]+)^\\#\\+END_CENTER\\n")
var centerHeaderReg = regexp.MustCompile("(?m)^\\#\\+BEGIN_CENTER\\s*\\n")
var centerFooterReg = regexp.MustCompile("(?m)^\\#\\+END_CENTER\\n")

var parReg = regexp.MustCompile("\\n\\n+(?P<text>[^\\n]+)")
var allPropsReg = regexp.MustCompile(":PROPERTIES:(?s).+:END:")
var propReg = regexp.MustCompile("(?m)^\\#\\+.*$")
var rawHTML = regexp.MustCompile("\\<A-Za-z[^\\>]+\\>")

//estilos de texto
var boldReg = regexp.MustCompile("(?P<prefix>[\\s|\\W]+)\\*(?P<text>[^\\s][^\\*]+)\\*(?P<suffix>[\\s|\\W]*)")
var italicReg = regexp.MustCompile("(?P<prefix>[\\s])/(?P<text>[^\\s][^/]+)/(?P<suffix>[^A-Za-z0-9]*)")
var ulineReg = regexp.MustCompile("(?P<prefix>[\\s|\\W]+)_(?P<text>[^\\s][^_]+)_(?P<suffix>[\\s|\\W]*)")
var codeLineReg = regexp.MustCompile("(?P<prefix>[\\s|\\W]+)=(?P<text>[^\\s][^\\=]+)=(?P<suffix>[\\s|\\W]*)")

func Org2HTML(content []byte, url string) string {

	// First remove all HTML raw tags for security
	out := rawHTML.ReplaceAll(content, []byte(""))

	// tags (remove tags)
	out = tagsReg.ReplaceAll(out, []byte(""))

	// headings (h1 is not admit in the post body)
	out = head1Reg.ReplaceAll(out, []byte("<h1 class=\"note $status\">$head</h1>"))
	out = head2Reg.ReplaceAll(out, []byte("<h2 class=\"note $status\">$head</h2>\n"))
	out = head3Reg.ReplaceAll(out, []byte("<h3 class=\"note $status\">$head</h3>\n"))
	out = head4Reg.ReplaceAll(out, []byte("<h4 class=\"note $status\">$head</h4>\n"))
	out = regexp.MustCompile("class=\"(.*)TODO(.*)\">").ReplaceAll(out, []byte("class=\"$1 todo $2\">"))
	out = regexp.MustCompile("class=\"(.*)DONE(.*)\">").ReplaceAll(out, []byte("class=\"$1 done $2\">"))

	// images
	out = imgReg.ReplaceAll(out, []byte("<a target=\"_blank\" href='"+url+"/img/$src'><img src='"+url+"/img/thumbs/$src'/></a>"))
	out = imgLinkReg.ReplaceAll(out, []byte("<a target=\"_blank\" href='"+url+"/img/$img'><img src='"+url+"/img/thumbs/$thumb'/></a>"))
	out = linkReg.ReplaceAll(out, []byte("<a target=\"_blank\" href='$url'>$text</a>"))

	// Extract blocks codes
	codeBlocks, out := extractBlocks(string(out),
		codeReg,
		codeHeaderReg,
		codeFooterReg,
		"code")

	quoteBlocks, out := extractBlocks(string(out),
		quoteReg,
		quoteHeaderReg,
		quoteFooterReg,
		"quote")

	out = centerReg.ReplaceAll(out, []byte("<span class=\"center\">$cite</span>\n"))
	out = parReg.ReplaceAll(out, []byte("\n\n<p/>$text"))
	out = allPropsReg.ReplaceAll(out, []byte("\n"))
	out = propReg.ReplaceAll(out, []byte("\n"))

	// font styles
	out = italicReg.ReplaceAll(out, []byte("$prefix<i>$text</i>$suffix"))
	out = boldReg.ReplaceAll(out, []byte("$prefix<b>$text</b>$suffix"))
	out = ulineReg.ReplaceAll(out, []byte("$prefix<u>$text</u>$suffix"))
	out = codeLineReg.ReplaceAll(out, []byte("$prefix<code>$text</code>$suffix"))

	// Reinsert block codes
	sout := string(out)
	sout = strings.Replace(sout, "\n", "<br>", -1)

	sout = insertBlocks(sout, codeBlocks, "<pre><code>", "</code></pre>", "code")
	sout = insertBlocks(sout, quoteBlocks, "<blockquote>", "</blockquote>", "quote")

	return sout
}

func extractBlocks(src string, fullReg, headerReg, footerReg *regexp.Regexp, name string) ([][]byte, []byte) {
	out := []byte(src)
	blocks := fullReg.FindAll(out, -1)
	for i := range blocks {
		bstring := string(blocks[i])
		blocks[i] = headerReg.ReplaceAll(blocks[i], []byte("\n"))
		blocks[i] = footerReg.ReplaceAll(blocks[i], []byte("\n"))
		out = []byte(strings.Replace(string(out), bstring, "block"+name+":"+strconv.Itoa(i)+"\n", 1))
	}

	return blocks, out
}

func insertBlocks(src string, blocks [][]byte, header, footer, name string) string {
	s := src
	for i := range blocks {
		bstring := string(blocks[i])
		s = strings.Replace(s, "block"+name+":"+strconv.Itoa(i)+"\n",
			header+bstring+footer+"\n", 1)
	}

	return s
}
