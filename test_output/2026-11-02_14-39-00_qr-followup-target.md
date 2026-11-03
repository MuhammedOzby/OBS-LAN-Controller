# QR inceleme düzeltmeleri: hedefli test

- Tarih: 2 Kasım 2026 14:39:00 +03, Linux, geçici Node ve headless Chrome.
- Komut: `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./Lib -run 'TestPairing(CrossSiteScannerLanding|RealMobileLayoutChrome)' -count=1`.
- Sonuç: çıkış 0. Yerel çapraz-site kod üretimi ve çapraz-site POST reddi, uzak GET'in durumsuz açılması; gerçek panel ve /pair için üç viewportta kaydırma ve taşma denetimleri geçti. Ek adres yerleşimi testte klonlanmış QR ile simüle edildi.
- OBS kullanılmadı; gerçek OBS ve fiziksel telefon doğrulanmadı. Çalışan kullanıcı sunucusuna dokunulmadı.
