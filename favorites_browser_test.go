package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	lib "obs-lan-controller/Lib"
)

func TestFavoritesBrowserEditor(t *testing.T) {
	chrome := os.Getenv("CHROME_BIN")
	if chrome == "" {
		t.Skip("CHROME_BIN required")
	}
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			h.ServeHTTP(w, r)
			return
		}
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, r)
		fixture := `<script>
let saved=null; const scenes=['A','B','C','D'];
window.fetch=async (url,opt)=> {
 if(url.startsWith('/api/status')) return {ok:true,json:async()=>({connected:true,scenes,favorites:saved||['A','B','C'],currentScene:'A'})};
 if(url==='/api/favorites' && opt?.method==='POST'){saved=JSON.parse(opt.body).slots;return {ok:true,status:204}};
 if(url==='/api/favorites') return {ok:true,json:async()=>({collection:'Demo',scenes,slots:saved||['A','B','C']})};
 return {ok:false};
};
window.addEventListener('load', async()=>{try {
 const wait=()=>new Promise(r=>setTimeout(r,100)); await wait();
 document.querySelector('#favorites-edit').click(); await wait();
 if(!document.querySelector('#favorites-editor').open) throw Error('dialog');
 document.querySelector('#favorite-0').value='D'; document.querySelector('#favorite-1').value='';
 document.querySelector('#favorites-form').requestSubmit(); await wait();
 const slots=[...document.querySelectorAll('[data-scene-slot]')];
 if(JSON.stringify(saved)!=='["D","","C"]'||slots[0].textContent!=='D'||!slots[1].disabled) throw Error('save/render');
 document.body.dataset.favoriteResult='PASS';
 }catch(e){document.body.dataset.favoriteResult='FAIL '+e.message}});
</script>`
		fmt.Fprint(w, strings.Replace(rr.Body.String(), "<head>", "<head>"+fixture, 1))
	}))
	defer server.Close()
	out, err := exec.Command(chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--user-data-dir="+t.TempDir(), "--virtual-time-budget=2000", "--dump-dom", server.URL).CombinedOutput()
	if err != nil || !strings.Contains(string(out), `data-favorite-result="PASS"`) {
		t.Fatalf("browser editor: %v\n%s", err, out)
	}
}
