(() => {
  const dialog = document.getElementById('settings-dialog');
  const field = name => document.getElementById(`settings-${name}`);
  let loading = false, saving = false, openRequested = false;
  async function load(open) {
    if (open) openRequested = true;
    if (loading || saving) return;
    loading = true;
    try {
      const response = await fetch('/api/admin/config', {cache: 'no-store'});
      if (!response.ok) throw new Error();
      const data = await response.json();
      field('host').value = data.host;
      field('port').value = data.port;
      field('listen').value = data.listeningIP;
      field('listen-port').value = data.listeningPort;
      field('password').value = '';
      field('clear').checked = false;
      field('message').textContent = data.passwordConfigured ? 'Kayıtlı parola var.' : 'Kayıtlı parola yok.';
      if ((openRequested || data.firstSetup) && !dialog.open) dialog.showModal();
    } catch (_) { field('message').textContent = 'Ayarlar okunamadı.'; if (openRequested && !dialog.open) dialog.showModal(); }
    finally { loading = false; openRequested = false; }
  }
  field('open').addEventListener('click', () => load(true));
  field('close').addEventListener('click', () => dialog.close());
  field('form').addEventListener('submit', async event => {
    event.preventDefault();
    if (loading || saving) return;
    const password = field('password').value;
    if (password && field('clear').checked) { field('message').textContent = 'Parolayı değiştirmek veya silmek için yalnız birini seçin.'; return; }
    const body = {host: field('host').value, port: Number(field('port').value), listeningIP: field('listen').value, listeningPort: Number(field('listen-port').value)};
    if (field('clear').checked) body.password = '';
    else if (password) body.password = password;
    saving = true;
    const controls = [...field('form').querySelectorAll('input, button')];
    controls.forEach(control => { control.disabled = true; });
    try {
      const response = await fetch('/api/admin/config', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)});
      if (!response.ok) throw new Error();
      field('password').value = '';
      field('clear').checked = false;
      field('message').textContent = 'Kaydedildi. Dinleme adresi değişikliği için sunucuyu yeniden başlatın.';
    } catch (_) { field('message').textContent = 'Ayarlar kaydedilemedi; değerleri kontrol edin.'; }
    finally { saving = false; controls.forEach(control => { control.disabled = false; }); }
  });
  load(false);
})();
