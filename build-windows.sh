mkdir -p windows-build
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o windows-build/farmrpg.exe -ldflags "-s -w" 
mkdir -p windows-build/assets
mkdir -p windows-build/data
cp -r assets/* windows-build/assets/
cp -r data/* windows-build/data/
cp raylib.dll windows-build/
tar -czvf farmrpg_windows.tar.gz -C windows-build .