package main

import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "bytes"
    "strings"
    "github.com/gorilla/mux"
)
func shell(w http.ResponseWriter, r *http.Request) {
vars := mux.Vars(r)
var out bytes.Buffer
workDir := "/Users/ichtar/git/helm-charts"
// force endpoint to accept only specific values send a forbidden otherwise
if strings.EqualFold(vars["cluster"],"eu1.bestmile.com") && strings.EqualFold(vars["deployment"],"test") {
    // get last version of code
    cmd := exec.Command("git","pull","--rebase")
    cmd.Dir = workDir
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    envPath := strings.Join([]string{"environments/",vars["cluster"],"/",vars["deployment"],".yaml"},"")
    // remove dry-run when want to push real update
    cmd = exec.Command("helm","--dry-run","--debug","--kube-context",vars["cluster"],"upgrade",vars["deployment"],"-f",envPath,"charts/bm-stack")
    cmd.Dir = workDir
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil {
	log.Fatal(err)
    }
        fmt.Fprintf(w, out.String())
    } else {
        http.Error(w, "Forbidden", http.StatusForbidden)
    }
}
func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/trigger/{cluster}/{deployment}", shell ).Methods("PUT")
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func main() {
    handleRequests()
}
