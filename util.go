package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Record represents an umeng error record
type Record struct {
	Title      string    `csv:"错误摘要"`
	Count      int       `csv:"发生次数"`
	FirstDate  time.Time `csv:"首次发生时间" format:"2006-01-02 03:04:05"`
	Version    string    `csv:"版本"`
	StackTrace string    `csv:"错误详情"`
}

// read umeng file, need convert utf16-le(with bom) to utf8
func readFile(path string) (reader io.Reader, err error) {
	e := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	file, err := os.Open(path)

	if err != nil {
		return
	}

	reader = transform.NewReader(file, e.NewDecoder())
	return
}

func readReport(input io.Reader) (result []Record, err error) {
	reader := csv.NewReader(input)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 5
	headerLine, err := reader.Read()
	if err != nil {
		return
	}
	fieldColumnIndex, err := getFieldIndexFromHeader(Record{}, headerLine)
	if err != nil {
		return
	}
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	for _, record := range records {
		s := Record{}
		err = fillStruct(&s, fieldColumnIndex, record)
		if err != nil {
			return
		}
		result = append(result, s)
	}
	return
}

func fillStruct(v interface{}, fieldColumnIndex map[int]int, data []string) (err error) {
	reflectStruct := reflect.ValueOf(v).Elem()

	for i := 0; i < reflectStruct.NumField(); i++ {
		field := reflectStruct.Field(i)
		value := data[fieldColumnIndex[i]]

		switch field.Type().Kind() {
		default:
			return fmt.Errorf("not implemented type %s", field.Type().Name())
			break
		case reflect.String:
			field.SetString(value)
			break
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			field.SetInt(int64(intValue))
			break
		case reflect.Struct:
			structName := field.Type().Name()
			switch structName {
			default:
				return fmt.Errorf("not implemented type %s", field.Type().Name())
				break
			case reflect.TypeOf(time.Time{}).Name():
				timeField := reflectStruct.Type().Field(i)
				timeFormat := timeField.Tag.Get("format")
				if timeFormat == "" {
					return fmt.Errorf("missing format key tag for field %s", timeField.Name)
				}
				timeValue, err := time.Parse(timeFormat, value)
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(timeValue))
				break
			}
			break
		}
	}

	return
}

func getFieldIndexFromHeader(v interface{}, header []string) (result map[int]int, err error) {
	headerMap := make(map[string]int, len(header))
	for index, value := range header {
		headerMap[value] = index
	}
	reflectStruct := reflect.TypeOf(v)

	result = make(map[int]int, reflectStruct.NumField())

	for index := 0; index < reflectStruct.NumField(); index++ {
		field := reflectStruct.Field(index)
		key := field.Tag.Get("csv")
		if key == "" {
			key = field.Name
		}
		value, ok := headerMap[key]
		if !ok {
			result = nil
			err = fmt.Errorf("%s not found in header", key)
			return
		}
		result[index] = value
	}
	return
}
