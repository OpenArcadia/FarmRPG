export CC="aarch64-linux-android${ANDROID_API}-clang"
export CGO_CFLAGS="--sysroot=${ANDROID_SYSROOT}"
export CGO_LDFLAGS="--sysroot=${ANDROID_SYSROOT}"
export CGO_ENABLED=1
export GOOS=android
export GOARCH=arm64
go build -buildmode=c-shared \
  -ldflags="-s -w -extldflags=-Wl,-soname,libfarmrpg.so" \
  -o=android/libs/arm64-v8a/libfarmrpg.so
