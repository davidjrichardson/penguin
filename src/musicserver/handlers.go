package musicserver

import (
	"strings"
	"net/http"
	"html/template"
	"encoding/json"
)

// Return homepage
func homeHandler(w http.ResponseWriter, req *http.Request) {
	_, aliasExists := Q.GetAlias(req.RemoteAddr)

	if !aliasExists {
		http.Redirect(w, req, "/alias", http.StatusSeeOther)
	} else {
		plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
		homeTemplate, _ := template.ParseFiles("templates/home.html")
		homeTemplate.Execute(w, plInfo)
	}
}

// Endpoint for setting alias and returns alias set page
func aliasHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		newAlias := req.PostFormValue("alias_value")

		// Check of alias is whitespace
		if len(strings.TrimSpace(newAlias)) == 0 {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
			return
		}

		Q.SetAlias(req.RemoteAddr, newAlias)

		http.Redirect(w, req, "/", http.StatusSeeOther)
	} else {
		aliasTemplate, _ := template.ParseFiles("templates/alias.html")
		aliasTemplate.Execute(w, nil)
	}
}

// Endpoint for queuing videos via link
func queueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			http.Redirect(w, req, "/alias", http.StatusSeeOther)
		}

		videoLink := req.PostFormValue("video_link")

		// If there is space in the playlist or video is whitespace
		if !Q.CanAddVideo(req.RemoteAddr) || len(strings.TrimSpace(videoLink)) == 0{
			vidNotAddedTempl, _ := template.ParseFiles("templates/not_added.html")
			vidNotAddedTempl.Execute(w, nil)
			return
		}

		// Started in new go routine to prevent response waiting
		go Q.DownloadAndAddVideo(req.RemoteAddr, videoLink)

		vidAddedTempl, _ := template.ParseFiles("templates/added.html")
		vidAddedTempl.Execute(w, nil)
	} else {
		// Redirect back to homepage if not a POST request)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}

// Endpoint to return playlist JSON
func playlistHandler(w http.ResponseWriter, req *http.Request) {
	info := Q.GetPlaylistInfo(req.RemoteAddr)
	json.NewEncoder(w).Encode(info)
}