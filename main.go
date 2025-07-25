// carnet_go/main.go
package main

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "os/exec"
    "time"
)

var (
    validUsername = "admin"
    validPassword = "reblochon"
    sessionCookieName = "session_auth"
    sessionDuration = 30 * 60 // 30 minutes en secondes
)

func isAuthenticated(r *http.Request) bool {
    cookie, err := r.Cookie(sessionCookieName)
    return err == nil && cookie.Value == "ok"
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl := template.Must(template.ParseFiles("templates/login.html"))
        tmpl.Execute(w, nil)
        return
    }
    // POST
    username := r.FormValue("username")
    password := r.FormValue("password")
    if username == validUsername && password == validPassword {
        http.SetCookie(w, &http.Cookie{
            Name: sessionCookieName,
            Value: "ok",
            Path: "/",
            HttpOnly: true,
            MaxAge: sessionDuration,
        })
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    tmpl := template.Must(template.ParseFiles("templates/login.html"))
    tmpl.Execute(w, map[string]string{"Error": "Identifiants invalides"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name: sessionCookieName,
        Value: "",
        Path: "/",
        MaxAge: -1,
        HttpOnly: true,
    })
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}

func main() {
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/", requireAuth(formHandler))
    http.HandleFunc("/generer", requireAuth(genererHandler))
    http.HandleFunc("/telecharger", requireAuth(telechargerHandler))
    http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    today := time.Now().Format("2006-01-02")
    start := "2025-01-01"
    tmpl.Execute(w, map[string]string{"Today": today, "Start": start})
}

func genererHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    start := r.FormValue("start_date")
    end := r.FormValue("end_date")
    if end == "" {
        end = time.Now().Format("02/01/2006")
    }

    outDir := "out"
    os.MkdirAll(outDir, 0755)

    layout := "2006-01-02"
    startDate, err := time.Parse(layout, start)
    if err != nil {
        http.Error(w, "Date invalide : "+start, http.StatusBadRequest)
        return
    }
    endDate, err := time.Parse(layout, end)
    if err != nil {
        http.Error(w, "Date invalide : "+start, http.StatusBadRequest)
        return
    }
    formattedStart := startDate.Format("02/01/2006")
    formattedEnd := endDate.Format("02/01/2006")

    cmd := exec.Command("java", "-Dworkdir=out", "-jar", "Carnet_vaccin-1.0-SNAPSHOT.jar", formattedEnd, formattedStart )
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        http.Error(w, "Erreur lors de l'exécution du programme Java", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/telecharger", http.StatusSeeOther)
}

func telechargerHandler(w http.ResponseWriter, r *http.Request) {
    filePath := "out/out.xlsx"
    f, err := os.Open(filePath)
    if err != nil {
        http.Error(w, "Fichier non trouvé", http.StatusNotFound)
        return
    }
    defer f.Close()

    w.Header().Set("Content-Disposition", "attachment; filename=out.xlsx")
    w.Header().Set("Content-Type", "officedocument.spreadsheetml.sheet")
    io.Copy(w, f)
}
