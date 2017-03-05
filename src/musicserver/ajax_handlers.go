package musicserver

import (
	"strings"
	"html/template"
	"encoding/json"
	"net/http"
)

// Endpoint for queueing video via ajax
func ajaxQueueHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	out := map[string]string{}

	if req.Method == http.MethodPost {
		_, aliasExists := Q.GetAlias(req.RemoteAddr)

		if !aliasExists {
			out["Message"] = "No user alias set"
			out["Type"] = "error"
			json.NewEncoder(w).Encode(out)
			return
		}

		videoLink := req.PostFormValue("video_link")

		// Submitted video link is blank
		if len(strings.TrimSpace(videoLink)) == 0 {
			out["Message"] = "No video link given"
			out["Type"] = "error"
			json.NewEncoder(w).Encode(out)
			return
		}

		// If user has max added videos
		if !Q.CanAddVideo(req.RemoteAddr) {
			out["Message"] = "Video not added, user has too many videos"
			out["Type"] = "warn"
			json.NewEncoder(w).Encode(out)
			return
		}

		// Start video downloader in new goroutine so 
		Q.QuickAddVideo(req.RemoteAddr, videoLink)

		out["Message"] = "Video added"
		out["Type"] = "success"
		json.NewEncoder(w).Encode(out)
	} else {
		out["Message"] = "Use POST method"
		out["Type"] = "error"
		json.NewEncoder(w).Encode(out)
	}
}

func ajaxPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
	templ, _ := template.ParseFiles("templates/playlist.html")
	templ.Execute(w, plInfo)
}

func ajaxAdminPlaylistHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	plInfo := Q.GetPlaylistInfo(req.RemoteAddr)
	templ, _ := template.ParseFiles("templates/admin_playlist.html")
	templ.Execute(w, plInfo)
}