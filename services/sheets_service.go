package services

import (
	"ClubTennis/models"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// returns os temp directory + /users.xlsx
func DEFAULT_SHEET_FILENAME() string {
	return os.TempDir() + "/users.xlsx"
}

const defaultSheetName string = "Sheet1"

// generates a file in /tmp (or OS temp directory) containing the users in .xlsx format. returns filename.
func UsersToSheet(users []models.User) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	f.SetActiveSheet(0)
	sort.Slice(users, func(i, j int) bool {
		return strings.Compare(users[i].LastName, users[j].LastName) < 0
	})

	valids := getValidTypes()
	go makeTitle(f, users, valids)
	makeBody(f, users, valids)

	filename := DEFAULT_SHEET_FILENAME()
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
		return "", err
	}
	return filename, nil
}

// builds an array of users using the .xlsx filename provided.
func SheetToUsers(filename string) ([]models.User, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	var m []models.User

	rows, err := f.GetRows(defaultSheetName)
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(rows); i++ {
		m = append(m, interpretUser(rows[i]))
	}
	return m, nil
}

func coord(col, row int) (s string) {
	s, e := excelize.CoordinatesToCellName(col, row)
	if e != nil {
		panic("fatal error parsing spreadsheet")
	}
	return
}

func getValidTypes() map[reflect.Kind]bool {
	v := make(map[reflect.Kind]bool)
	v[reflect.Int] = true
	v[reflect.Int8] = true
	v[reflect.Int16] = true
	v[reflect.Int32] = true
	v[reflect.Int64] = true
	v[reflect.Uint] = true
	v[reflect.Uint8] = true
	v[reflect.Uint16] = true
	v[reflect.Uint32] = true
	v[reflect.Uint64] = true
	v[reflect.Float32] = true
	v[reflect.Float64] = true
	v[reflect.String] = true
	v[reflect.Bool] = true
	return v
}

func makeTitle(f *excelize.File, users []models.User, valids map[reflect.Kind]bool) {
	val := reflect.ValueOf(users[0])
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	for i := 0; i < val.NumField(); i++ {
		if valids[val.Field(i).Kind()] {
			c := coord(i+1, 1)
			f.SetCellValue(defaultSheetName, c, val.Type().Field(i).Name)
			f.SetCellStyle(defaultSheetName, c, c, style)
		}
	}
}
func makeBody(f *excelize.File, users []models.User, valids map[reflect.Kind]bool) {
	var wg sync.WaitGroup
	for i, u := range users {
		wg.Add(1)
		go func(u models.User, i int) {
			v := reflect.ValueOf(u)
			for j := 0; j < v.NumField(); j++ {
				if valids[v.Field(j).Kind()] {
					f.SetCellValue(defaultSheetName, coord(j+1, i+2), v.Field(j).Interface())
				}
			}
			wg.Done()
		}(u, i)
	}
	wg.Wait()
}
func interpretUser(row []string) models.User {
	userType := reflect.TypeOf(models.User{})
	userPtr := reflect.New(userType)

	buildStruct(userType, userPtr, row)

	userVal := userPtr.Elem()
	userInterface := userVal.Interface()
	return userInterface.(models.User)
}

func buildStruct(t reflect.Type, v reflect.Value, valRow []string) {
	for i := 0; i < v.Elem().NumField(); i++ {
		f := v.Elem().Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.String:
			f.SetString(valRow[i])
		case reflect.Bool:
			f.SetBool(valRow[i] == "TRUE" || valRow[i] == "true" || valRow[i] == "True")
		case reflect.Uint:
			num, _ := strconv.ParseUint(valRow[i], 10, 0)
			f.SetUint(num)
		case reflect.Uint8:
			num, _ := strconv.ParseUint(valRow[i], 10, 8)
			f.SetUint(num)
		case reflect.Uint16:
			num, _ := strconv.ParseUint(valRow[i], 10, 16)
			f.SetUint(num)
		case reflect.Uint32:
			num, _ := strconv.ParseUint(valRow[i], 10, 32)
			f.SetUint(num)
		case reflect.Uint64:
			num, _ := strconv.ParseUint(valRow[i], 10, 64)
			f.SetUint(num)
		case reflect.Int:
			num, _ := strconv.ParseInt(valRow[i], 10, 0)
			f.SetInt(num)
		case reflect.Int8:
			num, _ := strconv.ParseInt(valRow[i], 10, 8)
			f.SetInt(num)
		case reflect.Int16:
			num, _ := strconv.ParseInt(valRow[i], 10, 16)
			f.SetInt(num)
		case reflect.Int32:
			num, _ := strconv.ParseInt(valRow[i], 10, 32)
			f.SetInt(num)
		case reflect.Int64:
			num, _ := strconv.ParseInt(valRow[i], 10, 64)
			f.SetInt(num)
		case reflect.Float32:
			num, _ := strconv.ParseFloat(valRow[i], 32)
			f.SetFloat(num)
		case reflect.Float64:

			num, _ := strconv.ParseFloat(valRow[i], 64)
			f.SetFloat(num)
		default:
		}
	}
}
