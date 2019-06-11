package diamond_client

import (
    "errors"
    . "github.com/CharLemAznable/gokits"
    "os"
    "strings"
)

var DefHolderPrefix = "${"
var DefHolderSuffix = "}"
var DefHolderPrefixBytes = []byte(DefHolderPrefix)
var DefHolderSuffixBytes = []byte(DefHolderSuffix)
var DefHolderPrefixLen = len(DefHolderPrefixBytes)
var DefHolderSuffixLen = len(DefHolderSuffixBytes)

func Substitute(strVal string, ignoreBadHolders bool, group string, dataId string, lastProperties *Properties) (string, error) {
    if "" == strVal {
        return "", nil
    }

    visitedHolders := NewHashset()
    return SubstituteRecursive(strVal, visitedHolders, ignoreBadHolders, group, dataId, lastProperties)
}

func SubstituteRecursive(strVal string, visitedHolders *Hashset, ignoreBadHolders bool, group string, dataId string, lastProperties *Properties) (string, error) {
    if "" == strVal {
        return "", nil
    }

    buf := string(strVal)
    startIndex := strings.Index(strVal, DefHolderPrefix)
    for -1 != startIndex {
        endIndex := findHolderEndIndex(buf, startIndex)
        if -1 != endIndex {
            holder := buf[startIndex+DefHolderPrefixLen : endIndex]
            defValue := ""
            defIndex := strings.LastIndex(holder, ":")
            if defIndex >= 0 {
                defValue = strings.TrimSpace(holder[defIndex+1:])
                holder = strings.TrimSpace(holder[:defIndex])
            }

            if !visitedHolders.Add(holder) {
                return "", errors.New("Circular PlaceHolder reference '" + holder + "' in property definitions")
            }

            // Recursive invocation, parsing Holders contained in the Holder key.
            holder, err := SubstituteRecursive(holder, visitedHolders, ignoreBadHolders, group, dataId, lastProperties)
            if nil != err {
                return "", err
            }
            // Now obtain the value for the fully resolved key...
            propVal, err := resolveHolder(strVal, holder, defValue, group, dataId, lastProperties)
            if nil != err {
                return "", err
            }
            if "" != propVal {
                // Recursive invocation, parsing Holders contained in the
                // previously resolved Holder value.
                propVal, err = SubstituteRecursive(propVal, visitedHolders, ignoreBadHolders, group, dataId, lastProperties)
                if nil != err {
                    return "", err
                }
                buf = strings.Replace(buf, buf[startIndex:endIndex+DefHolderSuffixLen], propVal, 1)
                propEndIndex := startIndex + len(propVal)
                startIndex = strings.Index(buf[propEndIndex:], DefHolderPrefix)
                if -1 != startIndex {
                    startIndex = startIndex + propEndIndex
                }
            } else if ignoreBadHolders {
                // Proceed with unprocessed value.
                startIndex = strings.Index(buf[endIndex+DefHolderSuffixLen:], DefHolderPrefix)
                if -1 != startIndex {
                    startIndex = startIndex + endIndex + DefHolderSuffixLen
                }
            } else {
                return "", errors.New("Could not resolve Placeholder '" + holder + "'")
            }
            visitedHolders.Remove(holder)
        }
    }
    return buf, nil
}

func findHolderEndIndex(buf string, startIndex int) int {
    index := startIndex + DefHolderPrefixLen
    withinNestedHolder := 0
    bufBytes := []byte(buf)
    for index < len(bufBytes) {
        if SubstringMatch(bufBytes, index, DefHolderSuffixBytes) {
            if withinNestedHolder > 0 {
                withinNestedHolder--
                index = index + DefHolderSuffixLen
            } else {
                return index
            }
        } else if SubstringMatch(bufBytes, index, DefHolderPrefixBytes) {
            withinNestedHolder++
            index = index + DefHolderPrefixLen
        } else {
            index++
        }
    }
    return -1
}

func SubstringMatch(bytes []byte, index int, subbytes []byte) bool {
    for j := 0; j < len(subbytes); j++ {
        i := index + j
        if i >= len(bytes) || bytes[i] != subbytes[j] {
            return false
        }
    }
    return true
}

func resolveHolder(strVal string, holder string, defaultValue string, group string, dataId string, lastProperties *Properties) (string, error) {
    propVal, err := resolveHolderInternal(strVal, holder, defaultValue, group, dataId, lastProperties)
    if nil != err {
        return "", err
    }
    if "" == propVal {
        propVal = resolveSystemProperty(holder)
    }
    return propVal, nil
}

func resolveHolderInternal(strVal string, holder string, defaultValue string, curGroup string, curDataId string, lastProperties *Properties) (string, error) {
    separated := strings.IndexByte(holder, '^')
    if separated < 0 && strings.HasPrefix(holder, "this.") {
        referKey := holder[len("this."):]
        return recursiveSubstitute(strVal, defaultValue, curGroup, curDataId, lastProperties, referKey)
    }

    group := DefaultGroupName
    dataId := holder
    if separated > 0 && separated < len(holder)-1 {
        group = holder[0:separated]
        dataId = holder[separated:]
    }

    value := ""
    propKeyPos := strings.IndexByte(dataId, '^')
    if propKeyPos > 0 {
        subDataId := dataId[0:propKeyPos]
        propKey := dataId[propKeyPos+1:]
        if isSameGroupAndDataId(curGroup, curDataId, group, subDataId) {
            res, err := recursiveSubstitute(strVal, defaultValue, curGroup, curDataId, lastProperties, propKey)
            if nil != err {
                return "", err
            }
            value = res
        } else {
            minerable, err := NewMiner().GetMiner(group, subDataId)
            value = Condition(nil == err, minerable.GetString(propKey), "").(string)
        }
    } else {
        if isSameGroupAndDataId(curGroup, curDataId, group, dataId) {
            // reference to itself
            return "", errors.New(curGroup + "^" + curDataId + "can not refer itself")
        }

        value, _ = NewMiner().GetStoneOrDefault(group, dataId, "")
    }

    return Condition("" != value, value, defaultValue).(string), nil
}

func recursiveSubstitute(strVal string, defaultValue string, curGroup string, curDataId string, lastProperties *Properties, referKey string) (string, error) {
    properties := Condition(nil != lastProperties, lastProperties, ParseStoneToProperties(strVal)).(*Properties)
    property := properties.GetPropertyDefault(referKey, defaultValue)
    substitute, err := Substitute(property, true, curGroup, curDataId, properties)
    properties.SetProperty(referKey, substitute)
    return substitute, err
}

func isSameGroupAndDataId(curGroup string, curDataId string, group string, dataId string) bool {
    return group == curGroup && dataId == curDataId
}

func resolveSystemProperty(key string) string {
    return os.Getenv(key)
}
