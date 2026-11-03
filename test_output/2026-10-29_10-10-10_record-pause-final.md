# Record pause final regression
UTC 2026-10-29T10:10:10.755671+00:00
Mock/Node tests and vet. Real OBS trials documented separately: shared Stream encoder cannot pause; alternative HQ encoders did not start. Original Stream/nvenc values restored and recording stopped, readback verified. No working real pause/resume claimed.

go test ./... -count=1 -timeout=60s Exit 0
ok  	obs-lan-controller	0.219s
ok  	obs-lan-controller/Lib	4.250s
?   	obs-lan-controller/web	[no test files]

go vet ./... Exit 0

go build -o /root/.hermes/profiles/coder/cache/scratch/obs-controller-pause . Exit 0
