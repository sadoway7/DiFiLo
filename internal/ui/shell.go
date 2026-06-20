package ui

import "fmt"

// acJS wires Google-style autocomplete onto every .df-bigsearch: a debounced
// /suggest fetch fills a dropdown of page-title matches, the first match is
// ghost-completed inline in the input, and arrow/enter/tab/esc keys navigate.
const acJS = `<script>(function(){
function setup(form){
  var input=form.querySelector('input[name="q"]');
  if(!input) return;
  var box=document.createElement('div'); box.className='df-ac'; form.appendChild(box);
  var items=[], sel=-1, last=null, timer=null, ghostOn=false, typing=true;
  function render(list){
    items=list; sel=-1; box.innerHTML='';
    if(!list.length){ box.classList.remove('open'); return; }
    var hd=document.createElement('div'); hd.className='df-ac-head'; hd.textContent='Jump to'; box.appendChild(hd);
    list.forEach(function(it,i){
      var row=document.createElement('div'); row.className='df-ac-item';
      var t=document.createElement('span'); t.className='df-ac-t'; t.textContent=it.t; row.appendChild(t);
      var s=document.createElement('span'); s.className='df-ac-sec'; s.textContent=it.s||''; row.appendChild(s);
      row.addEventListener('mousedown',function(e){e.preventDefault();pick(i);});
      row.addEventListener('mouseenter',function(){setSel(i);});
      box.appendChild(row);
    });
    box.classList.add('open');
  }
  function setSel(i){ sel=i; var rows=box.querySelectorAll('.df-ac-item'); for(var k=0;k<rows.length;k++){ rows[k].classList.toggle('sel',k===i);} }
  function pick(i){ if(i>=0&&i<items.length&&items[i].r) window.location.href=items[i].r; }
  function ghost(first){
    if(!typing||ghostOn) return;
    var val=input.value, trimmed=val.trim();
    if(!first||!first.t||!trimmed) return;
    if(input.selectionStart!==val.length) return;
    var fl=first.t.toLowerCase();
    if(fl===trimmed.toLowerCase()||fl.indexOf(trimmed.toLowerCase())!==0) return;
    input.value=first.t;
    try{ input.setSelectionRange(trimmed.length,first.t.length); }catch(e){}
    ghostOn=true;
  }
  function fetch(){
    var q=input.value; if(q===last) return; last=q;
    var trimmed=q.trim(); if(!trimmed){ render([]); return; }
    var x=new XMLHttpRequest();
    x.open('GET','/suggest?q='+encodeURIComponent(trimmed),true);
    x.onreadystatechange=function(){ if(x.readyState===4&&x.status===200){ try{ var arr=JSON.parse(x.responseText)||[]; var list=arr.map(function(o){return {t:o.Title,r:o.Route,s:o.Section};}); render(list); ghost(list[0]); }catch(e){ render([]);} } };
    x.send();
  }
  input.addEventListener('keydown',function(e){
    var k=e.key;
    if(k==='Backspace'||k==='Delete'){ typing=false; ghostOn=false; }
    else if(k.length===1){ typing=true; }
    if(!items.length) return;
    if(k==='ArrowDown'){ e.preventDefault(); setSel(sel<items.length-1?sel+1:0); }
    else if(k==='ArrowUp'){ e.preventDefault(); setSel(sel>0?sel-1:0); }
    else if(k==='Enter'){ if(sel>=0){ e.preventDefault(); pick(sel);} }
    else if(k==='Tab'){ if(ghostOn){ e.preventDefault(); ghostOn=false; var n=input.value.length; try{input.setSelectionRange(n,n);}catch(e){} } }
    else if(k==='Escape'){ box.classList.remove('open'); }
  });
  input.addEventListener('input',function(){ clearTimeout(timer); ghostOn=false; last=null; timer=setTimeout(fetch,110); });
  input.addEventListener('blur',function(){ setTimeout(function(){ box.classList.remove('open'); },130); });
  input.addEventListener('focus',function(){ if(items.length) box.classList.add('open'); });
}
function init(){ var fs=document.querySelectorAll('.df-bigsearch'); for(var i=0;i<fs.length;i++) setup(fs[i]); }
if(document.readyState==='loading') document.addEventListener('DOMContentLoaded',init); else init();
})();</script>`

// ShellHTML wraps rendered app pages (home/search/list/external) in a full
// document. bodyClass is applied to <body> (used to hide the panel search on
// home/search pages).
func ShellHTML(title, bodyClass, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8">`+
		`<meta name="viewport" content="width=device-width,initial-scale=1">`+
		`<title>%s</title>%s</head><body class="%s">%s%s</body></html>`,
		title, DifiCSS(), bodyClass, body, acJS)
}
