go build -o farmrpg -ldflags="-s -w"
flatpak-builder build-dir --force-clean farmrpg.flatpak.yaml 
flatpak build-export repo build-dir
flatpak build-bundle repo farmrpg.flatpak com.openarchadia.farmrpg
