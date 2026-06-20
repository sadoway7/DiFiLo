package app

import (
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"difilo/internal/auth"
)

// handleAdmin renders the admin user-management panel, download settings, and
// recent download log. Admin-only.
func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil || user.Role != auth.RoleAdmin {
		s.renderShell(w, r, "Access Denied", "",
			"<div class='df-external'><h1>Admin only</h1><p>You must be logged in as an admin to view this page.</p></div>")
		return
	}
	users, _ := s.DB.AllUsers()

	var b strings.Builder
	b.WriteString("<div class='df-admin'>")
	b.WriteString("<h1>User Management</h1>")
	b.WriteString(fmt.Sprintf("<p class='df-muted'>%d users</p>", len(users)))
	b.WriteString("<table class='table'><thead><tr><th>Username</th><th>Email</th><th>Role</th><th>Created</th><th>Actions</th></tr></thead><tbody>")
	for _, u := range users {
		b.WriteString("<tr>")
		uname := u.Username
		if uname == "" {
			uname = u.Email
		}
		b.WriteString(fmt.Sprintf("<td><b>%s</b></td>", html.EscapeString(uname)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(u.Email)))
		b.WriteString(fmt.Sprintf("<td><select class='df-role-select' data-uid='%d'>", u.ID))
		for _, role := range []string{auth.RoleGeneral, auth.RoleManager, auth.RoleAdmin} {
			selected := ""
			if u.Role == role {
				selected = " selected"
			}
			b.WriteString(fmt.Sprintf("<option value='%s'%s>%s</option>", role, selected, role))
		}
		b.WriteString("</select></td>")
		created := u.CreatedAt.Format("2006-01-02")
		if u.CreatedAt.IsZero() {
			created = "—"
		}
		b.WriteString(fmt.Sprintf("<td>%s</td>", created))
		b.WriteString(fmt.Sprintf("<td><a class='df-user-dl' href='/admin/user-downloads?uid=%d'>Downloads</a> ", u.ID))
		if u.ID != user.ID {
			b.WriteString(fmt.Sprintf("<button class='df-del-user' data-uid='%d'>Delete</button></td>", u.ID))
		} else {
			b.WriteString("<span class='df-muted'>(you)</span></td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	b.WriteString("</div>")

	// Download settings section.
	downloadURL := s.DB.GetSetting("download_url")
	b.WriteString("<div class='df-admin'>")
	b.WriteString("<h1>Offline Download Link</h1>")
	b.WriteString("<p class='df-muted'>Set the external download link for the offline archive. This URL will be used by the download button on the home page.</p>")
	b.WriteString(fmt.Sprintf("<div class='df-setting-row'><label>Download URL</label>"+
		"<input type='text' id='df-download-url' value='%s' placeholder='https://example.com/DiFilo-june12026.7z' style='flex:1'>"+
		"<button onclick='dfSaveDownload()'>Save</button></div>", html.EscapeString(downloadURL)))
	// Download password setting.
	downloadPW := s.DB.GetSetting("download_password")
	b.WriteString(fmt.Sprintf("<div class='df-setting-row'><label>Download password</label>"+
		"<input type='text' id='df-download-pw' value='%s' placeholder='(unset)' style='flex:1'>"+
		"<button onclick='dfSaveDownloadPw()'>Save</button></div>", html.EscapeString(downloadPW)))
	b.WriteString("</div>")

	// Page download log.
	dls, _ := s.DB.RecentDownloads(50)
	b.WriteString("<div class='df-admin'>")
	b.WriteString("<h1>Page Downloads</h1>")
	b.WriteString(fmt.Sprintf("<p class='df-muted'>%d recent downloads (newest first)</p>", len(dls)))
	b.WriteString("<table class='table'><thead><tr><th>User</th><th>Page</th><th>Title</th><th>When</th></tr></thead><tbody>")
	for _, d := range dls {
		uname := d.Username
		if uname == "" {
			uname = fmt.Sprintf("#%d", d.UserID)
		}
		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf("<td><b>%s</b></td>", html.EscapeString(uname)))
		b.WriteString(fmt.Sprintf("<td><a href='%s'>%s</a></td>", html.EscapeString(d.Route), html.EscapeString(d.Route)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(d.Title)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(d.CreatedAt)))
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	b.WriteString("</div>")

	// JS for role change + delete + settings.
	b.WriteString(`<script>
	document.querySelectorAll('.df-role-select').forEach(function(sel){
		sel.addEventListener('change',function(){
			fetch('/api/admin/role',{method:'POST',headers:{'Content-Type':'application/json'},
				body:JSON.stringify({uid:parseInt(this.dataset.uid),role:this.value})})
			.then(function(r){if(!r.ok)alert('Failed to update role')});
		});
	});
	document.querySelectorAll('.df-del-user').forEach(function(btn){
		btn.addEventListener('click',function(){
			if(!confirm('Delete this user and all their comments?'))return;
			fetch('/api/admin/delete-user',{method:'POST',headers:{'Content-Type':'application/json'},
				body:JSON.stringify({uid:parseInt(this.dataset.uid)})})
			.then(function(r){if(r.ok)location.reload();else alert('Failed to delete user')});
		});
	});
	window.dfSaveDownload=function(){
		var url=document.getElementById('df-download-url').value.trim();
		fetch('/api/admin/settings',{method:'POST',headers:{'Content-Type':'application/json'},
			body:JSON.stringify({key:'download_url',value:url})})
		.then(function(r){if(r.ok)alert('Saved');else alert('Failed to save')});
	};
	window.dfSaveDownloadPw=function(){
		var pw=document.getElementById('df-download-pw').value;
		fetch('/api/admin/settings',{method:'POST',headers:{'Content-Type':'application/json'},
			body:JSON.stringify({key:'download_password',value:pw})})
		.then(function(r){if(r.ok)alert('Saved');else alert('Failed to save')});
	};
	</script>`)

	s.renderShell(w, r, "Admin", "", b.String())
}

// handleAdminUserDownloads shows all page downloads for a single user
// (admin only).
func (s *Server) handleAdminUserDownloads(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil || user.Role != auth.RoleAdmin {
		s.renderShell(w, r, "Access Denied", "",
			"<div class='df-external'><h1>Admin only</h1><p>You must be logged in as an admin to view this page.</p></div>")
		return
	}
	uid, err := strconv.ParseInt(r.URL.Query().Get("uid"), 10, 64)
	if err != nil || uid < 1 {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}
	target, _ := s.DB.GetUserByID(uid)
	if target == nil {
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}
	dls, _ := s.DB.DownloadsByUser(uid)

	var b strings.Builder
	b.WriteString("<div class='df-admin'>")
	b.WriteString("<p class='df-muted'><a href='/admin'>&#8249; Back to User Management</a></p>")
	name := target.Username
	if name == "" {
		name = target.Email
	}
	b.WriteString(fmt.Sprintf("<h1>Downloads by %s</h1>", html.EscapeString(name)))
	b.WriteString(fmt.Sprintf("<p class='df-muted'>%d page downloads (newest first)</p>", len(dls)))
	b.WriteString("<table class='table'><thead><tr><th>Page</th><th>Title</th><th>When</th></tr></thead><tbody>")
	for _, d := range dls {
		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf("<td><a href='%s'>%s</a></td>", html.EscapeString(d.Route), html.EscapeString(d.Route)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(d.Title)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(d.CreatedAt)))
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	b.WriteString("</div>")
	s.renderShell(w, r, "User Downloads", "", b.String())
}

// handleAPIAdminRole updates a user's role (admin only).
func (s *Server) handleAPIAdminRole(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil || user.Role != auth.RoleAdmin {
		jsonError(w, "admin only", http.StatusForbidden)
		return
	}
	var body struct {
		UID  int64  `json:"uid"`
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Role != auth.RoleGeneral && body.Role != auth.RoleManager && body.Role != auth.RoleAdmin {
		jsonError(w, "invalid role", http.StatusBadRequest)
		return
	}
	if err := s.DB.UpdateUserRole(body.UID, body.Role); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleAPIAdminDeleteUser deletes a user and their content (admin only). An
// admin cannot delete themselves.
func (s *Server) handleAPIAdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil || user.Role != auth.RoleAdmin {
		jsonError(w, "admin only", http.StatusForbidden)
		return
	}
	var body struct {
		UID int64 `json:"uid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.UID == user.ID {
		jsonError(w, "cannot delete yourself", http.StatusBadRequest)
		return
	}
	if err := s.DB.DeleteUser(body.UID); err != nil {
		slog.Error("delete user failed", "uid", body.UID, "error", err)
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleAPIAdminSettings upserts a key/value pair in the settings table
// (admin only). Known keys include download_url and download_password.
func (s *Server) handleAPIAdminSettings(w http.ResponseWriter, r *http.Request) {
	user := s.currentUser(r)
	if user == nil || user.Role != auth.RoleAdmin {
		jsonError(w, "admin only", http.StatusForbidden)
		return
	}
	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Key == "" {
		jsonError(w, "key required", http.StatusBadRequest)
		return
	}
	if err := s.DB.SetSetting(body.Key, body.Value); err != nil {
		jsonError(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
