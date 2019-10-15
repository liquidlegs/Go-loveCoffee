package main

import (
  "fmt"
  "math/rand"
  "time"
  "os/exec"
  "bytes"
  "strings"
  "log"
  "encoding/base64"
  "net/http"
  "io"
  "os"
  "syscall"
)

// converts string into bytes and encodes with base64
func encodeBase64(encodeTo64 string) string {
  encodeString := base64.StdEncoding.EncodeToString([]byte(encodeTo64))
  returnStr := string(encodeString)
  return returnStr
}

// takes an exe that needs to run through cmd, along with a string containing the path
func runExe(ap string, path string) string {
  decodeAp := decodeBase64(ap)
  decodePath := decodeBase64(path)
  rStringAp := string(decodeAp)
  rStringPath := string(decodePath)
  exeCmd := exec.Command("cmd.exe", "/C", rStringAp, rStringPath)
  exeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
  
  exeCmd.Stdin = strings.NewReader("")
  var output bytes.Buffer
  exeCmd.Stdout = &output
  // Go needs to know how it can throw an error, otherwise it wont let you compile
  err := exeCmd.Run()
  if err != nil {
  // continue
  }

  return output.String()
}

// decodes a string from base64 to english
func decodeBase64(textToDecode string) []byte {
  stringDecode, err := base64.StdEncoding.DecodeString(textToDecode)
  if err != nil{
    log.Fatal(err)
  }
  return stringDecode
}

// downloads the file, takes in the path you want to save it to
// as well as the url you want to download it from
func DownloadFile(filepath string, url string) error {

  // waits for thr GET repsonse
  response, err := http.Get(url)
  if err != nil {
      return err
  }
  defer response.Body.Close()

  // Creates the file
  out, err := os.Create(filepath)
  if err != nil {
      return err
  }
  defer out.Close()

  // Writes the file content to the disk
  _, err = io.Copy(out, response.Body)
  return err
}

// decodes env variable from base64 to english and returns output
// of the executed command
func getEnv(value string) string {
  decodeString := decodeBase64(value)
  result := string(decodeString)
  cmdResult := os.Getenv(result)
  return cmdResult
}

// executes commands that only require cmd to execute
func executeCommand(cmd string) string {
  decodeString := decodeBase64(cmd)
  rString := string(decodeString)
  exeCmd := exec.Command("cmd.exe", "/C", rString)
  exeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

  exeCmd.Stdin = strings.NewReader("")
  var output bytes.Buffer
  exeCmd.Stdout = &output

  err := exeCmd.Run()
  if err != nil {
  // continue
  }

  return output.String()
}

// generates random character strings, adds a .extension on the end of each and put them into an array
// used to rename downloaded files with random character strings
func generateRandomString(chars int, numberOfStrings int, imgExtension string) []string{

  randomChar := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  buildStr := ""
  var outputArr []string
  rand.Seed(time.Now().UnixNano())

  loops := 0
  for loops < numberOfStrings{
    for index := 0; index < chars; index ++{
      buildStr += fmt.Sprintf("%c", randomChar[rand.Intn(len(randomChar))])
    }

    buildStr += ".jpeg"
    outputArr = append(outputArr, buildStr)
    buildStr = ""
    loops ++
  }

  return outputArr
}

// creates a 1 second delay
func timeOut(){
  time.Sleep(1*time.Second)
}

func main(){
  
  // raw string literal holding multiple urls
  urls := `
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNy8wNS8xMi8wOC8yOS9jb2ZmZWUtMjMwNjQ3MV85NjBfNzIwLmpwZw==
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNi8wMy8yNi8xMy8wOS9jdXAtb2YtY29mZmVlLTEyODA1MzdfOTYwXzcyMC5qcGc=
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxMy8wOC8xMS8xOS80Ni9jb2ZmZWUtMTcxNjUzXzk2MF83MjAuanBn
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNS8xMC8xMi8xNC81NC9jb2ZmZWUtOTgzOTU1Xzk2MF83MjAuanBn
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNi8wMy8zMC8yMS81OS9jb2ZmZWUtMTI5MTY1Nl85NjBfNzIwLmpwZw==
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxOC8wMi8xNi8xMC81Mi9iZXZlcmFnZS0zMTU3Mzk1Xzk2MF83MjAuanBn
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxMy8xMS8wNS8yMy81NS9jb2ZmZWUtMjA2MTQyXzk2MF83MjAuanBn
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNy8wNC8yNS8wOC8wMi9jb2ZmZWUtYmVhbnMtMjI1ODgzOV85NjBfNzIwLmpwZw==
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNy8wNi8wMi8xMS80OS9zdGlsbC1saWZlLTIzNjYwODRfOTYwXzcyMC5qcGc=
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNS8wNy8xMi8xNC8yNi9jb2ZmZWUtODQyMDIwXzk2MF83MjAuanBn
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNy8wOC8wNy8yMi81Ny9jb2ZmZWUtMjYwODg2NF85NjBfNzIwLmpwZw==
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNi8wMS8wMi8wNC81OS9jb2ZmZWUtMTExNzkzM185NjBfNzIwLmpwZw==
  aHR0cHM6Ly9jZG4ucGl4YWJheS5jb20vcGhvdG8vMjAxNC8xMi8xMS8wMi81Ni9jb2ZmZWUtNTYzNzk3Xzk2MF83MjAuanBn
    `
  urlLinks := strings.Fields(urls)
  // func generates random names with the chosen extnesion
  genRandomNames := generateRandomString(25, len(urlLinks) + 1, ".png")

  // store random names in an array
  storeRandomNames := []string{}
  for index := 0; index < len(genRandomNames); index ++ {
    fmt.Println(genRandomNames[index])
      storeRandomNames = append(storeRandomNames, genRandomNames[index])
  }
  // makes a folder to store files in
  makeFolder := executeCommand("bWQgJWFwcGRhdGElXE1TUGFja2FnZQ==")
  timeOut()
  getPath := getEnv("YXBwZGF0YQ==")
  folderName := decodeBase64("TVNQYWNrYWdl")
  tempPath := ""

  // builds the folder path, and downloads the files
  fmt.Println(makeFolder)
  for index := 0; index < len(urlLinks); index ++ {
    decodeUrl := decodeBase64(urlLinks[index])
    urlString := string(decodeUrl)
    tempPath += getPath + "\\" + string(folderName) + "\\" + storeRandomNames[index]
    writeToDisk := DownloadFile(tempPath, urlString)
    fmt.Println(index, tempPath, "\n", urlString)
    tempPath = ""
    fmt.Println(writeToDisk)
  }
  
  // creates an endless loop and opens files from 0 seconds to 3 minutes
  ex := "ZXhwbG9yZXIuZXhl"
  flag := true
  for flag != false {
      rand.Seed(time.Now().UnixNano())
      genTime := rand.Intn(180000)
      time.Sleep(time.Duration(genTime) * time.Millisecond)
      genValue := rand.Intn(len(storeRandomNames))
      buildStr := getPath + "\\" + string(folderName) + "\\" + storeRandomNames[genValue]
      buildStrEncode := encodeBase64(buildStr)
      runExe(ex, buildStrEncode)
    }
}
