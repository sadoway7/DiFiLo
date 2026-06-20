package ui

import (
	"fmt"
	"html"
)

// CommentsHTML returns the comment section injected before </body> on content pages.
func CommentsHTML(route string, viewer *Viewer) string {
	loggedIn := viewer != nil && viewer.LoggedIn
	canMod := viewer != nil && viewer.LoggedIn && (viewer.Role == "admin" || viewer.Role == "manager")
	routeEsc := html.EscapeString(route)
	userID := "0"
	userRole := ""
	if viewer != nil && viewer.LoggedIn {
		userID = fmt.Sprintf("%d", viewer.ID)
		userRole = viewer.Role
	}
	return fmt.Sprintf(`
<div id="df-comments" data-route="%s" data-loggedin="%t" data-canmod="%t" data-uid="%s" data-role="%s">
<h2>Comments</h2>
<div id="df-comments-list"></div>
<div id="df-comment-form"></div>
</div>
<script>
(function(){
var box=document.getElementById('df-comments');
if(!box)return;
var route=box.dataset.route;
var loggedIn=box.dataset.loggedin==='true';
var canMod=box.dataset.canmod==='true';
var uid=parseInt(box.dataset.uid)||0;
var role=box.dataset.role;

function esc(s){var d=document.createElement('div');d.textContent=s;return d.innerHTML}
function fmtTime(s){try{var d=new Date(s);if(isNaN(d))return s;return d.toLocaleDateString()+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})}catch(e){return s}}

function render(){
  fetch('/api/comments?route='+encodeURIComponent(route)).then(function(r){return r.json()}).then(function(comments){
    var list=document.getElementById('df-comments-list');
    if(!comments||comments.length===0){
      list.innerHTML='<p class="df-muted">No comments yet.</p>';
    } else {
      list.innerHTML='<ul class="df-comment-list">'+
        comments.map(function(c){
          var canDel=(role==='admin'||role==='manager')||(uid===c.UserID);
          var canEdit=(uid===c.UserID);
          var dn=c.Username||'(unknown)';
          return '<li class="df-comment" data-id="'+c.ID+'">'+
            '<div class="df-comment-head">'+
              '<span class="df-comment-author">'+esc(dn)+'</span>'+
              '<span class="df-comment-role">'+esc(c.Role)+'</span>'+
              '<span class="df-comment-date">'+fmtTime(c.CreatedAt)+'</span>'+
              (canEdit?'<button class="df-comment-edit" onclick="dfEditComment('+c.ID+',this)">Edit</button>':'')+
              (canDel?'<button class="df-comment-del" onclick="dfDelComment('+c.ID+')">Delete</button>':'')+
            '</div>'+
            '<div class="df-comment-body" id="df-cb-'+c.ID+'">'+esc(c.Body).replace(/\\n/g,'<br>')+'</div>'+
          '</li>';
        }).join('')+'</ul>';
    }
  });
}

function renderForm(){
  var form=document.getElementById('df-comment-form');
  if(loggedIn){
    form.innerHTML='<textarea id="df-comment-input" placeholder="Write a comment…" rows="3" maxlength="5000"></textarea>'+
      '<button onclick="dfPostComment()">Post Comment</button>';
  } else {
    form.innerHTML='<p class="df-muted"><a href="/login">Log in</a> to post a comment.</p>';
  }
}

window.dfPostComment=function(){
  var input=document.getElementById('df-comment-input');
  var body=input.value.trim();
  if(!body)return;
  fetch('/api/comment',{method:'POST',headers:{'Content-Type':'application/json'},
    body:JSON.stringify({route:route,body:body})})
  .then(function(r){if(!r.ok)throw 0;return r.json()})
  .then(function(){input.value='';render()})
  .catch(function(){alert('Failed to post comment')});
};

window.dfDelComment=function(id){
  if(!confirm('Delete this comment?'))return;
  fetch('/api/comment/delete/'+id,{method:'POST'})
  .then(function(r){if(!r.ok)throw 0;render()})
  .catch(function(){alert('Failed to delete comment')});
};

window.dfEditComment=function(id,btn){
  var bodyDiv=document.getElementById('df-cb-'+id);
  if(!bodyDiv)return;
  if(btn.textContent==='Edit'){
    var text=bodyDiv.textContent;
    bodyDiv.innerHTML='<textarea class="df-edit-input" rows="3" maxlength="5000">'+text.replace(/<br>/g,'\\n')+'</textarea>'+
      '<div class="df-edit-actions"><button class="df-edit-save" onclick="dfSaveEdit('+id+')">Save</button> '+
      '<button class="df-edit-cancel" onclick="dfCancelEdit('+id+')">Cancel</button></div>';
    btn.style.display='none';
  }
};
window.dfSaveEdit=function(id){
  var ta=document.querySelector('#df-cb-'+id+' .df-edit-input');
  if(!ta)return;
  var body=ta.value.trim();
  if(!body)return;
  fetch('/api/comment/edit/'+id,{method:'POST',headers:{'Content-Type':'application/json'},
    body:JSON.stringify({body:body})})
  .then(function(r){if(!r.ok)throw 0;return r.json()})
  .then(function(){render()})
  .catch(function(){alert('Failed to edit comment')});
};
window.dfCancelEdit=function(id){render()};

render();
renderForm();
})();
</script>
`, routeEsc, loggedIn, canMod, userID, userRole)
}
