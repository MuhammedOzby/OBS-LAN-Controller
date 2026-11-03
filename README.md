# OBS LAN Controller

> **AI-assisted development:** This project was developed with the help of AI tools.

**English** | [Türkçe](README-TR.md)

A Go-powered OBS dashboard using the Bootstrap Kehribar Gece theme. The **Control (Kontrol)** tab provides program preview, OBS metrics, three favorite scenes, and audio level/mute controls. **Deck** offers 12 customizable buttons for scenes, source visibility, mute, and media controls. **Broadcast Management (Yayın Yönetimi)** is the only place for streaming and recording actions. The fourth tab, **Mobile sign-in (Mobil cihaz ile giriş)**, pairs phones using QR codes. Trusted HTTPS and PWA support are not available yet.

The current interface uses Turkish labels; English translations below are accompanied by the original labels where needed.

## Download

```sh
git clone https://github.com/MuhammedOzby/OBS-LAN-Controller.git
cd OBS-LAN-Controller
```

For live preview, also install [OBS LAN Controller Live Preview Extension](https://github.com/MuhammedOzby/OBS-LAN-Controller-Live-Preview-Extension). Scene, audio, streaming, and recording controls use OBS's built-in WebSocket server.

> **Usage status — November 3, 2026:** The project owner reports having tested and actively using the application and plugin. This report is separate from the dated test reports below. The tested OBS version, device/browser, and numerical performance measurements still need to be documented.

## Run

If you are new to OBS, start with the [basic OBS setup and recording troubleshooting guide](RECOMMENDATIONS.md) (Turkish).

For testing with real OBS and a phone, see the [step-by-step manual checklist](MANUAL_TEST_CHECKLIST.md) (Turkish).

Requires Go 1.24 or later. Open OBS, then run this command in the project directory:

```sh
go run .
```

1. On the OBS computer, open **http://127.0.0.1:8080**.
2. In the first-run **Connection settings (Bağlantı ayarları)** dialog, enter and save OBS's WebSocket address, port, and password. For OBS on the same computer, use **127.0.0.1** and the default port **4455**.
3. To allow phone access, set **HTTP listening IP address (HTTP dinleme IP adresi)** to **0.0.0.0** in the same dialog. Save, stop the process with `Ctrl+C`, and run `go run .` again. Other devices on your trusted network can then connect.
4. On your phone, open **http://SERVER_LAN_IP:8080** and follow the pairing instructions below. Find the computer's LAN IP with `hostname -I`. **Do not enter 0.0.0.0 in the browser.**

Without a configuration file, first-run setup listens only on this computer. An existing file supplies the saved address and port. The **Settings (Ayarlar)** button is available only from the OBS computer through loopback; phones cannot change the OBS connection or password settings. If you choose a different HTTP port, replace 8080 in these URLs.

You can also run `go run . serve`. Settings are read from `configuration.json` in the working directory by default. When running from another directory, use `-config /absolute/path/configuration.json`.

### Configuration

- `obs.ip`: OBS computer's IP address; defaults to `127.0.0.1`.
- `obs.port`: OBS WebSocket port; defaults to `4455`.
- `obs.password`: OBS WebSocket password; enter your own password in the local file. An empty value is accepted for OBS without password authentication.
- `server.listening_ip`: HTTP listening IP; `0.0.0.0` listens on all IPv4 interfaces, while `127.0.0.1` allows access only from this computer.
- `server.port`: HTTP port; defaults to `8080`.

Normal setup does not require editing JSON manually. You may copy `configuration.example.json` to `configuration.json` and edit it. Replace the `OBS_WEBSOCKET_PAROLANIZ` placeholder with your OBS WebSocket password, or use an empty string for OBS without a password.

IP addresses and ports in the range 1–65535 are validated. An invalid existing file stops the application and is not overwritten. A missing file opens local first-run setup. Restart after changing the HTTP listening address or port. OBS connection settings saved in the dialog take effect when reconnecting on the next request. `-addr` overrides the HTTP address temporarily; only a loopback address is accepted during first-run setup.

`configuration.json` contains a password and is excluded from Git. The shareable template is `configuration.example.json`. The panel writes the file atomically with owner-only `0600` permissions. If you create it manually, protect it with `chmod 600 configuration.json`. The saved password is never sent back to the browser. Leaving the new password field empty preserves the saved password; selecting **Delete saved password (Kayıtlı parolayı sil)** clears it.

OBS's built-in WebSocket server must be enabled. On disconnect, stale data is cleared and subsequent requests retry after a two-second delay. Commands are not replayed automatically.

### Scenes and audio

- In the Control tab, use **Edit favorites (Favorileri düzenle)** to select and save three scene slots. Slots may be empty; the same scene cannot occupy two slots. Before the first save, the first three OBS scenes are used.
- Preferences are stored server-side in `configuration.json.favorites.json` beside the configuration file and shared by paired devices. Deleted scenes or a different scene collection are not silently substituted; edit favorites again if needed.
- The audio section lists OBS sources that support audio. The slider covers 0–100%; the button mutes or unmutes. Changes made in OBS are reflected by periodic polling. Controls are disabled when disconnected, and submitted changes are verified by reading OBS again.

Preview receives H.264 program video from the companion OBS plugin over a Unix socket and forwards it to the browser using **WebRTC without re-encoding**. The panel does not use a screenshot loop. Video and OBS control connections are tracked separately; losing video does not disable scene, Deck, or broadcast controls.

Bitrate is calculated from the change in bytes sent over time. No value is shown for the first sample, a counter reset, or an inactive/reconnecting stream. The quality label is a technical indicator based on measured dropped frames/network congestion, not a measurement of viewer experience.

### Run alongside the OBS plugin on the same computer (Linux)

1. Install **[OBS LAN Controller Live Preview Extension](https://github.com/MuhammedOzby/OBS-LAN-Controller-Live-Preview-Extension)** into native Linux OBS and restart OBS. Flatpak has not been validated for this release.
2. Run the Controller **as the same user who runs OBS**, from a terminal in the desktop session. Do not use `sudo`. `XDG_RUNTIME_DIR` must come from that session, usually `/run/user/UID`; do not point it to an arbitrary shared or temporary directory.
3. On Arch/CachyOS, `sudo pacman -Syu --needed go` provides the Controller build requirement. The Controller does not require the OBS SDK, CMake, Node.js, or FFmpeg to run. In this directory:

   ```sh
   # Only if you do not already have a configuration file:
   cp -n configuration.example.json configuration.json
   chmod 600 configuration.json
   # Set obs.ip=127.0.0.1, the port, and your OBS WebSocket password in the file.
   go build -o obs-lan-controller .
   ./obs-lan-controller serve
   ```

4. Open the **Control (Kontrol)** tab at `http://127.0.0.1:8080`. Video starts when the tab is visible and WebRTC is connected. It stops when the last viewer leaves or the page becomes hidden, and reconnects when visible again. The plugin reconnects automatically, so you may start the Controller after OBS.
5. For phone access, pair through the existing `/pair` flow and open the panel using the LAN IP. A modern browser with H.264/WebRTC support is required. If autoplay is blocked, tap the video. The application does not use internet STUN/TURN servers; it uses direct local/LAN ICE candidates only. In addition to HTTP TCP 8080, the firewall must allow WebRTC's dynamic UDP traffic. The application does not change firewall rules. Do not forward ports from the internet.

The video socket is a filesystem path, not an IP/port: `$XDG_RUNTIME_DIR/obs-lan-controller/preview.sock`. The Controller is the server and the plugin is the client. They use a 0700 directory, a 0600 socket, and same-UID peer checks. A second Controller does not take over the socket. If socket startup fails, HTTP/OBS control remains available and the reason is logged.

For diagnosis, open `/api/preview/status` from the local browser. `connected` indicates a plugin connection; `streaming` indicates receipt of H.264 video during the last three seconds. Neither measures browser-decoded frames or latency. The API's `fps` is the encoder configuration. The panel's preview FPS is measured from browser WebRTC receiver statistics and displays `—` when no data is available.

The initial version allows at most **3 concurrent preview sessions**. This is a resource limit, not a measured capacity guarantee. Unconnected sessions expire after 15 seconds; all sessions close after at most one hour, with visible panels reconnecting. Revoking or expiring LAN pairing also closes the video session. For a slow viewer, the bounded queue is cleared and waits for the next keyframe. Audio is not transmitted. The existing plugin protocol is unchanged, so this Controller update does not require reinstalling the plugin.

Automated validation used real H.264 FFmpeg test video → simulated plugin → real Unix socket → WebRTC/SRTP/RTP → H.264 decoding. This does not replace testing real OBS program video, phones, scene transitions, or performance on the target computer. FFmpeg is an optional integration-test tool, not an application dependency.

Both the shareable example and first-run setup without a file listen on loopback only (`127.0.0.1:8080`). For phone access, change the HTTP listening address to `0.0.0.0` in local settings and restart the Controller. Use only a trusted local network and do not forward ports from the internet. **HTTP is not encrypted; someone able to read network traffic could capture the pairing code or session.** Trusted HTTPS remains an open production requirement. If the firewall blocks TCP 8080, access must be allowed separately; the application does not change system or firewall settings.

### Phone/LAN pairing

1. On the OBS computer, open `http://127.0.0.1:8080` and select **Mobile sign-in (Mobil cihaz ile giriş)**. You can also open `http://127.0.0.1:8080/pair` directly. Use your configured HTTP port if different.
2. With the phone and computer on the same trusted network, scan the QR code using the phone's normal camera and open the link. The phone pairs and opens Control without typing a code. Browser camera permission is not required.
3. If several QR codes are shown, choose the one under the phone's Wi-Fi/LAN address. If no LAN address is available, set **Settings → HTTP listening IP address (Ayarlar → HTTP dinleme IP adresi)** to `0.0.0.0`, save, and restart only the Controller. Do not disable the entire firewall.
4. If the camera cannot scan the code, open `http://SERVER_LAN_IP:8080` on the phone and enter the 12-character code shown on the same page.

QR codes and pairing codes are single-use and expire after five minutes. Generate a new code for each additional phone. Refreshing `/pair` on the computer invalidates the previous code; the fourth tab shares the current valid code. Generate a new QR code after expiry or use. QR codes are shown only through loopback on the computer. Paired phones cannot view or generate a new pairing code/QR.

Sessions last 12 hours. **Revoke all LAN sessions (Tüm LAN oturumlarını iptal et)** or restarting the server removes access. Signing in again after a restart is expected; scan a new QR code. QR pairing does not provide permanent sign-in or encryption. The OBS password is not included in the QR code. QR codes are generated locally without an external QR service.

Unpaired LAN API requests receive 401; panel pages redirect to pairing. Session cookies use HttpOnly/SameSite=Strict and Secure when served directly over TLS. Limits are 128 sessions and 20 failed code attempts per minute. The OBS password is not the pairing code and is never sent to the browser.

Host protection accepts only localhost/loopback and interface IPs detected at startup. Arbitrary DNS names and proxy headers are not trusted. Restart if the IP changes. Reverse proxies are unsupported; a loopback proxy could bypass the local-access boundary. All request bodies have a five-second read timeout.

## Updating from the initial GitHub upload

The initially uploaded `configuration.json`, `configuration.json.favorites.json`, and compiled `obs-lan-controller` are no longer distributed through Git. Use `configuration.example.json` as the template and build the application from source. New favorite selections remain local and are excluded from Git.

If updating a working copy cloned from the initial upload, **before running `git pull`**, back up your configuration and favorites outside the repository in a private directory. The update removing these files from tracking may remove them from your working directory. Restore them after updating and protect `configuration.json` with `chmod 600 configuration.json`. Rebuild with `go build -o obs-lan-controller .`. Restart the running Controller at a suitable time.

## Development

The layout prioritizes 19:9 landscape screens and fits the available browser area (`100dvh`) without horizontal or vertical page scrolling. The header, tabs, and footer retain their positions while content uses the remaining space. Portrait screens are also supported. Check for overflow when adding components.

- `main.go`: application entry point and error exit only.
- `Lib/configuration.go`: `Configuration` and `ReadConfiguration`; JSON loading and IP/port validation.
- `Lib/settings.go`: loopback-only connection settings, password preservation, and atomic saving.
- `Lib/favorites.go`: persistent selection of three favorite scenes and related APIs.
- `Lib/audio.go`: audio sources and verified level/mute actions.
- `Lib/obs.go`: authenticated OBS WebSocket client, metrics, temporary screenshots, and scene switching/verification.
- `Lib/serve.go`: `Run` and `NewHandler`; Gin routes, templates, static serving, and clean shutdown through `http.Server`.
- `web/assets.go`: embeds templates and static files with `go:embed`, independent of the working directory.
- `web/templates/*.tmpl`: shared layout and page areas.
- `web/static/app.css`: layout and Kehribar Gece colors.
- `web/static/app.js`: periodic status refresh and scene controls.
- `web/static/audio.js`, `settings.js`: audio actions and the local settings dialog.
- `web/static/vendor/`: local theme CSS and its license.

Assets are compiled using `go:embed`; restart the server after template/CSS changes. Node.js, a CDN, and internet access are not needed at runtime. PWA support is not available yet.

```sh
go test ./...
go vet ./...
go build -o obs-lan-controller .
```

Test results are stored under `test_output/`. On real demo OBS, Deck scene switching, source show/hide, mute/unmute, and media play/pause/restart/stop were verified. The test creates and removes a temporary scene and local media input, restores the previous program scene, and verifies cleanup by reading OBS again. Streaming/recording state is not changed. Physical phone/tablet and long-running load tests were not performed in those reports.

The November 2, 2026 Terra delivery passed Go/Node/real headless Chrome tests, race checks, vet, and builds. The new build read real OBS status, audio sources, and favorite options without changing them. That delivery's acceptance of new scene/audio commands, physical devices, and recording remained open; an earlier quick recording test was observed as still active at the time. **Overall real-OBS acceptance for Milestone 1 was not granted.** See `PLAN.md` and the dated reports for context and subsequent observations.

Run the live Deck acceptance test only with explicit permission:
`OBS_DECK_LIVE=1 go test ./Lib -run '^TestDeckLiveAcceptance$' -v -count=1`.
First place a local video with audio lasting at least a few seconds at `test_output/deck-acceptance.mp4`. Normal tests do not connect to OBS. Visibility buttons are bound to collection and source identities; old buttons without identities must be reassigned in the editor. The layout is stored only in this browser's localStorage.

## Theme source

[Bootstrap Kehribar Gece](https://github.com/MuhammedOzby/bootstrap-kehribar-gece), commit `0a7d4423c00e3dab021ec70570224f7d41466f4f`, `dist/css/bootstrap.min.css`. The upstream file is included unchanged. License: `web/static/vendor/LICENSE.bootstrap`. Application-specific layout and amber/night colors are layered in `app.css`.
