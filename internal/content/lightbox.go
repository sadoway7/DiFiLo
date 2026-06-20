package content

// LightboxHTML is the fullscreen image modal plus the JavaScript that drives
// it (open on image click, close on Esc / backdrop / close button). It also
// wires up inline content images so they open the lightbox on click.
const LightboxHTML = `
<div id="df-lightbox" onclick="dfLightboxClose()">
  <div id="df-lightbox-close" onclick="event.stopPropagation();dfLightboxClose()">&times;</div>
  <img id="df-lightbox-img" src="" alt="" onclick="event.stopPropagation()">
  <div id="df-lightbox-cap"></div>
</div>
<script>
function dfLightbox(img){
  var lb=document.getElementById('df-lightbox');
  var lbImg=document.getElementById('df-lightbox-img');
  var lbCap=document.getElementById('df-lightbox-cap');
  lbImg.src=img.src;
  lbCap.textContent=img.dataset.caption||'';
  lb.classList.add('open');
  document.body.style.overflow='hidden';
}
function dfLightboxClose(){
  document.getElementById('df-lightbox').classList.remove('open');
  document.body.style.overflow='';
}
document.addEventListener('keydown',function(e){
  if(e.key==='Escape')dfLightboxClose();
});
// Make inline content images clickable — prevent parent link navigation
document.querySelectorAll('.df-wiki-content img').forEach(function(img){
  if(!img.classList.contains('df-wiki-gal')){
    img.style.cursor='pointer';
    img.addEventListener('click',function(e){
      e.preventDefault();
      e.stopPropagation();
      dfLightbox(this);
    });
    if(img.alt)img.dataset.caption=img.alt;
  }
});
</script>
`
