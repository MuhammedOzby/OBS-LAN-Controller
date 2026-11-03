# Disabled control tooltips
UTC 2026-10-29T10:14:54.201450+00:00
Linux Go/Node mock frontend and backend regression; no real OBS mutation.

go test ./... -count=1 -timeout=60s Exit 0
ok  	obs-lan-controller	0.257s
ok  	obs-lan-controller/Lib	4.252s
?   	obs-lan-controller/web	[no test files]

go vet ./... Exit 0

go build -o /root/.hermes/profiles/coder/cache/scratch/obs-controller-tooltip . Exit 0
