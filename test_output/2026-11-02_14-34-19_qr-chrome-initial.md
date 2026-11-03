# QR odaklı tarayıcı denemesi

- Tarih: 2 Kasım 2026 14:34:19 +03, Linux ve headless Chrome.
- Komut: `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./Lib -run '^TestPairingQR' -count=1`.
- Sonuç: çıkış 1; Chrome akışı başlamadan testin HTML'den çıkardığı base64 değerinde HTML karakter referanslarını çözmediği görüldü. Test düzeltilecek.
- OBS kullanılmadı; fiziksel telefon doğrulanmadı. Kullanıcı sunucusuna dokunulmadı.
