# Ajan talimatı güncellemesi — belge doğrulaması

- Tarih/saat: 2026-10-29 07:31:50 UTC.
- Amaç: Test raporlama kuralını `AGENTS.md` içine eklemek.
- Ortam: Yerel Controller Git deposu; yalnızca belge değişikliği.
- Değişiklik: Her test/test grubu sonucunun `test_output/` içinde Markdown raporuyla saklanması, başarısız ve engellenen denemelerin kaydı, gizli bilgilerin dışlanması ve `TESTING.md` referansı eklendi.

## Doğrulama

- `git -C 'OBS LAN Controller' diff --check`: çıktı yok, çıkış kodu 0; boşluk hatası bulunmadı.
- `git -C 'OBS LAN Controller' status --short`: ` M AGENTS.md`; çıkış kodu 0. Bu rapor oluşturulmadan önceki durumdur.
- Dosya düzenleme aracı değişikliğin uygulandığını doğruladı; Markdown için linter çalıştırılmadı.

## Sınırlamalar

Uygulama veya gerçek OBS/simülasyon testi çalıştırılmadı; bu yalnızca belge doğrulamasıdır. Önceki rapor tarihsel kayıt olarak korunmuştur; o rapordaki engellenen AGENTS.md düzenlemesi bu işlemde kullanıcı onayıyla uygulanmıştır. Yeni commit oluşturulmadı.
