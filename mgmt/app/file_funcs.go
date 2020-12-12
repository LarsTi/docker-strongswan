package main

import (
        "fmt"
        "strings"
        "bufio"
        "os"
        "strconv"
)
func getFiles() []string{
        var files []string
        f, err := os.Open("/app/config")
        if err != nil {
                fmt.Println(err)
                return files
        }
        defer f.Close()
        fileInfo, err := f.Readdir(-1)
        if err != nil{
                fmt.Println(err)
                return files
        }
        for _, file := range fileInfo {
                files = append(files, file.Name())
        }
        return files
}
func getStringValueFromPath(path string, value string) string{
        f, err := os.Open(strings.Join([]string{"/app","config",path}, "/"))
        if err != nil{
                fmt.Println(err)
                return ""
        }
        defer func(){
                if err = f.Close(); err != nil{
                        fmt.Println(err)
                }
        }()
        s := bufio.NewScanner(f)
        for s.Scan(){
                if(strings.HasPrefix(s.Text(), value)){
                        return strings.Split(s.Text(), "=")[1]
                }
        }
        return "";
}
func getIntValueFromPath(path string, value string) int {
        i, err := strconv.Atoi(getStringValueFromPath(path, value))
        if err != nil {
                fmt.Printf("Error in file %s, value %s\n", path, value)
                fmt.Println(err)
                return 0
        }
        return i
}
func getStringArrayFromPath(path string, value string) []string {
        return strings.Split(getStringValueFromPath(path, value),",")
}
