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

func TestDialogsBrowserResponsive(t *testing.T) {
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
 window.fetch=async url=>({ok:true,json:async()=>url==='/api/admin/config'?{host:'127.0.0.1',port:4455,listeningIP:'127.0.0.1',listeningPort:8080,firstSetup:false}:url==='/api/favorites'?{collection:'Demo',scenes:['A','B','C'],slots:['A','B','C']}:{connected:true,scenes:['A','B','C'],favorites:['A','B','C'],currentScene:'A',inputs:[]}});
 window.addEventListener('load',async()=>{try{
  const wait=()=>new Promise(r=>setTimeout(r,150)); await wait();
  const check=(dialog,button)=>{const a=dialog.getBoundingClientRect(); if(!dialog.open||a.left < -1||a.right>innerWidth+1||a.top < -1||a.bottom>innerHeight+1)throw Error('bounds'); button.scrollIntoView({block:'nearest'}); const b=button.getBoundingClientRect(),d=dialog.getBoundingClientRect(); if(b.bottom>d.bottom+1||b.top<d.top-1)throw Error('button inaccessible');};
  document.querySelector('#settings-open').click(); await wait(); check(document.querySelector('#settings-dialog'),document.querySelector('#settings-close')); document.querySelector('#settings-close').click();
  document.querySelector('#favorites-edit').click(); await wait(); check(document.querySelector('#favorites-editor'),document.querySelector('#favorites-close'));
  document.body.dataset.dialogResult='PASS';
 }catch(e){document.body.dataset.dialogResult='FAIL '+e.message}});
 </script>`
		fmt.Fprint(w, strings.Replace(rr.Body.String(), "<head>", "<head>"+fixture, 1))
	}))
	defer server.Close()
	for _, size := range []string{"844,390", "390,844"} {
		out, err := exec.Command(chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--user-data-dir="+t.TempDir(), "--window-size="+size, "--virtual-time-budget=3000", "--dump-dom", server.URL).CombinedOutput()
		if err != nil || !strings.Contains(string(out), `data-dialog-result="PASS"`) {
			t.Errorf("%s: %v\n%s", size, err, out)
		}
	}
}
