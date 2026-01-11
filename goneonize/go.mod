module github.com/krypton-byte/neonize

go 1.25.3

require (
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.33
	go.mau.fi/util v0.9.4
	go.mau.fi/whatsmeow v0.0.0-20260107124630-ccfa04f8e445
	google.golang.org/protobuf v1.36.11
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beeper/argo-go v1.1.2 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/elliotchance/orderedmap/v3 v3.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/petermattis/goid v0.0.0-20251121121749-a11dd1a45f9a // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/vektah/gqlparser/v2 v2.5.31 // indirect
	go.mau.fi/libsignal v0.2.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/exp v0.0.0-20251219203646-944ab1f22d93 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace go.mau.fi/libsignal => github.com/fawwaz37/libsignal-protocol-go v0.2.1-0.20250920101933-ae5672c024d5

replace go.mau.fi/whatsmeow => github.com/ZhanNexus/whatsmeow v0.0.0-20260102020949-f26186f153e6 // github.com/ginkohub/whatsmeow v0.0.0-20251202021103-f3779ce15345
