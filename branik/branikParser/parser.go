package branikParser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	branik       = 1
	branikBalik  = 6 * branik
	branikPaleta = branikBalik * 48

	branikCena       = 38
	branikBalikCena  = 6 * branikCena
	branikPaletaCena = branikBalikCena * 48

	mainRegexStr  = "[1-9][0-9]*[ ]{0,}(?:kc|kč|czk|mega|korun|euro|eur|k|[ ]{0,}(czk|euro|eur))"
	digitRegexStr = "\\d"
	unitRegexStr  = "[a-z]"

	mega = 1000000
	euro = 25
	k    = 1000

	resultsDigits []string
	resultsUnits  []string
)

func createResponse() string {
	response := ""
	num := 0
	for i := 0; i < len(resultsUnits); i++ {
		num, _ = strconv.Atoi(resultsDigits[i])
		if resultsUnits[i] == "mega" {
			num *= mega
		}
		if resultsUnits[i] == "euro" {
			num *= euro
		}
		if resultsUnits[i] == "k" {
			num *= k
		}

		count := num / branikCena

		response += "-----------------------\n"
		response += bottleMessage(i, count)
		if count >= branikBalik {
			response += balikMessage(i, count/branikBalik)
		}
		if count >= branikPaleta {
			response += paletaMessage(i, count/branikPaleta)
		}
		response += "-----------------------\n"
	}
	return response
}

func paletaMessage(index int, count int) string {
	paletaSuffix := "eta"
	if count > 1 && count < 5 {
		paletaSuffix = "ety"
	} else {
		paletaSuffix = "iet"
	}

	return fmt.Sprintf("respektive %d pal%s Branika\n", count, paletaSuffix)
}

func balikMessage(index int, count int) string {
	balikSuffix := ""
	if count < 5 && count > 1 {
		balikSuffix = "y"
	} else {
		balikSuffix = "ov"
	}
	return fmt.Sprintf("a to je %d balik%s Branika\n", count, balikSuffix)
}

func bottleMessage(index int, count int) string {
	bottleSuffix := "iek"
	countBottle := strconv.Itoa(count)
	verb := "sa da kupit"
	if count <= 1 {
		if count == 0 {
			verb = "nekupis ani "
			countBottle = "jednu"
			bottleSuffix = "ku"
		} else {
			bottleSuffix = "ka"
		}
	} else if count <= 5 {
		bottleSuffix = "ky"
	}
	return fmt.Sprintf("Za %s %s %s %s dvoulitrov%s Branika\n ", resultsDigits[index], resultsUnits[index], verb, countBottle, bottleSuffix)
}

func Parse(messageContent string) string {
	mainRegex, _ := regexp.Compile(mainRegexStr)
	digitRegex, _ := regexp.Compile(digitRegexStr)
	unitRegex, _ := regexp.Compile(unitRegexStr)

	content := strings.ToLower(messageContent)

	resultsMainBytes := mainRegex.FindAll([]byte(content), -1)
	var resultsDigitBytes [][]string
	var resultsUnitsBytes [][]string

	for i := 0; i < len(resultsMainBytes); i++ {
		resultsDigitBytes = append(resultsDigitBytes, digitRegex.FindAllString(string(resultsMainBytes[i]), -1))
		resultsUnitsBytes = append(resultsUnitsBytes, unitRegex.FindAllString(string(resultsMainBytes[i]), -1))
		resultsDigits = append(resultsDigits, strings.Join(resultsDigitBytes[i], ""))
		resultsUnits = append(resultsUnits, strings.Join(resultsUnitsBytes[i], ""))
	}

	response := createResponse()
	resultsDigits = nil
	resultsUnits = nil

	return response
}
