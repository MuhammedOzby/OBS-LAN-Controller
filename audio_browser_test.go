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

func TestAudioBrowserMobileLayout(t *testing.T) {
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
window.fetch=async (url,opt)=>{
 if(url.startsWith('/api/status'))return {ok:true,json:async()=>({connected:true,scenes:['A','B','C'],favorites:['A','B','C'],currentScene:'A'})};
 if(url==='/api/audio')return {ok:true,json:async()=>({connected:true,inputs:[{inputName:'Mikrofon',volume:.5,muted:false}]})};
 if(url==='/api/audio/action'){let a=JSON.parse(opt.body);return {ok:true,json:async()=>({inputName:a.inputName,volume:a.volume??.5,muted:a.muted??false})}};
 return {ok:false};
};
window.addEventListener('load',()=>setTimeout(()=>{try{
 const panel=document.querySelector('.audio-panel'), slider=panel.querySelector('input'), button=panel.querySelector('button');
 const rect=panel.getBoundingClientRect();
 if(!slider || slider.disabled || !button || button.disabled)throw Error('audio unavailable');
 if(rect.width<120 || rect.height<30 || rect.bottom>innerHeight+1)throw Error('panel bounds '+JSON.stringify(rect.toJSON()));
 if(document.documentElement.scrollWidth>innerWidth)throw Error('horizontal overflow');
 if(document.querySelectorAll('[data-scene-slot]').length!==3 || document.querySelectorAll('.app-nav .nav-link').length!==4)throw Error('navigation');
 document.body.dataset.audioResult='PASS';
}catch(e){document.body.dataset.audioResult='FAIL '+e.message}},200));
</script>`
		fmt.Fprint(w, strings.Replace(rr.Body.String(), "<head>", "<head>"+fixture, 1))
	}))
	defer server.Close()
	for _, size := range []string{"844,390", "390,844"} {
		out, err := exec.Command(chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--user-data-dir="+t.TempDir(), "--window-size="+size, "--virtual-time-budget=2000", "--dump-dom", server.URL).CombinedOutput()
		if err != nil || !strings.Contains(string(out), `data-audio-result="PASS"`) {
			t.Errorf("%s: %v\n%s", size, err, out)
		}
	}
}
