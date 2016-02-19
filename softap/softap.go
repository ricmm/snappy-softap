package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "encoding/json"
    "io"
    "io/ioutil"
    "bytes"
    "text/template"
    "time"

    "github.com/gorilla/mux"
)

var wlanInterface string

func main() {

    r := mux.NewRouter().StrictSlash(true)

    for _, err := os.Stat(os.Getenv("SNAP_APP_DATA_PATH") + "/interface"); err != nil; _, err = os.Stat(os.Getenv("SNAP_APP_DATA_PATH") + "/interface") {
        fmt.Printf("err: %s, looping\n", err)
        time.Sleep(5*time.Second)
    }

    // make sure there's no one still writing
    time.Sleep(2*time.Second)

    out, _ := ioutil.ReadFile(os.Getenv("SNAP_APP_DATA_PATH") + "/interface")
    //cmd := "iw dev | grep Interface | awk '{print $2}'"
    //out, _ := exec.Command("/bin/bash", "-c", cmd).Output()
    wlanInterface = string(out)
    wlanInterface = strings.TrimSpace(wlanInterface)
    fmt.Printf("iface: '%s'\n", wlanInterface)
    //os.Remove(os.Getenv("SNAP_APP_DATA_PATH") + "/interface")

    // TODO: website should still continue to show even if connected to a network, client reflect

    r.HandleFunc("/connect/", connectWifi).Methods("POST")
    r.HandleFunc("/scan/", scanWifi)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("SNAP_APP_PATH") + "/static/")))

    log.Fatal(http.ListenAndServe(":8888", r))
}

type Response struct {
    Aps []string    `json:"ssids"`
}

type Credentials struct {
    Ssid        string  `json:"wlan_essid"`
    Password    string  `json:"wlan_password"`
    Cookie      string  `json:"cookie"`
}

func scanWifi(w http.ResponseWriter, r *http.Request) {
    cmd := "iwlist " + wlanInterface + " scan | grep ESSID"
    fmt.Printf("cmd: %s", cmd)
    val, _ := exec.Command("/bin/bash", "-c", cmd).Output()
    out := string(val)
    aps := strings.Split(out, "ESSID")[1:]
    for i := range aps {
        aps[i] = strings.TrimSpace(aps[i])
        aps[i] = strings.TrimPrefix(aps[i], ":")
        aps[i] = strings.Trim(aps[i], "\"")
    }

    res := &Response{Aps: aps}
    json, _ := json.Marshal(res)
    fmt.Println(string(json))

    fmt.Fprintln(w, string(json))
}

func connectWifi(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &creds); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    json, _ := json.Marshal(creds)
    fmt.Println(string(json))
    fmt.Fprintln(w, string(json))

    interfaceFile := `allow-hotplug {{.WlanIface}}
iface {{.WlanIface}} inet dhcp
        wpa-ssid {{.WlanEssid}}
        wpa-psk {{.WlanPassword}}`

    var fileOut bytes.Buffer
    t := template.Must(template.New("wrapper").Parse(interfaceFile))
    wrapperData := struct {
        WlanIface       string
        WlanEssid       string
        WlanPassword    string
    }{
        WlanIface:      wlanInterface,
        WlanEssid:      creds.Ssid,
        WlanPassword:   creds.Password,
    }

    if err := t.Execute(&fileOut, wrapperData); err != nil {
        fmt.Printf("Unable to execute template: %v", err)
    }

    if err := ioutil.WriteFile("/etc/network/interfaces.d/" + wlanInterface, fileOut.Bytes(), 0644); err != nil {
        fmt.Printf("Unable to write iface file: %v", err)
    }

    // write cookie file
    if err := ioutil.WriteFile(os.Getenv("SNAP_APP_DATA_PATH") + "/cookie", []byte(creds.Cookie), 0644); err != nil {
        fmt.Printf("Unable to write iface file: %v", err)
    }

    exec.Command("/bin/sync").Output()

    fmt.Printf("%s", fileOut.String())

    rebootCmd := "sleep 5 && /sbin/shutdown -r now"
    reboot := exec.Command("/bin/bash", "-c", rebootCmd)

    // run async
    // here we die
    reboot.Start()
}
