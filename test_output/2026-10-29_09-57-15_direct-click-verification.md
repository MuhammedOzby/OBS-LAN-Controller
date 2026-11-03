# Direct click verification
UTC: 2026-10-29T09:57:15.456046+00:00
Linux Go/Node, mock OBS only. Previous grouped invocation hit terminal 120s timeout (124), incomplete report. Retry with explicit per-command limits. No real OBS mutation.

`go test ./... -count=1 -timeout=30s`
Exit: 1
```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandler.func1 (2 handlers)
[GIN-debug] POST   /api/broadcast            --> obs-lan-controller/Lib.newHandler.broadcastHandler.func4 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandler.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
panic: test timed out after 30s
	running tests:
		TestBroadcastFrontendBehavior (30s)

goroutine 10 [running]:
testing.(*M).startAlarm.func1()
	/usr/local/go/src/testing/testing.go:2484 +0x394
created by time.goFunc
	/usr/local/go/src/time/sleep.go:215 +0x2d

goroutine 1 [chan receive]:
testing.(*T).Run(0xc000003a40, {0xfc1eaa?, 0xc0000d7b30?}, 0x10556b8)
	/usr/local/go/src/testing/testing.go:1859 +0x431
testing.runTests.func1(0xc000003a40)
	/usr/local/go/src/testing/testing.go:2279 +0x37
testing.tRunner(0xc000003a40, 0xc0000d7c70)
	/usr/local/go/src/testing/testing.go:1792 +0xf4
testing.runTests(0xc00000eab0, {0x189eee0, 0x9, 0x9}, {0x18b0560?, 0x7?, 0x18ae340?})
	/usr/local/go/src/testing/testing.go:2277 +0x4b4
testing.(*M).Run(0xc0000e2aa0)
	/usr/local/go/src/testing/testing.go:2142 +0x64a
main.main()
	_testmain.go:61 +0x9b

goroutine 8 [syscall]:
syscall.Syscall6(0xf7, 0x3, 0xa, 0xc000098d20, 0x4, 0xc0001b81b0, 0x0)
	/usr/local/go/src/syscall/syscall_linux.go:95 +0x39
internal/syscall/unix.Waitid(0xc000098d4e?, 0xc000098e78?, 0x58018b?, 0xc00003c840?, 0x1?)
	/usr/local/go/src/internal/syscall/unix/waitid_linux.go:18 +0x39
os.(*Process).pidfdWait.func1(...)
	/usr/local/go/src/os/pidfd_linux.go:106
os.ignoringEINTR(...)
	/usr/local/go/src/os/file_posix.go:251
os.(*Process).pidfdWait(0xc00000fd28?)
	/usr/local/go/src/os/pidfd_linux.go:105 +0x209
os.(*Process).wait(0x18b0560?)
	/usr/local/go/src/os/exec_unix.go:27 +0x25
os.(*Process).Wait(...)
	/usr/local/go/src/os/exec.go:358
os/exec.(*Cmd).Wait(0xc0000ca480)
	/usr/local/go/src/os/exec/exec.go:922 +0x45
os/exec.(*Cmd).Run(0xc0000ca480)
	/usr/local/go/src/os/exec/exec.go:626 +0x2d
os/exec.(*Cmd).CombinedOutput(0xc0000ca480)
	/usr/local/go/src/os/exec/exec.go:1039 +0x85
obs-lan-controller.TestBroadcastFrontendBehavior(0xc0002d6000)
	/root/hermes-work/calisma/OBS yayın araçları geliştirme/OBS LAN Controller/broadcast_frontend_test.go:61 +0x10a
testing.tRunner(0xc0002d6000, 0x10556b8)
	/usr/local/go/src/testing/testing.go:1792 +0xf4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1851 +0x413

goroutine 9 [IO wait]:
internal/poll.runtime_pollWait(0x7f2456bc6670, 0x72)
	/usr/local/go/src/runtime/netpoll.go:351 +0x85
internal/poll.(*pollDesc).wait(0xc0001b3620?, 0xc000208000?, 0x1)
	/usr/local/go/src/internal/poll/fd_poll_runtime.go:84 +0x27
internal/poll.(*pollDesc).waitRead(...)
	/usr/local/go/src/internal/poll/fd_poll_runtime.go:89
internal/poll.(*FD).Read(0xc0001b3620, {0xc000208000, 0x200, 0x200})
	/usr/local/go/src/internal/poll/fd_unix.go:165 +0x27a
os.(*File).read(...)
	/usr/local/go/src/os/file_posix.go:29
os.(*File).Read(0xc00006a768, {0xc000208000?, 0x41691f?, 0xeb6d20?})
	/usr/local/go/src/os/file.go:124 +0x4f
bytes.(*Buffer).ReadFrom(0xc0002c4540, {0x11e11c0, 0xc000206000})
	/usr/local/go/src/bytes/buffer.go:211 +0x98
io.copyBuffer({0x11e1700, 0xc0002c4540}, {0x11e11c0, 0xc000206000}, {0x0, 0x0, 0x0})
	/usr/local/go/src/io/io.go:415 +0x151
io.Copy(...)
	/usr/local/go/src/io/io.go:388
os.genericWriteTo(0xc00006a768?, {0x11e1700, 0xc0002c4540})
	/usr/local/go/src/os/file.go:275 +0x4f
os.(*File).WriteTo(0xc00006a768, {0x11e1700, 0xc0002c4540})
	/usr/local/go/src/os/file.go:253 +0x9c
io.copyBuffer({0x11e1700, 0xc0002c4540}, {0x11e1240, 0xc00006a768}, {0x0, 0x0, 0x0})
	/usr/local/go/src/io/io.go:411 +0x9d
io.Copy(...)
	/usr/local/go/src/io/io.go:388
os/exec.(*Cmd).writerDescriptor.func1()
	/usr/local/go/src/os/exec/exec.go:596 +0x34
os/exec.(*Cmd).Start.func2(0x10556b8?)
	/usr/local/go/src/os/exec/exec.go:749 +0x2c
created by os/exec.(*Cmd).Start in goroutine 8
	/usr/local/go/src/os/exec/exec.go:748 +0x92b
FAIL	obs-lan-controller	30.032s
ok  	obs-lan-controller/Lib	0.534s
?   	obs-lan-controller/web	[no test files]
FAIL

```

`go vet ./...`
Exit: 0
```

```

`go build -o /root/.hermes/profiles/coder/cache/scratch/obs-controller-direct .`
Exit: 0
```

```
