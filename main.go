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

func main() {
    http.HandleFunc("/", formHandler)
    http.HandleFunc("/generer", genererHandler)
    http.HandleFunc("/telecharger", telechargerHandler)
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
