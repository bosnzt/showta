appVersion=1.0.1
ldflags="\
-X 'showta.cc/app/system/conf.AppVersion=$appVersion' \
"
go build -o showta.exe -ldflags="$ldflags" .