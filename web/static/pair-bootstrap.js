(() => {
  const code = location.hash.slice(1);
  if (!code) return;
  history.replaceState(null, '', location.pathname + location.search);
  if (!/^[A-F0-9]{12}$/.test(code)) return;
  const form = document.getElementById('pair-form');
  const input = document.getElementById('code');
  if (form && input) { input.value = code; form.requestSubmit(); }
})();
