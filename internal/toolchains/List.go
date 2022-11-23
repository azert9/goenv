package toolchains

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var goVersionRegexp = regexp.MustCompile(`^go(\d{1,3})\.(\d{1,3})(\.(\d{1,3}))?$`)

type versionTriplet struct {
	Major int
	Minor int
	Patch int
}

func (v *versionTriplet) String() string {

	if v.Patch == -1 {
		return fmt.Sprintf("go%d.%d", v.Major, v.Minor)
	} else {
		return fmt.Sprintf("go%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

type versionTripletList []versionTriplet

func (v versionTripletList) Len() int {
	return len(v)
}

func (v versionTripletList) Less(i, j int) bool {

	if v[i].Major < v[j].Major {
		return true
	}

	if v[i].Minor < v[j].Minor {
		return true
	}

	if v[i].Patch < v[j].Patch {
		return true
	}

	return false
}

func (v versionTripletList) Swap(i, j int) {
	tmp := v[i]
	v[i] = v[j]
	v[j] = tmp
}

// sortVersions will create a sorted slice of toolchain versions. Versions that cannot be parsed are kept in the result
// but are sorted before the parsed ones.
func sortVersions(versions []string) []string {

	// TODO: implement and use

	result := make([]string, 0, len(versions))

	parsed := make([]versionTriplet, 0, len(versions))

	for _, version := range versions {

		match := goVersionRegexp.FindStringSubmatch(version)

		if match == nil {
			result = append(result, version)
			continue
		}

		var triplet versionTriplet

		if val, err := strconv.ParseInt(match[1], 10, 16); err != nil {
			panic(err) // should have been caught by the regexp
		} else {
			triplet.Major = int(val)
		}

		if val, err := strconv.ParseInt(match[2], 10, 16); err != nil {
			panic(err) // should have been caught by the regexp
		} else {
			triplet.Minor = int(val)
		}

		if val, err := strconv.ParseInt(match[4], 10, 16); err != nil {
			// no patch version specified
			triplet.Patch = -1
		} else {
			triplet.Patch = int(val)
		}

		parsed = append(parsed, triplet)
	}

	sort.Sort(versionTripletList(parsed))

	for i := range parsed {
		result = append(result, parsed[i].String())
	}

	return result
}

func List() ([]string, error) {

	dirs, err := getDirs(false)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	entries, err := os.ReadDir(dirs.ToolchainsDir)
	if err != nil {
		return nil, err
	}

	results := make([]string, len(entries))

	for i := range entries {
		results[i] = entries[i].Name()
	}

	return sortVersions(results), nil
}
