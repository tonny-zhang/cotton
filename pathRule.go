package httpserver

import (
	"httpserver/utils"
	"regexp"
	"strings"
)

var regParam = regexp.MustCompile(":[^/]+")
var regParamDetail = regexp.MustCompile(":([a-z]+)(<num>|{([^}]+)})?")

const (
	ruleParamTypeString = 0
	ruleParamTypeNum    = 1
	ruleParamTypeRegexp = 2
)

type ruleParam struct {
	name string
	t    int
}

type pathRule struct {
	rule                    string
	child                   []string
	countSplit              int
	ruleRegx                *regexp.Regexp
	params                  []ruleParam
	handler                 *HandlerFunc
	middlewareHandlersIndex int
}

type matchResult struct {
	IsMatch bool
	Params  map[string]string
	rule    *pathRule
}

type pathRuleSlice []pathRule

func (s pathRuleSlice) Len() int { return len(s) }

func (s pathRuleSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s pathRuleSlice) Less(i, j int) bool {
	result := s[i].countSplit - s[j].countSplit

	if result == 0 {
		result = len(s[i].child) - len(s[j].child)
	}
	return result > 0
}

func newPathRule(rule string, handler *HandlerFunc) pathRule {
	rule = utils.CleanPath(rule)
	rp := pathRule{
		rule:                    rule,
		handler:                 handler,
		countSplit:              strings.Count(rule, "/"),
		middlewareHandlersIndex: -1,
	}
	rp.child = strings.Split(strings.Trim(rule, "/"), "/")
	if strings.Index(rp.rule, ":") > -1 {
		rp.params = make([]ruleParam, 0)
		ruleRegStr := regParam.ReplaceAllStringFunc(rp.rule, func(str string) string {
			resultDetail := regParamDetail.FindAllStringSubmatch(str, -1)[0]

			// fmt.Println(resultDetail, len(resultDetail))
			regexpStr := resultDetail[3]
			isNum := false
			isRegexp := "" != regexpStr

			if !isRegexp {
				isNum = resultDetail[2] == "<num>"
				if isNum {
					regexpStr = "\\d+"
				}
			}
			rparam := ruleParam{
				name: resultDetail[1],
			}
			if isNum {
				rparam.t = ruleParamTypeNum
			} else if isRegexp {
				rparam.t = ruleParamTypeRegexp
			} else {
				rparam.t = ruleParamTypeString
			}
			rp.params = append(rp.params, rparam)

			if "" == regexpStr {
				regexpStr = "[^/]+"
			}
			return "(" + regexpStr + ")"
		})
		r, e := regexp.Compile(ruleRegStr + "$")
		if e == nil {
			rp.ruleRegx = r
		}
	}
	return rp
}

// isConflictsWith is conflicts with other pathrule
func (pr *pathRule) isConflictsWith(other *pathRule) bool {
	if pr.countSplit == other.countSplit {
		if len(pr.child) != len(other.child) {
			return false
		}
		for i, v := range pr.child {
			ov := other.child[i]
			if (v[0] == ':') || (ov[0] == ':') {
				continue
			} else if v != ov {
				return false
			}
		}
		return true
	}
	return false
}

func (pr *pathRule) match(path string) (result matchResult) {
	path = utils.CleanPath(path)
	if pr.ruleRegx != nil {
		// TODO: 先用/个数去匹配，如果匹配到再处理，以提高效率
		arr := pr.ruleRegx.FindAllStringSubmatch(path, -1)
		if len(arr) > 0 {
			list := arr[0]
			result.Params = make(map[string]string)
			for i, k := range pr.params {
				result.Params[k.name] = list[i+1]
			}
			result.IsMatch = true
		}
	} else {
		result.IsMatch = path == pr.rule || strings.HasPrefix(path, pr.rule)
	}
	if result.IsMatch {
		result.rule = pr
	}
	return
}
