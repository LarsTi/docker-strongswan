package filewrapper

import (
        "fmt"
        "strings"
        "bufio"
        "os"
        "strconv"
)
func getAllFiles() []string{
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
func GetFilesForSecrets() []string{
	var files []string
	for _, file := range getAllFiles() {
		if strings.HasSuffix(file, ".secret"){
			files = append(files, file)
		}
	}
	return files
}
func GetFilesForConnections() []string{
	var files []string
	for _, file := range getAllFiles() {
		if strings.HasSuffix(file, ".secret"){
			continue
		}
		files = append(files, file)
	}
	return files
}
func GetStringValueFromPath(path string, value string) string{
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
func GetIntValueFromPath(path string, value string) int {
        i, err := strconv.Atoi(GetStringValueFromPath(path, value))
        if err != nil {
                fmt.Printf("Error in file %s, value %s\n", path, value)
                fmt.Println(err)
                return 0
        }
        return i
}
func GetStringArrayFromPath(path string, value string) []string {
        return strings.Split(GetStringValueFromPath(path, value),",")
}
